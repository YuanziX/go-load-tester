package main

import (
	"context"
	"sync"
)

func setupLoadTesterWorkers(config *Config, metrics *Metrics, ctx context.Context) {
	queue := make(chan RequestResult, config.queueChannelSize)
	var workersWg sync.WaitGroup
	var writerWg sync.WaitGroup

	workersWg.Add(config.requestWorkersCount)
	for range config.requestWorkersCount {
		go func() {
			defer workersWg.Done()
			requestWorker(ctx, config.url, config.requestsPerWorker, queue)
		}()
	}

	writerWg.Add(config.writerWorkersCount)
	for range config.writerWorkersCount {
		go func() {
			defer writerWg.Done()
			metrics.writeWorker(ctx, queue)
		}()
	}

	go func() {
		workersWg.Wait()
		close(queue)
		writerWg.Wait()

		metrics.Mux.Lock()
		defer metrics.Mux.Unlock()
		metrics.IsCompleted = true
	}()
}
