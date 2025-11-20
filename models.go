package main

import (
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
	totalRequests      int
	successfulRequests int
	failedRequests     int
	minLatency         time.Duration
	maxLatency         time.Duration
	avgLatency         time.Duration
	totalLatency       time.Duration
	errors             []RequestError
	mux                sync.Mutex
}
