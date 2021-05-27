package domain

import (
	"fmt"
	"github.com/poonman/entry-task/client/domain/gateway"
	"github.com/poonman/entry-task/client/infra/config"
)

type Service struct {
	kvGateway gateway.KvGateway
	conf *config.Config

	keys []string
	values []string
}

func NewService(
	kvGateway gateway.KvGateway,
	conf *config.Config,
	) *Service {

	s := &Service{
		kvGateway: kvGateway,
		conf: conf,
	}

	keys := make([]string, 0, 100)
	values := make([]string, 0, 100)

	for i:=1;i<=100;i++ {
		key := newKey(i)
		keys = append(keys, key)
		value := newValue(i)
		values = append(values, value)
	}

	s.keys = keys
	s.values = values

	return s
}

func newKey(id int) (key string) {
	key = fmt.Sprintf("%dxxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx",id)
	for i :=0; i<9;i++ {
		key+= "0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx0xxxxxxxxx"
	}
	
	return 
}

func newValue(id int) (key string) {
	key = fmt.Sprintf("%dyyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy", id)
	for i :=0; i<9;i++ {
		key+= "0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy"
	}

	return
}