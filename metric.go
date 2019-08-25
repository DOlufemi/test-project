package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Metric struct contains metric data
type Metric struct {
	TS    int64  `json:"ts"`
	Key   string `json:"key"`
	Value int64  `json:"val"`
}

var (
	metricPool = sync.Pool{
		New: func() interface{} {
			return &Metric{}
		},
	}
)

func (m *Metric) String() string {
	return fmt.Sprintf("%d %s %d", m.TS, m.Key, m.Value)
}

// ParseMetric creates Metric object from serialized data
func ParseMetric(data []byte) (*Metric, error) {
	result := metricPool.Get().(*Metric)

	err := json.Unmarshal(data, &result)
	if err != nil {
		result.Free()
		return nil, err
	}
	if len(result.Key) == 0 {
		result.Free()
		return nil, fmt.Errorf("empty metric key")
	}

	return result, nil
}

// Free releases Metric struct to pool
func (m *Metric) Free() {
	m.TS = 0
	m.Key = ""
	m.Value = 0
	metricPool.Put(m)
}
