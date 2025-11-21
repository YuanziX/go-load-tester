package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var url string = "http://localhost:3000"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// setup config
	config := getDefaultConfig(url)

	queue := make(chan RequestResult, 10_000)

	// setup channel for sigint
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT)

	go func() {
		<-sigCh
		fmt.Println("\nReceived interrupt signal, shutting down")
		cancel()
	}()

	var workersWg sync.WaitGroup
	var writerWg sync.WaitGroup

	metrics := getMetricsObject()

	fmt.Printf("Starting load test\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Workers: %d\n", config.requestWorkersCount)
	fmt.Printf("Requests per worker: %d\n", config.requestsPerWorker)

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

	metrics.Print()

	if err := metrics.WriteErrorsToFile("errors.log"); err != nil {
		fmt.Printf("Failed to write errors: %v\n", err)
	}
}
