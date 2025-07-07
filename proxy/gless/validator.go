package gless

import (
	"sync"
	"time"

	"github.com/xtls/xray-core/common/protocol"
)

// MemoryValidator is an in-memory form of GLess user validator.
type MemoryValidator struct {
	users map[string]*MemoryAccount
	mu    sync.RWMutex
}

// NewMemoryValidator creates a new MemoryValidator.
func NewMemoryValidator() *MemoryValidator {
	return &MemoryValidator{
		users: make(map[string]*MemoryAccount),
	}
}

// Add adds a new user.
func (v *MemoryValidator) Add(user *protocol.MemoryUser) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	account := user.Account.(*MemoryAccount)
	v.users[account.ID.String()] = account
	return nil
}

// Del removes a user by its ID.
func (v *MemoryValidator) Del(id string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if _, found := v.users[id]; !found {
		return false
	}
	delete(v.users, id)
	return true
}

// Get returns a user by its ID.
func (v *MemoryValidator) Get(id string) (*MemoryAccount, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	user, found := v.users[id]
	return user, found
}

// GetAll returns all users.
func (v *MemoryValidator) GetAll() []*MemoryAccount {
	v.mu.RLock()
	defer v.mu.RUnlock()

	users := make([]*MemoryAccount, 0, len(v.users))
	for _, user := range v.users {
		users = append(users, user)
	}
	return users
}

// GetByEmail returns a user by email (ID in this case).
func (v *MemoryValidator) GetByEmail(email string) *protocol.MemoryUser {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if user, found := v.users[email]; found {
		return &protocol.MemoryUser{
			Account: user,
		}
	}
	return nil
}

// GetCount returns the number of users.
func (v *MemoryValidator) GetCount() int64 {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return int64(len(v.users))
}

// Close implements common.Closable.
func (v *MemoryValidator) Close() error {
	return nil
}

// TimedUserValidator is a GLess user validator with time-based validation.
type TimedUserValidator struct {
	users    map[string]*MemoryAccount
	lastSeen map[string]time.Time
	mu       sync.RWMutex
}

// NewTimedUserValidator creates a new TimedUserValidator.
func NewTimedUserValidator() *TimedUserValidator {
	return &TimedUserValidator{
		users:    make(map[string]*MemoryAccount),
		lastSeen: make(map[string]time.Time),
	}
}

// Add adds a new user.
func (v *TimedUserValidator) Add(user *protocol.MemoryUser) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	account := user.Account.(*MemoryAccount)
	v.users[account.ID.String()] = account
	v.lastSeen[account.ID.String()] = time.Now()
	return nil
}

// Del removes a user by its ID.
func (v *TimedUserValidator) Del(id string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if _, found := v.users[id]; !found {
		return false
	}
	delete(v.users, id)
	delete(v.lastSeen, id)
	return true
}

// Get returns a user by its ID and updates last seen time.
func (v *TimedUserValidator) Get(id string) (*MemoryAccount, bool) {
	v.mu.Lock()
	defer v.mu.Unlock()

	user, found := v.users[id]
	if found {
		v.lastSeen[id] = time.Now()
	}
	return user, found
}

// GetAll returns all users.
func (v *TimedUserValidator) GetAll() []*MemoryAccount {
	v.mu.RLock()
	defer v.mu.RUnlock()

	users := make([]*MemoryAccount, 0, len(v.users))
	for _, user := range v.users {
		users = append(users, user)
	}
	return users
}

// Close implements common.Closable.
func (v *TimedUserValidator) Close() error {
	return nil
} 