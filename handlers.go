package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) getMetrics(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	job, ok := s.jobs[id]
	if !ok {
		json.NewEncoder(w).Encode(HttpResponse{
			Success: false,
			Data:    "No such job exists in the system",
		})
		return
	}

	m := &job.metrics

	w.Header().Set("Content-Type", "application/json")
	m.mux.RLock()
	defer m.mux.RUnlock()
	json.NewEncoder(w).Encode(&m)
}

func (s *Server) createJob(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Url               string `json:"url"`
		RequestsPerSecond int    `json:"rps"`
		RequestsPerWorker int    `json:"rpw"`
	}

	decoder := json.NewDecoder(r.Body)
	params := payload{}
	if err := decoder.Decode(&params); err != nil {
		json.NewEncoder(w).Encode(HttpResponse{
			Success: false,
			Data:    "Failed to decode payload",
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
		done:    make(chan struct{}),
	}
	s.jobs[key] = &job

	setupLoadTesterWorkers(&job.config, &job.metrics, job.ctx, job.done)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Id      string   `json:"id"`
		Metrics *Metrics `json:"metrics"`
	}{
		Id:      key,
		Metrics: &job.metrics,
	})
}

func (s *Server) cancelJob(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	job, ok := s.jobs[id]
	if !ok {
		json.NewEncoder(w).Encode(HttpResponse{
			Success: false,
			Data:    "No such job exists in the system",
		})
		return
	}

	job.cancel()
	<-job.done // wait for workers to stop

	w.Header().Set("Content-Type", "application/json")
	job.metrics.mux.RLock()
	defer job.metrics.mux.RUnlock()
	json.NewEncoder(w).Encode(&job.metrics)
}
