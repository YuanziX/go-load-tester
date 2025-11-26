package main

type Config struct {
	requestWorkersCount int
	requestsPerWorker   int

	writerWorkersCount int
	queueChannelSize   int
	url                string
}

func getConfigWithParams(url string, rwc, rpw int) Config {
	return Config{
		requestWorkersCount: rwc,
		requestsPerWorker:   rpw,
		writerWorkersCount:  max(5, min(rwc/100, 50)),
		queueChannelSize:    DefaultQueueSize,
		url:                 url,
	}
}
