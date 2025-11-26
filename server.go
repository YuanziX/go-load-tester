package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func (s *Server) getMetrics(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	job, ok := s.jobs[id]
	if !ok {
		json.NewEncoder(w).Encode(HttpResponse{
			success: false,
			data:    "No such job exists in the system",
		})
	}

	m := &job.metrics

	w.Header().Set("Content-Type", "application/json")
	m.Mux.RLock()
	defer m.Mux.RUnlock()
	json.NewEncoder(w).Encode(&m)
}

func (s *Server) startWorking(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Url               string `json:"url"`
		RequestsPerSecond int    `json:"rps"`
		RequestsPerWorker int    `json:"rpw"`
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

	config := getConfigWithParams(
		params.Url,
		params.RequestsPerSecond,
		params.RequestsPerWorker,
	)

	ctx, cancel := context.WithCancel(context.Background())
	key := GenerateID(16)
	job := Job{
		config:  config,
		metrics: getMetricsObject(),
		ctx:     ctx,
		cancel:  cancel,
	}
	s.jobs[key] = &job

	setupLoadTesterWorkers(&job.config, &job.metrics, job.ctx)

	json.NewEncoder(w).Encode(struct {
		Id      string
		Metrics *Metrics
	}{
		Id:      key,
		Metrics: &job.metrics,
	})
}

func serveUI() {
	mux := http.NewServeMux()

	// setup server with in memory user sessions
	server := Server{
		jobs: make(map[string]*Job),
	}

	mux.HandleFunc("GET /metrics", server.getMetrics)
	mux.HandleFunc("POST /metrics", server.startWorking)

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
