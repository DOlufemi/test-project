package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var (
	debug bool
	spath string
	addr  = "[::1]:8080"

	storage *Storage
)

func main() {
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.StringVar(&spath, "s", "storage.txt", "Path to storage file")
	flag.Parse()

	if flag.NArg() > 0 {
		addr = flag.Arg(0)
	}

	var err error
	storage, err = NewStorage(spath, 1000, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to open storage: %s\n", err)
	}

	http.HandleFunc("/api/v1/write", handleAPIv1Write)

	log.Printf("Starting server on %s\n", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Failed to listen http socket: %s\n", err)
	}
}
