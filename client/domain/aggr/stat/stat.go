package stat

import (
	"bytes"
	"encoding/json"
	"time"
)

type Stat struct {
	Success   int
	Failure   int
	MaxRT     time.Duration
	MinRT     time.Duration
	TotalRT   time.Duration
	SuccessRT time.Duration
	FailureRT time.Duration
}

type Report struct {
	Concurrency int           `json:"Concurrency"`
	Success     int           `json:"Success"`
	Failure     int           `json:"Failure"`
	QPS         float32       `json:"QPS"`
	MaxRT       time.Duration `json:"MaxRT"`
	MinRT       time.Duration `json:"MinRT"`
	AvgRT       time.Duration `json:"AvgRT"`
	TotalRT     time.Duration `json:"TotalRT"`
	SuccessRT   time.Duration `json:"SuccessRT"`
	FailureRT   time.Duration `json:"FailureRT"`
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