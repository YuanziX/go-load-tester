package main

import (
	"sync"
	"time"
)

type RequestResult struct {
	success   bool
	timeTaken time.Duration
}

type Metrics struct {
	totalRequests      int
	successfulRequests int
	failedRequests     int
	minLatency         time.Duration
	maxLatency         time.Duration
	avgLatency         time.Duration
	totalLatency       time.Duration
	mux                sync.Mutex
}
