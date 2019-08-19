package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func handleAPIv1Write(storage *Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		if stopping {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		if debug {
			log.Printf("[DEBUG] Got metric from %s", r.RemoteAddr)
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if debug {
				log.Printf("[DEBUG] Failed to read body: %s", err)
			}
			return
		}

		m, err := ParseMetric(data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if debug {
				log.Printf("[DEBUG] Failed to parse data: %s", err)
			}
			return
		}
		defer m.Free()

		storage.Write(m)
	}
}
