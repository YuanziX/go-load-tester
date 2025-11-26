package main

import (
	"fmt"
	"os"
	"time"
)

func (m *Metrics) update(reqResults []RequestResult) {
	m.Mux.Lock()
	defer m.Mux.Unlock()

	for _, reqRes := range reqResults {
		m.TotalRequests++

		if reqRes.success {
			m.SuccessfulRequests++
		} else {
			m.FailedRequests++
			if reqRes.errorInfo != nil {
				m.Errors = append(m.Errors, *reqRes.errorInfo)
			}
		}

		if reqRes.timeTaken < m.MinLatency {
			m.MinLatency = reqRes.timeTaken
		}

		if reqRes.timeTaken > m.MaxLatency {
			m.MaxLatency = reqRes.timeTaken
		}

		m.TotalLatency += reqRes.timeTaken
		m.AvgLatency = m.TotalLatency / time.Duration(m.TotalRequests)
	}
}

func (m *Metrics) WriteErrorsToFile(filename string) error {
	m.Mux.RLock()
	defer m.Mux.RUnlock()

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
	m.Mux.Lock()
	defer m.Mux.Unlock()

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
