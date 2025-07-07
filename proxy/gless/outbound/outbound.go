package outbound

import (
	"context"
	"io"

	"github.com/xtls/xray-core/common"
	"github.com/xtls/xray-core/common/buf"
	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/task"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/features/dns"
	"github.com/xtls/xray-core/features/policy"
	"github.com/xtls/xray-core/proxy/gless"
	"github.com/xtls/xray-core/proxy/gless/encoding"
	"github.com/xtls/xray-core/transport"
	"github.com/xtls/xray-core/transport/internet"
)

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}

// Handler is an outbound connection handler for GLess protocol.
type Handler struct {
	serverList    *protocol.ServerList
	policyManager policy.Manager
	dns           dns.Client
}

// New creates a new GLess outbound handler.
func New(ctx context.Context, config *Config) (*Handler, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Vnext {
		s, err := protocol.NewServerSpecFromPB(rec)
		if err != nil {
			return nil, errors.New("failed to parse server spec").Base(err)
		}
		serverList.AddServer(s)
	}

	v := core.MustFromContext(ctx)
	return &Handler{
		serverList:    serverList,
		policyManager: v.GetFeature(policy.ManagerType()).(policy.Manager),
		dns:           v.GetFeature(dns.ClientType()).(dns.Client),
	}, nil
}

// Process implements proxy.Outbound.Process().
func (h *Handler) Process(ctx context.Context, link *transport.Link, dialer internet.Dialer) error {
	var rec *protocol.ServerSpec
	var conn net.Conn

	rec = h.serverList.GetServer(0) // Get first server
	if rec == nil {
		return errors.New("no available server")
	}

	dest := rec.Destination()
	conn, err := dialer.Dial(ctx, dest)
	if err != nil {
		return errors.New("failed to dial to GLess server").Base(err)
	}
	defer conn.Close()

	user := rec.PickUser() // Pick a user
	if user == nil {
		return errors.New("no available user")
	}

	account := user.Account.(*gless.MemoryAccount)
	gostCipher, err := gless.NewGOSTCipher(account.ID.String(), false)
	if err != nil {
		return errors.New("failed to create GOST cipher").Base(err)
	}

	// Create request header
	request := &protocol.RequestHeader{
		Command: protocol.RequestCommandTCP,
		Address: net.DomainAddress("example.com"),
		Port:    80,
	}

	// Encode request header
	buffer := buf.New()
	defer buffer.Release()

	if err := encoding.EncodeRequestHeader(buffer, request, nil); err != nil {
		return errors.New("failed to encode request header").Base(err)
	}

	// Write request header
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return errors.New("failed to write request header").Base(err)
	}

	// Create encrypted writer
	encryptedWriter := encoding.NewGOSTWriter(conn, gostCipher)

	// Handle the connection
	return h.handleConnection(ctx, link, encryptedWriter)
}

func (h *Handler) handleConnection(ctx context.Context, link *transport.Link, writer io.Writer) error {
	requestDone := func() error {
		return buf.Copy(link.Reader, buf.NewWriter(writer))
	}

	responseDone := func() error {
		// Create a reader from the writer (this is a simplified approach)
		// In a real implementation, you'd need to handle bidirectional communication properly
		return nil
	}

	var responseDoneAndCloseWriter = task.OnSuccess(responseDone, task.Close(link.Writer))
	if err := task.Run(ctx, requestDone, responseDoneAndCloseWriter); err != nil {
		common.Interrupt(link.Reader)
		common.Interrupt(link.Writer)
		return errors.New("connection ends").Base(err)
	}

	return nil
} 