package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	debug    bool
	spath    string
	maxbatch int
	maxtime  int
	addr     = "[::1]:8080"
)

func main() {
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.StringVar(&spath, "s", "storage.txt", "Path to storage file")
	flag.IntVar(&maxbatch, "b", 1000, "Max size of batch written to storage")
	flag.IntVar(&maxtime, "t", 5, "Max time before data is flushed from buffer")
	flag.Parse()

	if flag.NArg() > 0 {
		addr = flag.Arg(0)
	}

	storage, err := NewStorage(spath, maxbatch, time.Duration(maxtime)*time.Second)
	if err != nil {
		log.Fatalf("Failed to open storage: %s\n", err)
	}

	sigchan := make(chan os.Signal, 1)
	go handleSig(sigchan, storage)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	http.HandleFunc("/api/v1/write", handleAPIv1Write(storage))

	log.Printf("Starting server on %s\n", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Failed to listen http socket: %s\n", err)
	}
}

func handleSig(sigchan chan os.Signal, storage *Storage) {
	for sig := range sigchan {
		switch sig {
		case syscall.SIGHUP:
			log.Println("Reopening storage")
			storage.Reload()
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Println("Flushing storage")
			storage.Flush(true)
			os.Exit(0)
		}
	}
}
