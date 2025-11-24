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
		writerWorkersCount:  max(1000, rwc/10),
		queueChannelSize:    DefaultQueueSize,
		url:                 url,
	}
}

func getDefaultConfig(url string) Config {
	return Config{
		requestWorkersCount: 500,
		requestsPerWorker:   500,
		writerWorkersCount:  50,
		queueChannelSize:    10_000,
		url:                 url,
	}
}
