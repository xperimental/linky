package main

import (
	"log"

	"github.com/spf13/pflag"
)

func makeWorkers(number int, linkChan chan string) []*worker {
	workers := make([]*worker, number)
	for i := range workers {
		workers[i] = newWorker(linkChan)
	}
	return workers
}

func main() {
	var startURL string
	var concurrency int
	pflag.StringVarP(&startURL, "url", "u", "", "URL to use as base for linkchecker.")
	pflag.IntVarP(&concurrency, "concurrency", "c", 1, "Number of workers to use concurrently.")
	pflag.Parse()

	if len(startURL) == 0 {
		log.Fatalln("Need to provide a start URL.")
	}

	if concurrency < 1 {
		log.Fatalln("Need at least one worker.")
	}

	log.Printf("URL: %s", startURL)
	log.Printf("Workers: %d", concurrency)

	linkChan := make(chan string)
	workers := makeWorkers(concurrency, linkChan)

	log.Println("Checking...")
	linkChan <- startURL

	log.Println("Stop workers...")
	for _, w := range workers {
		w.Stop()
	}
}
