package main

import (
	"sync"

	"github.com/czertbytes/fullworker/pipeline"
	"github.com/czertbytes/fullworker/worker"
)

var (
	urls = []string{
		"https://www.google.com",
		"https://www.apple.com",
		"https://www.golang.org",
		"https://www.microsoft.com",
		"https://www.yahoo.com",
		"https://www.amazon.com",
		"https://www.netflix.com",
		"https://www.spotify.com",
		"https://www.dropbox.com",
		"https://www.teslamotors.com",
	}
)

func main() {
	//jobPipeline()
	workerQueue()
}

func jobPipeline() {
	jobPipeline := pipeline.NewPipeline()

	downloader := pipeline.NewDownloader(urls)
	sizer := pipeline.NewSizer()
	printer := pipeline.NewPrinter()

	// connect pipelines
	sizer.In = downloader.Out
	printer.In = sizer.Out

	jobPipeline.AddProcesses(downloader, sizer, printer)
	jobPipeline.Run()
}

func workerQueue() {
	dispatcher := worker.NewDispatcher(5)
	dispatcher.Run()

	var wg sync.WaitGroup
	wg.Add(len(urls))

	// add job
	for _, url := range urls {
		dispatcher.StartJob(&wg, url)
	}

	wg.Wait()
}
