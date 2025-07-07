// Package gless contains the implementation of GLess protocol and transportation.
//
// GLess contains both inbound and outbound connections. GLess inbound is usually used on servers
// together with 'freedom' to talk to final destination, while GLess outbound is usually used on
// clients with 'socks' for proxying.
//
// GLess is a GOST-based protocol similar to VLESS but uses Russian GOST cryptography standards.
package gless

import (
	"github.com/xtls/xray-core/common/protocol"
)

const (
	// GOST-based flow identifier
	GRV = "gost-rprx-vision"
)

// Validator interface for GLess user validation
type Validator interface {
	Add(user *protocol.MemoryUser) error
	Del(id string) bool
	Get(id string) (*MemoryAccount, bool)
	GetByEmail(email string) *protocol.MemoryUser
	GetAll() []*MemoryAccount
	GetCount() int64
	Close() error
} 