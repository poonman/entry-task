package stat

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type Report struct {
	Concurrency    int            `json:"Concurrency"`
	Success        int            `json:"Success"`
	Failure        int            `json:"Failure"`
	QPS            float32        `json:"QPS"`
	MaxRT          time.Duration  `json:"MaxRT"`
	MinRT          time.Duration  `json:"MinRT"`
	AvgRT          time.Duration  `json:"AvgRT"`
	LatencySummary LatencySummary `json:"LatencySummary"`
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
	return fmt.Sprintf("%s : %d milliseconds (cumulative count %d)", l.Percentile, l.AvgDelay.Milliseconds(), l.Count)
}

func getIndex(total, i int) int {
	return int(float32(total*i) / float32(100))
}

func (r *Report) Statistic(stats []*Stat) {
	sort.Sort(Stats(stats))

	total := len(stats)
	r.MaxRT = 0
	r.MinRT = math.MaxInt64

	var totalRT time.Duration

	p50Index := getIndex(total, 50)
	p60Index := getIndex(total, 60)
	p70Index := getIndex(total, 70)
	p80Index := getIndex(total, 80)
	p90Index := getIndex(total, 90)
	p95Index := getIndex(total, 95)
	p99Index := getIndex(total, 99)
	p100Index := total

	//log.Debugf("p50Index:%d", p50Index)
	//log.Debugf("p60Index:%d", p60Index)
	//log.Debugf("p70Index:%d", p70Index)
	//log.Debugf("p80Index:%d", p80Index)
	//log.Debugf("p90Index:%d", p90Index)
	//log.Debugf("p95Index:%d", p95Index)
	//log.Debugf("p99Index:%d", p99Index)

	for k, st := range stats {

		if st.Success {
			r.Success++
		} else {
			r.Failure++
		}

		if r.MaxRT < st.RT {
			r.MaxRT = st.RT
		}

		if r.MinRT > st.RT {
			r.MinRT = st.RT
		}

		totalRT += st.RT

		if k+1 == p100Index {
			lat := &Latency{
				Percentile: "100%",
				AvgDelay:   totalRT / time.Duration(k),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p99Index {
			lat := &Latency{
				Percentile: "99%",
				AvgDelay:   totalRT / time.Duration(k),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p95Index {
			lat := &Latency{
				Percentile: "95%",
				AvgDelay:   totalRT / time.Duration(k),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p90Index {
			lat := &Latency{
				Percentile: "90%",
				AvgDelay:   totalRT / time.Duration(k),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p80Index {
			lat := &Latency{
				Percentile: "80%",
				AvgDelay:   totalRT / time.Duration(k),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p70Index {
			lat := &Latency{
				Percentile: "70%",
				AvgDelay:   totalRT / time.Duration(k),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p60Index {
			lat := &Latency{
				Percentile: "60%",
				AvgDelay:   totalRT / time.Duration(k),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		} else if k+1 == p50Index {
			lat := &Latency{
				Percentile: "50%",
				AvgDelay:   totalRT / time.Duration(k),
				Count:      k + 1,
			}

			r.LatencySummary.Latencies = append(r.LatencySummary.Latencies, lat.String())
		}
	}

	r.QPS = float32(r.Concurrency) * float32(r.Success) / (float32(totalRT) / float32(time.Second))
	r.AvgRT = totalRT / time.Duration(total)

}

func (r *Report) String() string {
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
