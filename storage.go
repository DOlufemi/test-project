package main

import (
	"bytes"
	"io"
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
	mutex    sync.Mutex
}

// NewStorage creates new Storage object
func NewStorage(path string, maxbatch int, maxtime time.Duration) (*Storage, error) {
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		fd:       fd,
		buf:      bytes.NewBuffer(nil),
		maxtime:  maxtime,
		maxbatch: maxbatch,
	}

	s.timer = time.NewTimer(maxtime)
	s.timer.Stop()
	go s.flushOnTimeout()

	return s, nil
}

// Reload reopens storage file
func (s *Storage) Reload() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Flush(false)
	fd, err := os.OpenFile(s.fd.Name(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	s.fd = fd
	return err
}

func (s *Storage) Write(m *Metric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.buf.Len() == 0 {
		s.timer.Reset(s.maxtime)
	}

	s.buf.WriteString(m.String())
	s.buf.WriteString("\n")

	s.curbatch++
	if s.curbatch >= s.maxbatch {
		return s.Flush(false)
	}
	return nil
}

// Flush flushes buffered data into storage file
func (s *Storage) Flush(lock bool) error {
	if lock {
		s.mutex.Lock()
		defer s.mutex.Unlock()
	}
	_, err := io.Copy(s.fd, s.buf)
	s.buf.Reset()
	s.timer.Stop()
	return err
}

func (s *Storage) flushOnTimeout() {
	for range s.timer.C {
		_ = s.Flush(true)
	}

}
