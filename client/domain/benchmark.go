package domain

import (
	"github.com/poonman/entry-task/client/domain/aggr/benchmark"
	"github.com/poonman/entry-task/client/domain/aggr/user"
	"github.com/poonman/entry-task/dora/misc/log"
	"strconv"
	"sync"
	"time"
)

func (s *Service) Benchmark(bm *benchmark.Benchmark, users []*user.User) {

	var (
		wg sync.WaitGroup
	)

	wg.Add(bm.Concurrency)

	for _, u := range users {
		err := s.Login(u)
		if err != nil {
			log.Errorf("failed to login. err:[%v]", err)
			return
		}
	}

	log.Info("login success.")

	startAt := time.Now()
	level := log.GetLevel()
	// set log level FATAL before benchmark
	log.SetLevel(log.FATAL)

	for i := 0; i < bm.Concurrency; i++ {
		u := users[i]
		go func(no int, u *user.User) {
			s.request(u, bm, no)
			wg.Done()
		}(i, u)
	}

	wg.Wait()

	bm.Duration.Duration = time.Since(startAt)

	log.SetLevel(level)

	bm.Statistic()

	log.Infof("Benchmark Report:[%s]", bm)
}

func (s *Service) request(u *user.User, bm *benchmark.Benchmark, no int) {

	var stats []*benchmark.Stat

	if no == bm.Concurrency {
		stats = bm.Stats[no*bm.Requests:]
	} else {
		stats = bm.Stats[no*bm.Requests : (no+1)*bm.Requests]
	}

	var err error
	for i := 1; i <= bm.Requests; i++ {
		if u == nil {
			tmp := strconv.Itoa(no*1000 + i)
			u = &user.User{
				Name:     tmp,
				Password: tmp,
				Token:    "",
			}

			_ = s.Login(u)
		}

		st := &benchmark.Stat{}
		stats[i-1] = st

		before := time.Now()
		if bm.Method == "read" {
			_, err = s.kvGateway.Get(u, bm.Key)
		} else if bm.Method == "write" {
			err = s.kvGateway.Set(u, bm.Key, bm.Value)
		}

		st.RT = time.Since(before)

		if err != nil {
			st.Success = false
		} else {
			st.Success = true
		}

	}
}
