package main

import (
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/dora/misc/helper"
	"github.com/poonman/entry-task/dora/server"
	"github.com/poonman/entry-task/server/api"
	"github.com/poonman/entry-task/server/api/interceptor"
	"github.com/poonman/entry-task/server/app"
	"github.com/poonman/entry-task/server/domain/factory"
	"github.com/poonman/entry-task/server/idl/kv"
	"github.com/poonman/entry-task/server/infra/config"
	"github.com/poonman/entry-task/server/infra/driver/redis"
	"github.com/poonman/entry-task/server/infra/repo/account"
	"github.com/poonman/entry-task/server/infra/repo/limiter"
	"github.com/poonman/entry-task/server/infra/repo/quota"
	"github.com/poonman/entry-task/server/infra/repo/session"
	"github.com/poonman/entry-task/server/infra/repo/store"
	"go.uber.org/dig"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "golang.org/x/time/rate"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	helper.MustContainerProvide(c, config.NewConfig)
	helper.MustContainerProvide(c, redis.NewRedisPool)

	helper.MustContainerProvide(c, factory.NewFactory)
	helper.MustContainerProvide(c, app.NewService)
	helper.MustContainerProvide(c, api.NewHandler)
	helper.MustContainerProvide(c, interceptor.NewInterceptor)

	// repo
	helper.MustContainerProvide(c, account.NewRepo)
	helper.MustContainerProvide(c, quota.NewRepo)
	helper.MustContainerProvide(c, store.NewRepo)
	helper.MustContainerProvide(c, session.NewRepo)
	helper.MustContainerProvide(c, limiter.NewRepo)

	return c
}

func main() {
	c := BuildContainer()

	helper.MustContainerInvoke(c, func(conf *config.Config, interceptor *interceptor.Interceptor, h kv.StoreServer) {

		log.Debug("start...")
		var dora *server.Server

		if conf.ServerConfig.EnableTls {
			tlsConfig := conf.LoadTLSConfig()
			dora = server.NewServer(server.WithTlsConfig(tlsConfig), server.WithInterceptor(interceptor.Do))
		} else {
			dora = server.NewServer(server.WithInterceptor(interceptor.Do))
		}

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
