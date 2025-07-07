package inbound

import (
	"context"
	"io"

	"github.com/xtls/xray-core/common"
	"github.com/xtls/xray-core/common/buf"
	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/log"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/session"
	"github.com/xtls/xray-core/common/task"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/features/dns"
	feature_inbound "github.com/xtls/xray-core/features/inbound"
	"github.com/xtls/xray-core/features/policy"
	"github.com/xtls/xray-core/features/routing"
	"github.com/xtls/xray-core/proxy/gless"
	"github.com/xtls/xray-core/proxy/gless/encoding"
	"github.com/xtls/xray-core/transport/internet/stat"
)

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		var dc dns.Client
		if err := core.RequireFeatures(ctx, func(d dns.Client) error {
			dc = d
			return nil
		}); err != nil {
			return nil, err
		}

		c := config.(*Config)

		validator := gless.NewMemoryValidator()
		for _, user := range c.Clients {
			u, err := user.ToMemoryUser()
			if err != nil {
				return nil, errors.New("failed to get GLess user").Base(err).AtError()
			}
			if err := validator.Add(u); err != nil {
				return nil, errors.New("failed to initiate user").Base(err).AtError()
			}
		}

		return New(ctx, c, dc, validator)
	}))
}

// Handler is an inbound connection handler that handles messages in GLess protocol.
type Handler struct {
	inboundHandlerManager feature_inbound.Manager
	policyManager         policy.Manager
	validator             gless.Validator
	dns                   dns.Client
	fallbacks             map[string]map[string]map[string]*Fallback
}

// New creates a new GLess inbound handler.
func New(ctx context.Context, config *Config, dc dns.Client, validator gless.Validator) (*Handler, error) {
	v := core.MustFromContext(ctx)
	handler := &Handler{
		inboundHandlerManager: v.GetFeature(feature_inbound.ManagerType()).(feature_inbound.Manager),
		policyManager:         v.GetFeature(policy.ManagerType()).(policy.Manager),
		dns:                   dc,
		validator:             validator,
	}

	if config.Fallbacks != nil {
		handler.fallbacks = make(map[string]map[string]map[string]*Fallback)
		for _, fb := range config.Fallbacks {
			if handler.fallbacks[fb.Name] == nil {
				handler.fallbacks[fb.Name] = make(map[string]map[string]*Fallback)
			}
			if handler.fallbacks[fb.Name][fb.Alpn] == nil {
				handler.fallbacks[fb.Name][fb.Alpn] = make(map[string]*Fallback)
			}
			handler.fallbacks[fb.Name][fb.Alpn][fb.Path] = fb
		}
	}

	return handler, nil
}

// Close implements common.Closable.Close().
func (h *Handler) Close() error {
	return errors.Combine(common.Close(h.validator))
}

// AddUser implements proxy.UserManager.AddUser().
func (h *Handler) AddUser(ctx context.Context, u *protocol.MemoryUser) error {
	return h.validator.Add(u)
}

// RemoveUser implements proxy.UserManager.RemoveUser().
func (h *Handler) RemoveUser(ctx context.Context, e string) error {
	if h.validator.Del(e) {
		return nil
	}
	return errors.New("user not found")
}

// GetUser implements proxy.UserManager.GetUser().
func (h *Handler) GetUser(ctx context.Context, email string) *protocol.MemoryUser {
	return h.validator.GetByEmail(email)
}

// GetUsers implements proxy.UserManager.GetUsers().
func (h *Handler) GetUsers(ctx context.Context) []*protocol.MemoryUser {
	accounts := h.validator.GetAll()
	users := make([]*protocol.MemoryUser, 0, len(accounts))
	for _, acc := range accounts {
		users = append(users, &protocol.MemoryUser{Account: acc})
	}
	return users
}

// GetUsersCount implements proxy.UserManager.GetUsersCount().
func (h *Handler) GetUsersCount(context.Context) int64 {
	return h.validator.GetCount()
}

// Network implements proxy.Inbound.Network().
func (*Handler) Network() []net.Network {
	return []net.Network{net.Network_TCP, net.Network_UNIX}
}

// Process implements proxy.Inbound.Process().
func (h *Handler) Process(ctx context.Context, network net.Network, connection stat.Connection, dispatcher routing.Dispatcher) error {
	iConn := connection
	if statConn, ok := iConn.(*stat.CounterConnection); ok {
		iConn = statConn.Connection
	}

	// Use iConn (io.Reader) for DecodeRequestHeader
	request, err := encoding.DecodeRequestHeader(iConn, h.validator)
	if err != nil {
		if errors.Cause(err) != io.EOF {
			log.Record(&log.AccessMessage{
				From:   iConn.RemoteAddr(),
				To:     "",
				Status: log.AccessRejected,
				Reason: err,
			})
			err = errors.New("invalid request from ", iConn.RemoteAddr()).Base(err)
		}
		return err
	}

	if request.Command != protocol.RequestCommandTCP {
		return errors.New("unsupported command: ", request.Command)
	}

	dest := request.Destination()
	ctx = session.ContextWithInbound(ctx, &session.Inbound{
		Source:  net.DestinationFromAddr(iConn.RemoteAddr()),
		Gateway: dest,
		Tag:     "",
	})
	ctx = log.ContextWithAccessMessage(ctx, &log.AccessMessage{
		From:   iConn.RemoteAddr(),
		To:     dest,
		Status: log.AccessAccepted,
		Reason: "",
	})

	userSettings := ctx.Value("user").(*protocol.MemoryUser)
	account := userSettings.Account.(*gless.MemoryAccount)

	// Create GOST cipher for decryption
	gostCipher, err := gless.NewGOSTCipher(account.ID.String(), true)
	if err != nil {
		return errors.New("failed to create GOST cipher").Base(err)
	}

	decryptedReader := encoding.NewGOSTReader(iConn, gostCipher)

	return h.handleConnection(ctx, decryptedReader, dest, dispatcher)
}

func (h *Handler) handleConnection(ctx context.Context, reader io.Reader, dest net.Destination, dispatcher routing.Dispatcher) error {
	link, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return errors.New("failed to dispatch request").Base(err)
	}

	requestDone := func() error {
		return buf.Copy(buf.NewReader(reader), link.Writer)
	}

	responseDone := func() error {
		return buf.Copy(link.Reader, link.Writer)
	}

	var responseDoneAndCloseWriter = task.OnSuccess(responseDone, func() error { return common.Close(link.Writer) })
	if err := task.Run(ctx, requestDone, responseDoneAndCloseWriter); err != nil {
		common.Interrupt(link.Reader)
		common.Interrupt(link.Writer)
		return errors.New("connection ends").Base(err)
	}

	return nil
} 