package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const storName = "mstor-test.txt"

func TestStorageFlush(t *testing.T) {
	spath := filepath.Join(os.TempDir(), storName)
	_ = os.Remove(spath)
	s, err := NewStorage(spath, 5, 5*time.Second)
	assert.Nil(t, err)
	defer os.Remove(spath)

	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})
	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})
	s.Flush(true)
	s.Wait()

	buf, err := ioutil.ReadFile(spath)
	assert.Nil(t, err)
	count := bytes.Count(buf, []byte{'\n'})
	assert.Equal(t, 2, count)
}

func TestStorageBatch(t *testing.T) {
	spath := filepath.Join(os.TempDir(), storName)
	_ = os.Remove(spath)
	s, err := NewStorage(spath, 5, 5*time.Second)
	assert.Nil(t, err)
	defer os.Remove(spath)

	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})
	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})
	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})
	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})
	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})
	s.Wait()

	buf, err := ioutil.ReadFile(spath)
	assert.Nil(t, err)
	count := bytes.Count(buf, []byte{'\n'})
	assert.Equal(t, 5, count)
}

func TestStorageTime(t *testing.T) {
	spath := filepath.Join(os.TempDir(), storName)
	_ = os.Remove(spath)
	s, err := NewStorage(spath, 100, 500*time.Millisecond)
	assert.Nil(t, err)
	defer os.Remove(spath)

	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})
	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})
	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})
	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})
	s.Write(&Metric{TS: 0, Key: "metric1", Value: 10})

	time.Sleep(700 * time.Millisecond)
	s.Wait()

	buf, err := ioutil.ReadFile(spath)
	assert.Nil(t, err)
	count := bytes.Count(buf, []byte{'\n'})
	assert.Equal(t, 5, count)
}
