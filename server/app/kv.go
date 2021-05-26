package app

import (
	"github.com/poonman/entry-task/server/domain/aggr/kv"
)

func (s *Service) WriteSecureMessage(username, key, value string) (err error) {
	err = kv.ValidateKey(key)
	if err != nil {
		return
	}

	tmpKey := kv.NewKey(username, key)

	err = s.kvRepo.Set(tmpKey, value)
	if err != nil {

	}

	return
}

func (s *Service) ReadSecureMessage(username, key string) (value string, err error) {
	err = kv.ValidateKey(key)
	if err != nil {
		return
	}

	tmpKey := kv.NewKey(username, key)

	value, err = s.kvRepo.Get(tmpKey)
	if err != nil {

	}

	return
}
