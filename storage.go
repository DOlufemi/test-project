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
	timer    *time.Timer
	maxtime  time.Duration
	maxbatch int
	curbatch int
	mlock    sync.Mutex
	dlock    sync.Mutex
	aiochan  chan bool
}

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
		aiochan:  make(chan bool, 10),
	}

	s.timer = time.NewTimer(maxtime)
	s.timer.Stop()
	go s.flushOnTimeout()

	return s, nil
}

// Reload reopens storage file
func (s *Storage) Reload() error {
	s.mlock.Lock()
	defer s.mlock.Unlock()
	s.Flush(false)
	fd, err := os.OpenFile(s.fd.Name(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	s.fd = fd
	return err
}

func (s *Storage) Write(m *Metric) error {
	s.mlock.Lock()
	defer s.mlock.Unlock()

	if s.buf.Len() == 0 {
		s.timer.Reset(s.maxtime)
	}

	s.buf.WriteString(m.String())
	s.buf.WriteString("\n")

	s.curbatch++
	if s.curbatch >= s.maxbatch {
		s.Flush(false)
	}
	return nil
}

// Flush flushes buffered data into storage file
func (s *Storage) Flush(lock bool) {
	if lock {
		s.mlock.Lock()
		defer s.mlock.Unlock()
	}
	s.timer.Stop()
	s.curbatch = 0

	oldbuf := s.buf
	s.buf = bufPool.Get().(*bytes.Buffer)

	s.aiochan <- true
	go func() {
		s.dlock.Lock()
		_, err := io.Copy(s.fd, oldbuf)
		if err != nil {
			log.Printf("[ERROR] Failed to flush buffer to storage: %s", err)
		}
		s.dlock.Unlock()
		err = s.fd.Sync()
		if err != nil {
			log.Printf("[ERROR] Failed to fsync storage: %s", err)
		}
		oldbuf.Reset()
		bufPool.Put(oldbuf)
		_ = <-s.aiochan
	}()
}

func (s *Storage) flushOnTimeout() {
	for range s.timer.C {
		s.Flush(true)
	}
}
