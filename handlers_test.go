package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	expect200 = `{"ts": 10000, "key": "metric1", "val": 10}`
	expect400 = `{"ts": 10000, "key": "metric1", "val": 10`
)

func TestAPIV1Write200(t *testing.T) {
	spath := filepath.Join(os.TempDir(), storName)
	_ = os.Remove(spath)
	storage, err := NewStorage(spath, 1, 1*time.Second)
	assert.Nil(t, err)
	defer os.Remove(spath)
	defer storage.Flush(true)

	ts := httptest.NewServer(http.HandlerFunc(handleAPIv1Write(storage)))
	defer ts.Close()

	res200, err := http.Post(ts.URL, "application/json", bytes.NewBufferString(expect200))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res200.StatusCode)
}

func TestAPIV1Write400(t *testing.T) {
	spath := filepath.Join(os.TempDir(), storName)
	_ = os.Remove(spath)
	storage, err := NewStorage(spath, 1, 1*time.Second)
	assert.Nil(t, err)
	defer os.Remove(spath)
	defer storage.Wait()
	defer storage.Flush(true)

	ts := httptest.NewServer(http.HandlerFunc(handleAPIv1Write(storage)))
	defer ts.Close()

	res400, err := http.Post(ts.URL, "application/json", bytes.NewBufferString(expect400))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res400.StatusCode)
}

func BenchmarkAPIv1Write(b *testing.B) {
	spath := filepath.Join(os.TempDir(), storName)
	_ = os.Remove(spath)
	storage, err := NewStorage(spath, 1000, 5*time.Second)
	assert.Nil(b, err)
	defer os.Remove(spath)
	defer storage.Flush(true)

	ts := httptest.NewServer(http.HandlerFunc(handleAPIv1Write(storage)))
	defer ts.Close()

	defTrans := http.DefaultTransport.(*http.Transport)
	defTrans.MaxIdleConnsPerHost = 100

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			res200, err := http.Post(ts.URL, "application/json", bytes.NewBufferString(expect200))
			if err != nil {
				b.Fatal(err)
			}
			if res200.StatusCode != http.StatusOK {
				b.Fatal("Response status is not 200 OK")
			}
		}
	})
}
