package app

import (
	"fmt"
	"github.com/poonman/entry-task/client/domain"
	"github.com/poonman/entry-task/client/domain/aggr/benchmark"
	"github.com/poonman/entry-task/client/domain/aggr/user"
	"github.com/poonman/entry-task/client/domain/aggr/userm"
	"github.com/poonman/entry-task/dora/misc/log"
	"strconv"
)

type Service struct {
	domainSvc *domain.Service

	userManager *userm.UserManager
}

func NewService(domainSvc *domain.Service) *Service {
	return &Service{
		domainSvc:   domainSvc,
		userManager: userm.NewUserManager(),
	}
}

func (s *Service) Login(username, password string) (err error) {
	u := user.NewUser(username, password)

	err = s.domainSvc.Login(u)
	if err != nil {
		log.Errorf("Failed to login. username:[%s], err:[%v]\n\n", username, err)
		return
	}

	s.userManager.AddUser(u)

	log.Infof("login success. username:[%s]\n\n", username)
	return
}

func (s *Service) WriteSecureMessage(username, password string, k, v byte) (err error) {

	u := s.userManager.GetUser(username)
	if u == nil {
		u = user.NewUser(username, password)
		s.userManager.AddUser(u)
	}

	err = s.domainSvc.SetKV(u, k, v)
	if err != nil {
		log.Errorf("failed to set kv. err:[%v]\n\n", err)
		return
	}

	log.Infof("write secure message success. username:[%s], k:[%c], v:[%c]\n\n", username, k, v)
	return
}

func (s *Service) ReadSecureMessage(username, password string, k byte) (err error) {
	u := s.userManager.GetUser(username)

	if u == nil {
		u = user.NewUser(username, password)
		s.userManager.AddUser(u)
	}

	var vv string
	var v byte
	vv, err = s.domainSvc.GetKV(u, k)
	if err != nil {
		log.Errorf("failed to get kv. user:[%+v], err:[%v]\n\n", u, err)
		return
	}

	if len(vv) > 0 {
		v = vv[0]
	}

	log.Infof("read secure message success. username:[%s], k:[%c], v:[%c]\n\n", username, k, v)

	return
}

func (s *Service) Benchmark(concurrency, requests int, username, password, method string, k, v byte) {

	if concurrency <= 0{
		concurrency = 10
	}

	if requests <= 0 {
		requests = 10
	}

	if method == "" {
		method = "read"
	}

	if k < 'a' || k > 'z' {
		k = 'a'
	}

	if v < 'a' || v > 'z' {
		k = 'a'
	}

	bm := benchmark.NewBenchmark(concurrency, requests, username, method, k, v)

	var users []*user.User

	_, err := strconv.ParseInt(username, 10, 64)
	if err != nil {

		users = batchNewUser(concurrency)

	}else {
		u := user.NewUser(username, password)

		for i:=0; i<concurrency;i++ {
			users = append(users, u)
		}
	}

	s.domainSvc.Benchmark(bm, users)
}

func batchNewUser(num int) (users []*user.User) {

	users = make([]*user.User, 0, num)

	for i:=1;i<=num;i++ {
		username := fmt.Sprintf("%d", 100000+i)

		u := user.NewUser(username, username)

		users = append(users, u)
	}

	return
}