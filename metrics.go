package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Metrics struct {
	TotalRequests      int            `json:"totalRequests"`
	SuccessfulRequests int            `json:"successfulRequests"`
	FailedRequests     int            `json:"failedRequests"`
	MinLatency         time.Duration  `json:"minLatency"`
	MaxLatency         time.Duration  `json:"maxLatency"`
	AvgLatency         time.Duration  `json:"avgLatency"`
	TotalLatency       time.Duration  `json:"totalLatency"`
	Errors             []RequestError `json:"errors"`

	IsCompleted bool `json:"isCompleted"`
	mux         sync.RWMutex
}

func getMetricsObject() (metrics Metrics) {
	metrics = Metrics{MinLatency: MaxDuration}
	return
}

func (m *Metrics) update(reqResults []RequestResult) {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, reqRes := range reqResults {
		m.TotalRequests++

		if reqRes.Success {
			m.SuccessfulRequests++
		} else {
			m.FailedRequests++
			if reqRes.ErrorInfo != nil {
				m.Errors = append(m.Errors, *reqRes.ErrorInfo)
			}
		}

		if reqRes.TimeTaken < m.MinLatency {
			m.MinLatency = reqRes.TimeTaken
		}

		if reqRes.TimeTaken > m.MaxLatency {
			m.MaxLatency = reqRes.TimeTaken
		}

		m.TotalLatency += reqRes.TimeTaken
		m.AvgLatency = m.TotalLatency / time.Duration(m.TotalRequests)
	}
}

func (m *Metrics) WriteErrorsToFile(filename string) error {
	m.mux.RLock()
	defer m.mux.RUnlock()

	if len(m.Errors) == 0 {
		fmt.Println("No errors to write")
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create error file: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "Load Test Error Log\n")
	fmt.Fprintf(file, "===================\n\n")
	fmt.Fprintf(file, "Total Errors: %d\n\n", len(m.Errors))

	for i, e := range m.Errors {
		fmt.Fprintf(file, "[%d] %s\n", i+1, e.Timestamp.Format("2006-01-02 15:04:05.000"))
		fmt.Fprintf(file, "    Error: %s\n", e.Error)
		if e.StatusCode > 0 {
			fmt.Fprintf(file, "    Status Code: %d\n", e.StatusCode)
		}
		fmt.Fprintf(file, "    Latency: %v\n\n", e.Latency)
	}

	fmt.Printf("Wrote %d errors to %s\n", len(m.Errors), filename)
	return nil
}

func (m *Metrics) Reset() {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.TotalRequests = 0
	m.SuccessfulRequests = 0
	m.FailedRequests = 0
	m.Errors = nil
	m.MinLatency = MaxDuration
	m.MaxLatency = 0
	m.TotalLatency = 0
	m.AvgLatency = 0
	m.IsCompleted = false
}
