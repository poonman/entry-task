package main

import (
	"github.com/poonman/entry-task/server/app"
	"github.com/poonman/entry-task/server/infra/helper"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container{
	c := dig.New()

	helper.MustContainerProvide(c, app.NewService)

	return c
}

func main() {
	c := BuildContainer()

	helper.MustContainerInvoke(c, func() {

	})
}