package main

import (
	_ "bufio"
	_ "encoding/json"
	_ "io"
	"log"
	_ "net"
	"net/url"
	"os"
	"path/filepath"
	_ "runtime"
	"strings"
	"sync"

	"github.com/godstopme/myonmyon/dvach"
	"github.com/godstopme/myonmyon/fetch"
	"github.com/godstopme/myonmyon/utils"
)

const (
	baseURLString     = "https://2ch.hk"
	boardURLTemplate  = "https://2ch.hk/%s/threads.json"
	threadURLTemplate = "https://2ch.hk/%s/res/%d.json"
)

var BaseURL *url.URL

func init() {
	BaseURL, _ = url.Parse("https://2ch.hk")

	log.SetOutput(os.Stdout)
}

func fetchThreadFiles(posts []dvach.Post) {
	var wg sync.WaitGroup

	// postsCount := len(posts)
	worker := fetch.NewFetchingWorker(4)

	wg.Add(2)
	go worker.Work(true)

	go func() {
		defer close(worker.Urls)

		for _, post := range posts {
			for _, file := range post.Files {
				relFileURL, _ := url.Parse(file.Path)
				fileURL := BaseURL.ResolveReference(relFileURL)

				if file.Type == 6 {
					relThumbnailURL, _ := url.Parse(file.ThumbnailPath)
					thumbnailURL := BaseURL.ResolveReference(relThumbnailURL)

					worker.Urls <- thumbnailURL.String()
				}

				worker.Urls <- fileURL.String()
			}
		}
	}()

	go func(w *fetch.FetchingWorker) {
		defer wg.Done()

		for result := range w.Results {
			log.Printf("Got result from %s\n", result.URL)
			wg.Add(1)

			go func() {
				defer wg.Done()

				fname := strings.Replace(strings.Replace(result.URL, "/", "_", -1), ":", "", -1)
				fileName := filepath.Join("data/", fname)

				utils.SaveFile(fileName, result.Result)
			}()
		}
		log.Println("Processed all results")
	}(worker)

	go func(w *fetch.FetchingWorker) {
		defer wg.Done()

		for err := range w.Errors {
            log.Printf("URL %s failed. %v", err.URL, err.Error)
            // w.Urls <- err.URL
		}
		log.Println("Processed all errors")
	}(worker)

	wg.Wait()
	log.Println("Done fetching")
}

func main() {
	var wg sync.WaitGroup

	urls := []string{
		// "https://2ch.hk/wp/res/60491.html",
		"https://2ch.hk/wp/res/54986.html",
		// "https://2ch.hk/aa/res/85799.html",
		// "https://2ch.hk/aa/res/67884.html",
	}
	wg.Add(2)

	worker := fetch.NewFetchingWorker(uint32(len(urls)))
	go worker.Work(true)

	go func() {
		defer wg.Done()

		for result := range worker.Results {
			wg.Add(1)
			log.Println("Fetched thread", result.URL)
			posts, err := dvach.UnmarshalPosts(result.Result)
			if err != nil {
				log.Println("Failed to unmarshal thread json", result.URL)
			}

			go func() {
				defer wg.Done()

				fetchThreadFiles(posts)
			}()
		}
	}()
	go func() {
		defer wg.Done()

		for result := range worker.Errors {
			log.Println("Got error from ", result.URL, result.Error)
		}
	}()

	for _, url := range urls {
		log.Println("Sending thread url", url)
		worker.Urls <- strings.Replace(url, ".html", ".json", -1)
	}
	close(worker.Urls)
	wg.Wait()

	log.Println("Done")
}
