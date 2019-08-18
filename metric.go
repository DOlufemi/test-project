package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Metric struct{
	TS int64 `json:"ts"`
	Key string `json:"key"`
	Value int64 `json:"val"`
}

var (
	metricPool = sync.Pool{
		New: func() interface{} {
			return Metric{}
		},
	}
)

func (m *Metric) String() string {
	return fmt.Sprintf("%d %s %d", m.TS, m.Key, m.Value)
}

func ParseMetric(data []byte) (*Metric, error) {
	result := metricPool.Get().(Metric)

	err := json.Unmarshal(data, &result)
	if err != nil {
		result.Free()
		return nil, err
	}
	if len(result.Key) == 0 {
		result.Free()
		return nil, fmt.Errorf("Empty metric key")
	}

	return &result, nil
}

func (m Metric) Free() {
	metricPool.Put(m)
}
