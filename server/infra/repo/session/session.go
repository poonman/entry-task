package session

import (
	"github.com/poonman/entry-task/server/domain/aggr/session"
	"sync"
)

type repo struct {
	mu sync.RWMutex
	userTokens map[string]string
}

func (r *repo) Save(username, token string) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.userTokens[username] = token

	return nil
}

func (r *repo) Get(username string) (token string, err error) {
	r.mu.RLock()
	r.mu.RUnlock()

	token, ok := r.userTokens[username]
	if !ok {
		return "", nil
	}

	return token, nil
}

func NewRepo() session.Repo {
	return &repo{
		mu:    sync.RWMutex{},
		userTokens: make(map[string]string),
	}
}