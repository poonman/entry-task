package benchmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/poonman/entry-task/client/infra/helper"
	"math"
	"sort"
	"time"
)

type Stat struct {
	Success bool
	RT      time.Duration
}

type Stats []*Stat

func (s Stats) Len() int {
	return len(s)
}

func (s Stats) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Stats) Less(i, j int) bool {
	if s[i].RT <= s[j].RT {
		return true
	}

	return false
}

func NewBenchmark(concurrency, requests int, username, method string, k, v byte) *Benchmark {
	return &Benchmark{
		Concurrency:    concurrency,
		Requests:       requests,
		Success:        0,
		Failure:        0,
		QPS:            0,
		MaxRT:          helper.Duration{},
		MinRT:          helper.Duration{},
		AvgRT:          helper.Duration{},
		LatencySummary: LatencySummary{},
		Username:       username,
		Method:         method,
		K:              k,
		V:              v,
		Key:            helper.NewString(k),
		Value:          helper.NewString(v),
		Stats:          make([]*Stat, concurrency*requests),
	}
}

type Benchmark struct {
	Concurrency    int             `json:"Concurrency"`
	Requests       int             `json:"requests"`
	Success        int             `json:"Success"`
	Failure        int             `json:"Failure"`
	QPS            float32         `json:"QPS"`
	MaxRT          helper.Duration `json:"MaxRT"`
	MinRT          helper.Duration `json:"MinRT"`
	AvgRT          helper.Duration `json:"AvgRT"`
	Duration       helper.Duration `json:"Duration"`
	LatencySummary LatencySummary  `json:"LatencySummary"`

	Username string `json:"username"`
	Method   string `json:"method"`
	K        byte   `json:"key"`
	V        byte   `json:"value"`

	Key   string  `json:"-"`
	Value string  `json:"-"`
	Stats []*Stat `json:"-"`
}

type LatencySummary struct {
	Latencies []string
}

type Latency struct {
	Percentile string        `json:"-"`
	AvgDelay   time.Duration `json:"-"`
	Count      int           `json:"-"`
}

func (l *Latency) String() string {
	return fmt.Sprintf("%4s : %12s (cumulative count %5d)", l.Percentile, l.AvgDelay.String(), l.Count)
}

func getIndex(total, i int) int {
	return int(float32(total*i) / float32(100))
}

func (r *Benchmark) Statistic() {
	sort.Sort(Stats(r.Stats))

	total := len(r.Stats)
	r.MaxRT.Duration = 0
	r.MinRT.Duration = math.MaxInt64

	var totalRT time.Duration

	p50Index := getIndex(total, 50)
	p60Index := getIndex(total, 60)
	p70Index := getIndex(total, 70)
	p80Index := getIndex(total, 80)
	p90Index := getIndex(total, 90)
	p95Index := getIndex(total, 95)
	p99Index := getIndex(total, 99)
	p100Index := total

	for k, st := range r.Stats {

		if st.Success {
			r.Success++
		} else {
			r.Failure++
		}

		if r.MaxRT.Duration < st.RT {
			r.MaxRT.Duration = st.RT
		}

		if r.MinRT.Duration > st.RT {
			r.MinRT.Duration = st.RT
		}

		totalRT += st.RT

		if k+1 == p100Index {
			lat := &Latency{
				Percentile: "100%",
				AvgDelay:   totalRT / time.Duration(k+1),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p99Index {
			lat := &Latency{
				Percentile: "99%",
				AvgDelay:   totalRT / time.Duration(k+1),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p95Index {
			lat := &Latency{
				Percentile: "95%",
				AvgDelay:   totalRT / time.Duration(k+1),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p90Index {
			lat := &Latency{
				Percentile: "90%",
				AvgDelay:   totalRT / time.Duration(k+1),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p80Index {
			lat := &Latency{
				Percentile: "80%",
				AvgDelay:   totalRT / time.Duration(k+1),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p70Index {
			lat := &Latency{
				Percentile: "70%",
				AvgDelay:   totalRT / time.Duration(k+1),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p60Index {
			lat := &Latency{
				Percentile: "60%",
				AvgDelay:   totalRT / time.Duration(k+1),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p50Index {
			lat := &Latency{
				Percentile: "50%",
				AvgDelay:   totalRT / time.Duration(k+1),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		}
	}

	r.QPS = float32(r.Concurrency) * float32(r.Success+r.Failure) / (float32(totalRT) / float32(time.Second))
	r.AvgRT.Duration = totalRT / time.Duration(total)

}

func (r *Benchmark) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return ""
	}
	return out.String()
}
