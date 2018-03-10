package main

import (
	_ "bufio"
	_ "encoding/json"
	"fmt"
	_ "io"
	"io/ioutil"
	"log"
	_ "net"
	"net/http"
	"os"
	_ "runtime"
	"sync"
)

const (
	boardURLTemplate  = "https://2ch.hk/%s/threads.json"
	threadURLTemplate = "https://2ch.hk/%s/res/%d.json"
)

func init() {
	log.SetOutput(os.Stdout)
}

// func fetchJSON(url) {

// }

// func fetchBoard(boardName string) {
//     return fetchJSON(fmt.Sprintf(boardURLTemplate, boardName))
// }

func fetch(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	read, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	return read
}

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

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

	wg.Wait()

	if closeResults {
        log.Println("Closing errors & results")
        
        close(worker.Errors)
		close(worker.Results)
	}
}

func (worker *FetchingWorker) Close() {
	close(worker.Results)
	close(worker.Errors)
    
    if _, ok := <-worker.Urls; ok {
        close(worker.Urls)
    }
}

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	read, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return read, nil
}

func dump(name, contents string) {
	f, err := os.Create(name)
	if err != nil {
		log.Printf("Failed to create file %s\n", name)

		return
	}

	defer f.Close()

	if written, err := f.WriteString(contents); err != nil {
		log.Printf("Error writing to file %v %s\n", written, name)

		return
	}

}

func main() {
	var wg sync.WaitGroup

	worker := FetchingWorker{
		Urls:    make(chan string, 100),
		Results: make(chan UrlResultPair, 100),
		Errors:  make(chan UrlErrorPair, 100),
	}

	worker.Urls <- "https://2ch.hk/b/res/172168647.json"
	worker.Urls <- "https://2ch.hk/b/threads.json"

	wg.Add(2)

	go worker.Work(true)

	go func(w *FetchingWorker) {
		defer wg.Done()

		i := 0
		for result := range w.Results {
			log.Printf("Got result from %s\n", result.URL)

			i++
			dump(fmt.Sprintf("dump%v.json", i), string(result.Result))
		}
		log.Println("Processed all results")
	}(&worker)

	go func(w *FetchingWorker) {
		defer wg.Done()

		for err := range w.Errors {
			log.Printf("URL %s failed. %v", err.URL, err.Error)
		}
		log.Println("Processed all errors")
	}(&worker)

	close(worker.Urls)

	wg.Wait()

	log.Println("Done")
}
