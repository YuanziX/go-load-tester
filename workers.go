package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: DefaultTimeout,
	Transport: &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 500,
		MaxConnsPerHost:     500,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	},
}

func requestWorker(ctx context.Context, url string, rpw int, queue chan<- RequestResult) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range rpw {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
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
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
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
				Success:   success,
				TimeTaken: timeTaken,
				ErrorInfo: errorInfo,
			}
		}
	}
}

func (m *Metrics) writeWorker(ctx context.Context, queue <-chan RequestResult) {
	list := make([]RequestResult, 0, 520)

	for {
		select {
		case <-ctx.Done():
			if len(list) > 0 {
				m.update(list)
				list = list[:0]
			}
			return
		case reqRes, ok := <-queue:
			if !ok {
				if len(list) > 0 {
					m.update(list)
				}
				return
			}
			list = append(list, reqRes)

			if len(list) >= 500 {
				m.update(list)
				list = list[:0]
			}
		}
	}
}
