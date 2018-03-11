package fetch

import (
	"log"
	"sync"
)

type UrlResultPair struct {
	URL    string
	Result []byte
}

type UrlErrorPair struct {
	URL   string
	Error error
}

type FetchingWorker struct {
	Urls    chan string
	Results chan UrlResultPair
	Errors  chan UrlErrorPair
}

func NewFetchingWorker(buffSize uint32) *FetchingWorker {
	return &FetchingWorker{
		Urls:    make(chan string, buffSize),
		Results: make(chan UrlResultPair, buffSize),
		Errors:  make(chan UrlErrorPair, buffSize),
	}
}

func (worker *FetchingWorker) Work(closeResults bool) {
	var wg sync.WaitGroup

	for url := range worker.Urls {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()

			log.Printf("Fetching url %s\n", u)

			bytes, err := Fetch(u)
			if err != nil {
				worker.Errors <- UrlErrorPair{u, err}
			} else {
				worker.Results <- UrlResultPair{u, bytes}
			}
		}(url)
	}

	go func(wrk *FetchingWorker, w *sync.WaitGroup) {
		w.Wait()

		if closeResults {
			log.Println("Closing errors & results")

			close(wrk.Errors)
			close(wrk.Results)
		}
	}(worker, &wg)
}

func (worker *FetchingWorker) Close() {
	close(worker.Results)
	close(worker.Errors)

	if _, ok := <-worker.Urls; ok {
		close(worker.Urls)
	}
}
