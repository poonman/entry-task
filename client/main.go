package main

import (
	"github.com/poonman/entry-task/client/domain"
	"github.com/poonman/entry-task/client/infra/config"
	"github.com/poonman/entry-task/client/infra/gateway"
	"github.com/poonman/entry-task/dora/misc/helper"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	helper.MustContainerProvide(c, config.NewConfig)
	helper.MustContainerProvide(c, gateway.NewKvGateway)
	//helper.MustContainerProvide(c, app.NewService)
	helper.MustContainerProvide(c, domain.NewService)

	return c
}

func main() {
	c := BuildContainer()

	helper.MustContainerInvoke(c, func(svc *domain.Service) {
		svc.RunJob()
	})
}
