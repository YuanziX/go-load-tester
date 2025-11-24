package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: DefaultTimeout,
}

func requestWorker(ctx context.Context, url string, rps int, queue chan<- RequestResult) {
	for range rps {
		select {
		case <-ctx.Done():
			return
		default:
			curr := time.Now()
			resp, err := httpClient.Get(url)
			timeTaken := time.Since(curr)
			success := true
			var errorInfo *RequestError

			if err != nil {
				success = false
				errorInfo = &RequestError{
					Timestamp:  time.Now(),
					Error:      err.Error(),
					StatusCode: 0,
					Latency:    timeTaken,
				}
			} else {
				if resp == nil {
					success = false
					errorInfo = &RequestError{
						Timestamp:  time.Now(),
						Error:      "nil response",
						StatusCode: 0,
						Latency:    timeTaken,
					}
				} else {
					_ = resp.Body.Close()
					if resp.StatusCode > 299 {
						success = false
						errorInfo = &RequestError{
							Timestamp:  time.Now(),
							Error:      fmt.Sprintf("HTTP %d", resp.StatusCode),
							StatusCode: resp.StatusCode,
							Latency:    timeTaken,
						}
					}
				}
			}

			queue <- RequestResult{
				success:   success,
				timeTaken: timeTaken,
				errorInfo: errorInfo,
			}
		}
	}
}

func (m *Metrics) writeWorker(ctx context.Context, queue <-chan RequestResult) {
	for {
		select {
		case <-ctx.Done():
			return
		case reqRes, ok := <-queue:
			if !ok {
				return
			}
			m.update(reqRes)
		}
	}
}
