package main

import (
	"github.com/poonman/entry-task/client/app"
	"github.com/poonman/entry-task/client/cmd/internal"
	"github.com/poonman/entry-task/client/domain/gateway"
	"github.com/poonman/entry-task/client/infra/config"
	"github.com/poonman/entry-task/dora/misc/helper"
	"github.com/poonman/entry-task/dora/misc/log"
)

func main() {

	arg := internal.CommandLine()

	log.Infof("arguments:[%+v]", arg)

	c := internal.BuildContainer()

	helper.MustContainerInvoke(c, func(conf *config.Config, appSvc *app.Service, kvGateway gateway.KvGateway) {

		log.SetLevelByString(conf.LogConfig.Level)

		var (
			k, v byte
		)

		if len(arg.Key) == 1 {
			k = []byte(arg.Key)[0]
		}

		if len(arg.Value) == 1 {
			v = []byte(arg.Value)[0]
		}

		for _, cmd := range arg.Commands {
			switch cmd {
			case "benchmark":

				appSvc.Benchmark(arg.Concurrency, arg.Requests, arg.Username, arg.Password, arg.Method, k, v)

			case "login":
				err := appSvc.Login(arg.Username, arg.Password)
				if err != nil {
					break
				}
			case "write":

				err := appSvc.WriteSecureMessage(arg.Username, arg.Password, k, v)
				if err != nil {
					break
				}
			case "read":
				err := appSvc.ReadSecureMessage(arg.Username, arg.Password, k)
				if err != nil {
					break
				}
			}
		}

		kvGateway.Stop()
	})
}
