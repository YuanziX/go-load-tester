package main

import (
	"context"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second,
}

func requestWorker(ctx context.Context, url string, queue chan<- RequestResult) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			curr := time.Now()
			resp, err := httpClient.Get(url)
			timeTaken := time.Since(curr)
			success := true

			if err != nil {
				success = false
			} else {
				if resp == nil {
					success = false
				} else {
					_ = resp.Body.Close()
					if resp.StatusCode > 299 {
						success = false
					}
				}
			}

			queue <- RequestResult{
				success:   success,
				timeTaken: timeTaken,
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
