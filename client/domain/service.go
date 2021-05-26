package domain

import "github.com/poonman/entry-task/client/domain/gateway"

type Service struct {
	kvGateway gateway.KvGateway
}

func NewService(
	kvGateway gateway.KvGateway,
	) *Service {

	s := &Service{
		kvGateway: kvGateway,
	}

	return s
}