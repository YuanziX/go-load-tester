package main

import (
	"context"
	"time"
)

type RequestError struct {
	Timestamp  time.Time
	Error      string
	StatusCode int
	Latency    time.Duration
}

type RequestResult struct {
	Success   bool          `json:"success"`
	TimeTaken time.Duration `json:"timeTaken"`
	ErrorInfo *RequestError `json:"errorInfo"`
}

type HttpResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
}

type Job struct {
	config  Config
	metrics Metrics
	ctx     context.Context
	cancel  context.CancelFunc
	done    chan struct{}
}

type Server struct {
	jobs map[string]*Job
}
