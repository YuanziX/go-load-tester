package main

import "fmt"

func main() {
	// setup globally required vars
	metrics := getMetricsObject()

	serveUI(&metrics)

	if err := metrics.WriteErrorsToFile("errors.log"); err != nil {
		fmt.Printf("Failed to write errors: %v\n", err)
	}
}
