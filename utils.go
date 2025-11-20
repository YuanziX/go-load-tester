package main

import (
	"fmt"
	"time"
)

func (m *Metrics) update(reqRes RequestResult) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.totalRequests++

	if reqRes.success {
		m.successfulRequests++
	} else {
		m.failedRequests++
	}

	if reqRes.timeTaken < m.minLatency {
		m.minLatency = reqRes.timeTaken
	}

	if reqRes.timeTaken > m.maxLatency {
		m.maxLatency = reqRes.timeTaken
	}

	m.totalLatency += reqRes.timeTaken
	m.avgLatency = m.totalLatency / time.Duration(m.totalRequests)
}

func (m *Metrics) Print() {
	m.mux.Lock()
	defer m.mux.Unlock()

	fmt.Println("Load Test Results")
	fmt.Printf("Total Requests:      %d\n", m.totalRequests)
	fmt.Printf("Successful:          %d\n", m.successfulRequests)
	fmt.Printf("Failed:              %d\n", m.failedRequests)
	fmt.Printf("Success Rate:        %.2f%%\n", float64(m.successfulRequests)/float64(m.totalRequests)*100)
	fmt.Printf("Min Latency:         %v\n", m.minLatency)
	fmt.Printf("Max Latency:         %v\n", m.maxLatency)
	fmt.Printf("Avg Latency:         %v\n", m.avgLatency)
	fmt.Printf("Requests/sec:        %.2f\n", float64(m.totalRequests)/5.0)
}
