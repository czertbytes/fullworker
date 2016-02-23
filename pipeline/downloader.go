package pipeline

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type downloader struct {
	process
	Out chan Payload

	urls []string
}

func NewDownloader(urls []string) *downloader {
	return &downloader{
		urls: urls,
		Out:  make(chan Payload, jobBuffer),
	}
}

func (p *downloader) Run() {
	defer close(p.Out)

	var wg sync.WaitGroup

	wg.Add(len(p.urls))

	for _, s := range p.urls {
		go func(s string) {
			defer wg.Done()
			start := time.Now()

			fmt.Printf("started downloader %s\n", s)
			resp, err := http.Get(s)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			p.Out <- Payload{
				URL:     s,
				Content: string(body),
				Time:    time.Now().Sub(start),
			}
			fmt.Printf("completed downloader %s\n", s)
		}(s)
	}

	wg.Wait()
}
