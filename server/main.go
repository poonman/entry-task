package main

import (
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/dora/misc/helper"
	"github.com/poonman/entry-task/dora/server"
	"github.com/poonman/entry-task/server/api"
	"github.com/poonman/entry-task/server/api/interceptor"
	"github.com/poonman/entry-task/server/app"
	"github.com/poonman/entry-task/server/idl/kv"
	"github.com/poonman/entry-task/server/infra/config"
	"github.com/poonman/entry-task/server/infra/driver/redis"
	"go.uber.org/dig"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	helper.MustContainerProvide(c, config.NewConfig)
	helper.MustContainerProvide(c, redis.NewRedisPool)

	helper.MustContainerProvide(c, app.NewService)
	helper.MustContainerProvide(c, api.NewHandler)
	helper.MustContainerProvide(c, interceptor.NewAuthInterceptor)

	return c
}

func main() {
	c := BuildContainer()

	helper.MustContainerInvoke(c, func(conf *config.Config, interceptor *interceptor.AuthInterceptor, h kv.StoreServer) {

		log.Debug("start...")
		//tlsConfig := conf.LoadTLSConfig()
		//dora := server.NewServer(server.WithTlsConfig(tlsConfig))

		dora := server.NewServer(server.WithInterceptor(interceptor.Auth))

		kv.RegisterStoreServer(dora, h)

		go func() {
			err := dora.Serve(conf.ServerConfig.ListenAddress)
			if err != nil {
				log.Errorf("Failed to serve address. address:[%s], err:[%v]", conf.ServerConfig.ListenAddress, err)
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			sig := <-c
			log.Infof("capture a signal. signal:[%s]", sig.String())
			switch sig {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:

				dora.Stop()

				time.Sleep(time.Second)
				return
			case syscall.SIGHUP:
			default:
				return
			}
		}
	})
}
