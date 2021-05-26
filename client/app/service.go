package app

import "github.com/poonman/entry-task/client/domain"

type Service struct {
	domainSvc *domain.Service
}

func NewService(domainSvc *domain.Service) *Service {
	return &Service{
		domainSvc: domainSvc,
	}
}
