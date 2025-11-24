package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
