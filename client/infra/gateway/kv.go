package gateway

import (
	"context"
	"github.com/poonman/entry-task/client/domain/aggr/user"
	"github.com/poonman/entry-task/client/domain/gateway"
	"github.com/poonman/entry-task/client/idl/kv"
	"github.com/poonman/entry-task/client/infra/config"
	"github.com/poonman/entry-task/dora/client"
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/dora/metadata"
)

type kvGateway struct {
	client kv.StoreClient
}

func (g *kvGateway) Login(u *user.User) (err error) {
	req := &kv.LoginReq{
		Username: u.Name,
		Password: u.Password,
	}

	rsp, err := g.client.Login(context.TODO(), req)
	if err != nil {
		log.Errorf("Failed to login. username:[%s], err:[%s]", u.Name, err)
		return
	}

	if rsp.Status.Code != kv.CODE_OK {
		return
	}

	u.Token = rsp.Token

	return
}

func (g *kvGateway) Set(u *user.User, key, value string) (err error) {
	req := &kv.WriteSecureMessageReq{
		Key:   key,
		Value: value,
	}

	ctx := metadata.NewOutgoingContext(context.TODO(), map[string]string{
		"username": u.Name,
		"token":    u.Token,
	})

	_, err = g.client.WriteSecureMessage(ctx, req)
	if err != nil {
		log.Errorf("Failed to write secure message. user:[%+v], err:[%v]",
			u, err)
		return
	}
	return
}

func (g *kvGateway) Get(u *user.User, key string) (value string, err error) {
	req := &kv.ReadSecureMessageReq{
		Key: key,
	}

	ctx := metadata.NewOutgoingContext(context.TODO(), map[string]string{
		"username": u.Name,
		"token":    u.Token,
	})

	rsp, err := g.client.ReadSecureMessage(ctx, req)
	if err != nil {
		log.Errorf("Failed to read secure message. user:[%+v], err:[%v]",
			u, err)
		return
	}

	value = rsp.Value

	return
}

func NewKvGateway(conf *config.Config) gateway.KvGateway {
	log.Infof("NewKvGateway begin...")

	cli := client.NewClient(conf.ServerConfig.Address, client.WithConnSize(conf.ServerConfig.MaxActiveConn))

	kvClient := kv.NewStoreClient(cli)

	g := &kvGateway{
		client: kvClient,
	}

	log.Infof("NewKvGateway success...")

	return g
}
