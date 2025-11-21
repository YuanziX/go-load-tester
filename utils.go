package main

import (
	"fmt"
	"os"
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
		if reqRes.errorInfo != nil {
			m.errors = append(m.errors, *reqRes.errorInfo)
		}
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

	if m.totalRequests > 0 {
		fmt.Printf("Success Rate:        %.2f%%\n", float64(m.successfulRequests)/float64(m.totalRequests)*100)
		// Check if minLatency was actually updated (not still at max value)
		if m.minLatency != maxDuration {
			fmt.Printf("Min Latency:         %v\n", m.minLatency)
			fmt.Printf("Max Latency:         %v\n", m.maxLatency)
		} else {
			fmt.Printf("Min Latency:         N/A\n")
			fmt.Printf("Max Latency:         N/A\n")
		}
		fmt.Printf("Avg Latency:         %v\n", m.avgLatency)
	} else {
		fmt.Printf("Success Rate:        N/A\n")
		fmt.Printf("Min Latency:         N/A\n")
		fmt.Printf("Max Latency:         N/A\n")
		fmt.Printf("Avg Latency:         N/A\n")
	}
}

func (m *Metrics) WriteErrorsToFile(filename string) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if len(m.errors) == 0 {
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
	fmt.Fprintf(file, "Total Errors: %d\n\n", len(m.errors))

	for i, e := range m.errors {
		fmt.Fprintf(file, "[%d] %s\n", i+1, e.Timestamp.Format("2006-01-02 15:04:05.000"))
		fmt.Fprintf(file, "    Error: %s\n", e.Error)
		if e.StatusCode > 0 {
			fmt.Fprintf(file, "    Status Code: %d\n", e.StatusCode)
		}
		fmt.Fprintf(file, "    Latency: %v\n\n", e.Latency)
	}

	fmt.Printf("Wrote %d errors to %s\n", len(m.errors), filename)
	return nil
}
