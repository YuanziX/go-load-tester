package main

import (
	"context"
	"sync"
	"time"
)

type RequestError struct {
	Timestamp  time.Time
	Error      string
	StatusCode int
	Latency    time.Duration
}

type RequestResult struct {
	success   bool
	timeTaken time.Duration
	errorInfo *RequestError
}

type Metrics struct {
	TotalRequests      int
	SuccessfulRequests int
	FailedRequests     int
	MinLatency         time.Duration
	MaxLatency         time.Duration
	AvgLatency         time.Duration
	TotalLatency       time.Duration
	Errors             []RequestError

	IsCompleted bool
	Mux         sync.RWMutex
}

func getMetricsObject() (metrics Metrics) {
	metrics = Metrics{MinLatency: MaxDuration}
	return
}

type HttpResponse struct {
	success bool
	data    string
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
