package userm

import (
	"github.com/poonman/entry-task/client/domain/aggr/user"
	"sync"
)

type UserManager struct {
	sync.RWMutex
	userMap map[string]*user.User
}

func NewUserManager() *UserManager {
	return &UserManager{
		RWMutex: sync.RWMutex{},
		userMap: make(map[string]*user.User),
	}
}

func (m *UserManager) AddUser(u *user.User) {
	m.Lock()
	defer m.Unlock()

	m.userMap[u.Name] = u
}

func (m *UserManager) GetUser(username string) (u *user.User) {
	m.RLock()
	defer m.RUnlock()

	var ok bool

	u, ok = m.userMap[username]
	if !ok {
		return nil
	}

	return
}
