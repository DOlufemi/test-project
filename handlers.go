package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func handleAPIv1Write(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if debug {
		log.Printf("got metric from %s", r.RemoteAddr)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if debug {
			log.Printf("Failed to read body: %s", err)
		}
		return
	}

	m, err := ParseMetric(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if debug {
			log.Printf("Failed to parse data: %s", err)
		}
		return
	}
	defer m.Free()

	err = storage.Write(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to save metric: %s", err)
		return
	}
}
