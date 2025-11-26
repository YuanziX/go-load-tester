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
	success   bool
	timeTaken time.Duration
	errorInfo *RequestError
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
