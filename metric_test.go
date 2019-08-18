package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testData1 = `{"ts": 0, "key": "metric1", "val": 10}`
)

func TestParseMetric(t *testing.T) {
	m, err := ParseMetric([]byte(testData1))
	assert.Nil(t, err)

	assert.Equal(t, m.TS, int64(0))
	assert.Equal(t, m.Key, "metric1")
	assert.Equal(t, m.Value, int64(10))
}

func TestString(t *testing.T) {
	m, err := ParseMetric([]byte(`{"ts": 0, "key": "metric1", "val": 10}`))
	assert.Nil(t, err)

	assert.Equal(t, m.String(), "0 metric1 10")
}
