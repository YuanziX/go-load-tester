package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var requestWorkerCount int = 500
var writerWorkerCount int = 50
var url string = "YOUR_URL_HERE"
var testingDuration time.Duration = 5 * time.Second

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := make(chan RequestResult, 10_000)

	var workersWg sync.WaitGroup
	var writerWg sync.WaitGroup

	metrics := Metrics{minLatency: time.Duration(1<<63 - 1)}

	fmt.Printf("Starting load test\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Workers: %d\n", requestWorkerCount)
	fmt.Printf("Duration: %v\n\n", testingDuration)

	workersWg.Add(requestWorkerCount)
	for range requestWorkerCount {
		go func() {
			defer workersWg.Done()
			requestWorker(ctx, url, queue)
		}()
	}

	writerWg.Add(writerWorkerCount)
	for range writerWorkerCount {
		go func() {
			defer writerWg.Done()
			metrics.writeWorker(ctx, queue)
		}()
	}

	time.Sleep(testingDuration)
	cancel()
	workersWg.Wait()
	close(queue)
	writerWg.Wait()

	metrics.Print()

	if err := metrics.WriteErrorsToFile("errors.log"); err != nil {
		fmt.Printf("Failed to write errors: %v\n", err)
	}
}
