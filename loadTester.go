package main

import (
	"context"
	"sync"
)

func setupLoadTesterWorkers(config *Config, metrics *Metrics, ctx context.Context) {
	queue := make(chan RequestResult, 10_000)
	var workersWg sync.WaitGroup
	var writerWg sync.WaitGroup

	workersWg.Add(config.requestWorkersCount)
	for range config.requestWorkersCount {
		go func() {
			defer workersWg.Done()
			requestWorker(ctx, url, config.requestsPerWorker, queue)
		}()
	}

	writerWg.Add(config.writerWorkersCount)
	for range config.writerWorkersCount {
		go func() {
			defer writerWg.Done()
			metrics.writeWorker(ctx, queue)
		}()
	}

	workersWg.Wait()
	close(queue)
	writerWg.Wait()
}
