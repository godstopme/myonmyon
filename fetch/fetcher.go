package fetch

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

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

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("Not found %s", url)
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Non 200 http code")
	}
	read, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return read, nil
}
