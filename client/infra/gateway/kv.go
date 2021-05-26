package gateway

import (
	"context"
	"github.com/poonman/entry-task/client/domain/aggr/user"
	"github.com/poonman/entry-task/client/domain/gateway"
	"github.com/poonman/entry-task/client/idl/kv"
	"github.com/poonman/entry-task/client/infra/config"
	"github.com/poonman/entry-task/dora/client"
	"github.com/poonman/entry-task/dora/log"
)

type kvGateway struct {
	client kv.StoreClient
}

func (g *kvGateway) Auth(u *user.User) (err error) {
	return
}

func (g *kvGateway) Set(u *user.User, key, value string) (err error) {
	req := &kv.WriteSecureMessageReq{
		Key:                  key,
		Value:                value,
	}

	ctx := context.TODO()

	_, err = g.client.WriteSecureMessage(ctx, req)
	if err != nil {
		log.Errorf("Failed to write secure message. user:[%+v], key:[%s], value:[%s], err:[%v]",
			u, key, value, err)
		return
	}
	return
}

func (g *kvGateway) Get(u *user.User, key string) (value string, err error) {
	req := &kv.ReadSecureMessageReq{
		Key:                  key,
	}

	ctx := context.TODO()

	rsp, err := g.client.ReadSecureMessage(ctx, req)
	if err != nil {
		log.Errorf("Failed to read secure message. user:[%+v], key:[%s], err:[%v]",
			u, key, err)
		return
	}

	value = rsp.Value

	return
}

func NewKvGateway(conf *config.Config) gateway.KvGateway {
	log.Infof("NewKvGateway begin...")

	cli := client.NewClient(conf.ServerConfig.Address, client.WithConnSize(10))

	kvClient := kv.NewStoreClient(cli)

	g := &kvGateway{
		client: kvClient,
	}

	log.Infof("NewKvGateway success...")

	return g
}
