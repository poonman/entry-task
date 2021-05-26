package app

import (
	"github.com/poonman/entry-task/server/domain/aggr/kv"
)

func (s *Service) WriteSecureMessage(uid uint64, key, value string) (err error) {
	err = kv.ValidateKey(key)
	if err != nil {
		return
	}

	tmpKey := kv.NewKey(uid, key)

	err = s.kvRepo.Set(tmpKey, value)
	if err != nil {

	}

	return
}

func (s *Service) ReadSecureMessage(uid uint64, key string) (value string, err error) {
	err = kv.ValidateKey(key)
	if err != nil {
		return
	}

	tmpKey := kv.NewKey(uid, key)

	value, err = s.kvRepo.Get(tmpKey)
	if err != nil {

	}

	return
}
