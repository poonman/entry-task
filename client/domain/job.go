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

	stats := make([]*stat.Stat, 0, s.conf.BenchmarkConfig.Concurrency)

	for i:=0; i<s.conf.BenchmarkConfig.Concurrency; i++ {
		st := &stat.Stat{}

		stats = append(stats, st)

		go func(no int) {

			s.RequestRead(no, st)
			wg.Done()
		}(i)
	}

	wg.Wait()

	rep := &stat.Report{
		Concurrency: s.conf.BenchmarkConfig.Concurrency,
		Success:     0,
		Failure:     0,
		QPS:         0,
		MaxRT:       0,
		MinRT:       1000 * time.Second,
		AvgRT:       0,
		TotalRT:     0,
		SuccessRT:   0,
		FailureRT:   0,
	}

	for _, s := range stats {
		if s.MinRT < rep.MinRT {
			rep.MinRT = s.MinRT
		}

		if s.MaxRT > rep.MaxRT {
			rep.MaxRT = s.MaxRT
		}

		rep.Success += s.Success
		rep.Failure += s.Failure

		rep.TotalRT += s.TotalRT
		rep.SuccessRT += s.SuccessRT
		rep.FailureRT += s.FailureRT
	}



	rep.QPS = float32(s.conf.BenchmarkConfig.Concurrency) * float32(rep.Success) / (float32(rep.TotalRT) / float32(time.Second))
	rep.AvgRT = rep.TotalRT / time.Duration(rep.Success+rep.Failure)

	log.Infof("report:[%s]", rep)
}

func (s *Service) RequestRead(concurrencyNo int, stat *stat.Stat) {


	for i := 1; i <= s.conf.BenchmarkConfig.RequestNumPerConcurrency; i++ {
		tmp := strconv.Itoa(concurrencyNo*1000+i)
		u := &user.User{
			Name:     tmp,
			Password: tmp,
			Token:    "",
		}
		before := time.Now()
		_, err := s.kvGateway.Get(u, s.keys[0])
		rt := time.Now().Sub(before)
		if err != nil {
			stat.Failure++
		} else {
			stat.Success++
		}

		if rt > stat.MaxRT {
			stat.MaxRT = rt
		}

		if rt < stat.MinRT {
			stat.MinRT = rt
		}

		stat.TotalRT += rt
	}
}
