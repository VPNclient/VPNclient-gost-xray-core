package conf

import (
	"encoding/json"
	"strings"

	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/common/uuid"
	"github.com/xtls/xray-core/proxy/gmess"
	"github.com/xtls/xray-core/proxy/gmess/inbound"
	"github.com/xtls/xray-core/proxy/gmess/outbound"
	"google.golang.org/protobuf/proto"
)

type GMessAccount struct {
	ID          string `json:"id"`
	Security    string `json:"security"`
	Experiments string `json:"experiments"`
}

// Build implements Buildable
func (a *GMessAccount) Build() *gmess.Account {
	var st protocol.SecurityType
	switch strings.ToLower(a.Security) {
	case "gost-28147":
		st = protocol.SecurityType_AES128_GCM // Using GOST equivalent
	case "gost-89":
		st = protocol.SecurityType_CHACHA20_POLY1305 // Using GOST equivalent
	case "auto":
		st = protocol.SecurityType_AUTO
	case "none":
		st = protocol.SecurityType_NONE
	case "zero":
		st = protocol.SecurityType_ZERO
	default:
		st = protocol.SecurityType_AUTO
	}
	return &gmess.Account{
		Id: a.ID,
		SecuritySettings: &protocol.SecurityConfig{
			Type: st,
		},
		TestsEnabled: a.Experiments,
	}
}

type GMessDetourConfig struct {
	ToTag string `json:"to"`
}

// Build implements Buildable
func (c *GMessDetourConfig) Build() *inbound.DetourConfig {
	return &inbound.DetourConfig{
		To: c.ToTag,
	}
}

type GMessDefaultConfig struct {
	Level byte `json:"level"`
}

// Build implements Buildable
func (c *GMessDefaultConfig) Build() *inbound.DefaultConfig {
	config := new(inbound.DefaultConfig)
	config.Level = uint32(c.Level)
	return config
}

type GMessInboundConfig struct {
	Users        []json.RawMessage   `json:"clients"`
	Defaults     *GMessDefaultConfig `json:"default"`
	DetourConfig *GMessDetourConfig  `json:"detour"`
}

// Build implements Buildable
func (c *GMessInboundConfig) Build() (proto.Message, error) {
	config := &inbound.Config{}

	if c.Defaults != nil {
		config.Default = c.Defaults.Build()
	}

	if c.DetourConfig != nil {
		config.Detour = c.DetourConfig.Build()
	}

	config.User = make([]*protocol.User, len(c.Users))
	for idx, rawData := range c.Users {
		user := new(protocol.User)
		if err := json.Unmarshal(rawData, user); err != nil {
			return nil, errors.New("invalid GMess user").Base(err)
		}
		account := new(GMessAccount)
		if err := json.Unmarshal(rawData, account); err != nil {
			return nil, errors.New("invalid GMess user").Base(err)
		}

		u, err := uuid.ParseString(account.ID)
		if err != nil {
			return nil, err
		}
		account.ID = u.String()

		user.Account = serial.ToTypedMessage(account.Build())
		config.User[idx] = user
	}

	return config, nil
}

type GMessOutboundTarget struct {
	Address *Address          `json:"address"`
	Port    uint16            `json:"port"`
	Users   []json.RawMessage `json:"users"`
}

type GMessOutboundConfig struct {
	Receivers []*GMessOutboundTarget `json:"vnext"`
}

// Build implements Buildable
func (c *GMessOutboundConfig) Build() (proto.Message, error) {
	config := new(outbound.Config)

	if len(c.Receivers) == 0 {
		return nil, errors.New("0 GMess receiver configured")
	}
	serverSpecs := make([]*protocol.ServerEndpoint, len(c.Receivers))
	for idx, rec := range c.Receivers {
		if len(rec.Users) == 0 {
			return nil, errors.New("0 user configured for GMess outbound")
		}
		if rec.Address == nil {
			return nil, errors.New("address is not set in GMess outbound config")
		}
		spec := &protocol.ServerEndpoint{
			Address: rec.Address.Build(),
			Port:    uint32(rec.Port),
		}
		for _, rawUser := range rec.Users {
			user := new(protocol.User)
			if err := json.Unmarshal(rawUser, user); err != nil {
				return nil, errors.New("invalid GMess user").Base(err)
			}
			account := new(GMessAccount)
			if err := json.Unmarshal(rawUser, account); err != nil {
				return nil, errors.New("invalid GMess user").Base(err)
			}

			u, err := uuid.ParseString(account.ID)
			if err != nil {
				return nil, err
			}
			account.ID = u.String()

			user.Account = serial.ToTypedMessage(account.Build())
			spec.User = append(spec.User, user)
		}
		serverSpecs[idx] = spec
	}
	config.Receiver = serverSpecs
	return config, nil
}
