package main

import (
	"context"
	"fmt"
)

var url string = "http://localhost:3000"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// setup globally required vars
	config := getDefaultConfig(url)
	metrics := getMetricsObject()

	setupShutdownWorker(cancel)
	setupLoadTesterWorkers(&config, &metrics, ctx)

	printStartup(config)
	metrics.Print()

	if err := metrics.WriteErrorsToFile("errors.log"); err != nil {
		fmt.Printf("Failed to write errors: %v\n", err)
	}
}
