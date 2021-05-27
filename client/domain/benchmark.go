package domain

import (
	"github.com/poonman/entry-task/client/domain/aggr/stat"
	"github.com/poonman/entry-task/client/domain/aggr/user"
	"github.com/poonman/entry-task/dora/log"
	"strconv"
	"sync"
	"time"
)

func (s *Service) BenchmarkRead() {

	var (
		wg sync.WaitGroup
	)

	wg.Add(s.conf.BenchmarkConfig.Concurrency)

	stats := make([]*stat.Stat, s.conf.BenchmarkConfig.Concurrency*s.conf.BenchmarkConfig.RequestNumPerConcurrency)

	u := &user.User{
		Name:     "1",
		Password: "1",
		Token:    "",
	}

	err := s.Login(u)
	if err != nil {
		log.Errorf("failed to login. err:[%v]", err)
		return
	}

	for i := 0; i < s.conf.BenchmarkConfig.Concurrency; i++ {

		go func(no int, u *user.User) {
			var tmp []*stat.Stat

			if no == s.conf.BenchmarkConfig.Concurrency {
				tmp = stats[no*s.conf.BenchmarkConfig.Concurrency:]
			} else {
				tmp = stats[no*s.conf.BenchmarkConfig.RequestNumPerConcurrency : (no+1)*s.conf.BenchmarkConfig.RequestNumPerConcurrency]
			}
			s.RequestRead(u, no, tmp)
			wg.Done()
		}(i, u)
	}

	wg.Wait()

	rep := &stat.Report{
		Concurrency: s.conf.BenchmarkConfig.Concurrency,
	}

	rep.Statistic(stats)

	log.Infof("report:[%s]", rep)
}

func (s *Service) RequestRead(u *user.User, concurrencyNo int, stats []*stat.Stat) {

	for i := 1; i <= s.conf.BenchmarkConfig.RequestNumPerConcurrency; i++ {
		if u == nil {
			tmp := strconv.Itoa(concurrencyNo*1000 + i)
			u = &user.User{
				Name:     tmp,
				Password: tmp,
				Token:    "",
			}

			_ = s.Login(u)
		}

		st := &stat.Stat{}
		stats[i-1] = st

		before := time.Now()
		_, err := s.kvGateway.Get(u, s.keys[0])
		rt := time.Now().Sub(before)
		if err != nil {
			st.Success = false
		} else {
			st.Success = true
		}

		st.RT = rt
	}
}
