package gateway

import (
	"context"
	"github.com/poonman/entry-task/client/domain/aggr/user"
	"github.com/poonman/entry-task/client/domain/gateway"
	"github.com/poonman/entry-task/client/idl/kv"
	"github.com/poonman/entry-task/client/infra/config"
	"github.com/poonman/entry-task/dora/client"
	"github.com/poonman/entry-task/dora/metadata"
	"github.com/poonman/entry-task/dora/misc/log"
	"github.com/poonman/entry-task/dora/status"
)

type kvGateway struct {
	client *client.Client
	kvClient kv.StoreClient
}

func Status2Error(st *kv.Status) (err error) {
	if st.Code == 0 {
		return nil
	}

	return status.New(status.Code(st.Code), st.Message)
}

func (g *kvGateway) Login(u *user.User) (err error) {
	req := &kv.LoginReq{
		Username: u.Name,
		Password: u.Password,
	}

	rsp, err := g.kvClient.Login(context.TODO(), req)
	if err != nil {
		log.Errorf("Failed to login. username:[%s], err:[%s]", u.Name, err)
		return
	}

	if rsp.Status.Code != kv.CODE_OK {
		err = Status2Error(rsp.Status)
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

	_, err = g.kvClient.WriteSecureMessage(ctx, req)
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

	log.Debugf("get begin...")
	rsp, err := g.kvClient.ReadSecureMessage(ctx, req)
	if err != nil {
		log.Errorf("Failed to read secure message. user:[%+v], err:[%v]",
			u, err)
		return
	}

	log.Debugf("get rsp:%+v\n", rsp)
	if rsp.Status.Code != kv.CODE_OK {
		err = Status2Error(rsp.Status)
		return
	}

	value = rsp.Value

	return
}

func (g *kvGateway) Stop() {
	g.client.Stop()
}

func NewKvGateway(conf *config.Config) gateway.KvGateway {

	var cli *client.Client

	if conf.ServerConfig.EnableTls {
		cli = client.NewClient(conf.ServerConfig.Address, client.WithConnSize(conf.ServerConfig.MaxActiveConn),
			client.WithTlsConfig(conf.LoadTlsConfig()))
	} else {
		cli = client.NewClient(conf.ServerConfig.Address, client.WithConnSize(conf.ServerConfig.MaxActiveConn))
	}

	kvClient := kv.NewStoreClient(cli)

	g := &kvGateway{
		client: cli,
		kvClient: kvClient,
	}

	return g
}
