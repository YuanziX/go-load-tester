# Go-Load-Tester

A very simple load tester written in golang to learn the concurrency model.

Uses workers to make requests on some url, writes responses into a channel, other workers process these responses into a metrics struct using a mutex lock.
