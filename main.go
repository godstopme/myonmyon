package main

import (
	_ "bufio"
	_ "encoding/json"
	"fmt"
	_ "io"
	"log"
	_ "net"
	"os"
	_ "runtime"
	"sync"

	"github.com/godstopme/myonmyon/fetch"
    "github.com/godstopme/myonmyon/utils"
    "github.com/godstopme/myonmyon/dvach"
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

func main() {
	var wg sync.WaitGroup

	worker := fetch.NewFetchingWorker(1)

	wg.Add(2)

	go worker.Work(true)
    worker.Urls <- "https://2ch.hk/b/res/172253674.json"
	go func(w *fetch.FetchingWorker) {
		defer wg.Done()

		i := 0
		for result := range w.Results {
			log.Printf("Got result from %s\n", result.URL)

			i++
            utils.Dump(fmt.Sprintf("dump%v.json", i), string(result.Result))
            posts, err := dvach.UnmarshalPosts(result.Result)
            if err != nil {
                log.Println("Error unmarshaling json!", err)
            } else {
                log.Println(posts)
            }
		}
		log.Println("Processed all results")
	}(worker)

	go func(w *fetch.FetchingWorker) {
		defer wg.Done()

		for err := range w.Errors {
			log.Printf("URL %s failed. %v", err.URL, err.Error)
		}
		log.Println("Processed all errors")
	}(worker)

	close(worker.Urls)

	wg.Wait()

	log.Println("Done")
}
