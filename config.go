package main

type Config struct {
	requestWorkersCount int
	requestsPerWorker   int

	writerWorkersCount int
	queueChannelSize   int
	url                string
}

func getDefaultConfig(url string) (c Config) {
	c = Config{
		requestWorkersCount: 500,
		requestsPerWorker:   500,
		writerWorkersCount:  50,
		queueChannelSize:    10_000,
		url:                 url,
	}
	return
}
