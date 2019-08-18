package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testData1 = `{"ts": 0, "key": "metric1", "val": 10}`
)

func TestParseMetric(t *testing.T) {
	m, err := ParseMetric([]byte(testData1))
	assert.Nil(t, err)

	assert.Equal(t, int64(0), m.TS)
	assert.Equal(t, "metric1", m.Key)
	assert.Equal(t, int64(10), m.Value)
}

func TestFormat(t *testing.T) {
	m, err := ParseMetric([]byte(`{"ts": 0, "key": "metric1", "val": 10}`))
	assert.Nil(t, err)

	assert.Equal(t, m.String(), "0 metric1 10")
}
