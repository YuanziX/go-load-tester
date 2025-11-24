package main

import "time"

const (
	DefaultQueueSize = 10_000
	DefaultTimeout   = 5 * time.Second
	MaxDuration      = time.Duration(1<<63 - 1)
)
