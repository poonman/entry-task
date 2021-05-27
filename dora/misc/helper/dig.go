package helper

import (
	"github.com/poonman/entry-task/dora/log"
	"go.uber.org/dig"
)

func MustContainerProvide(c *dig.Container, constructor interface{}, opts ...dig.ProvideOption) {
	err := c.Provide(constructor, opts...)
	if err != nil {
		log.Fatalf("failed to provide constructor. %v", err)
	}
}

func MustContainerInvoke(c *dig.Container, function interface{}, opts ...dig.InvokeOption) {
	err := c.Invoke(function, opts...)
	if err != nil {
		log.Fatalf("invoke error [%v]", err)
	}
}
