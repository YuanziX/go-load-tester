package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func (m *Metrics) getMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	m.Mux.RLock()
	defer m.Mux.RUnlock()
	json.NewEncoder(w).Encode(m)
}

func (m *Metrics) startWorking(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Url                 string `json:"url"`
		RequestWorkersCount int    `json:"rwc"`
		RequestsPerWorker   int    `json:"rpw"`
	}

	decoder := json.NewDecoder(r.Body)
	params := payload{}
	if err := decoder.Decode(&params); err != nil {
		json.NewEncoder(w).Encode(HttpResponse{
			success: false,
			data:    "Failed to decode payload",
		})
		log.Println(err)
		return
	}

	m.Reset()
	config := getConfigWithParams(
		params.Url,
		params.RequestWorkersCount,
		params.RequestsPerWorker,
	)

	setupLoadTesterWorkers(&config, m, context.Background())

	json.NewEncoder(w).Encode(m)
}

func serveUI(metrics *Metrics) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /metrics", metrics.getMetrics)
	mux.HandleFunc("POST /metrics", metrics.startWorking)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./static/ui.html")
	})

	handler := cors.AllowAll().Handler(mux)

	log.Println("Serving static page at '/'")
	log.Panicln(http.ListenAndServe(":8000", handler))
}
