package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func printStartup(config Config) {
	fmt.Printf("Starting load test\n")
	fmt.Printf("URL: %s\n", config.url)
	fmt.Printf("Workers: %d\n", config.requestWorkersCount)
	fmt.Printf("Requests per worker: %d\n", config.requestsPerWorker)

}

func setupShutdownWorker(cancel context.CancelFunc) {
	// setup channel for sigint
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT)

	go func() {
		<-sigCh
		fmt.Println("\nReceived interrupt signal, shutting down")
		cancel()
	}()
}
