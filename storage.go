package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// Storage struct wraps storage file for recieved metrics
type Storage struct {
	fd       *os.File
	buf      *bytes.Buffer
	bufwg    sync.WaitGroup
	fchan    chan *bytes.Buffer
	timer    *time.Timer
	maxtime  time.Duration
	maxbatch int
	curbatch int
	lock     sync.Mutex
}

const maxAsyncWrites = 10

var (
	bufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(nil)
		},
	}
)

// NewStorage creates new Storage object
func NewStorage(path string, maxbatch int, maxtime time.Duration) (*Storage, error) {
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		fd:       fd,
		buf:      bufPool.Get().(*bytes.Buffer),
		maxtime:  maxtime,
		maxbatch: maxbatch,
		fchan:    make(chan *bytes.Buffer, maxAsyncWrites),
	}

	s.timer = time.NewTimer(maxtime)
	s.timer.Stop()
	go s.flusher()
	go s.flushOnTimeout()

	return s, nil
}

// Reload reopens storage file
func (s *Storage) Reload() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Flush(false)
	s.Wait()
	fd, err := os.OpenFile(s.fd.Name(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	s.fd = fd
	return err
}

func (s *Storage) Write(m *Metric) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.buf.Len() == 0 {
		s.timer.Reset(s.maxtime)
	}

	s.buf.WriteString(m.String())
	s.buf.WriteString("\n")

	s.curbatch++
	if s.curbatch >= s.maxbatch {
		s.Flush(false)
	}
}

// Flush flushes buffered data into storage file
func (s *Storage) Flush(lock bool) {
	if lock {
		s.lock.Lock()
		defer s.lock.Unlock()
	}
	s.timer.Stop()
	s.curbatch = 0

	s.bufwg.Add(1)
	s.fchan <- s.buf

	s.buf = bufPool.Get().(*bytes.Buffer)
}

func (s *Storage) Wait() {
	s.bufwg.Wait()
}

func (s *Storage) flushOnTimeout() {
	for range s.timer.C {
		s.Flush(true)
	}
}

func (s *Storage) flusher() {
	for buf := range s.fchan {
		_, err := io.Copy(s.fd, buf)
		if err != nil {
			log.Fatalf("[ERROR] Failed to flush buffer to storage: %s", err)
		}
		err = s.fd.Sync()
		if err != nil {
			log.Fatalf("[ERROR] Failed to fsync storage: %s", err)
		}
		buf.Reset()
		bufPool.Put(buf)
		s.bufwg.Done()
	}
}
