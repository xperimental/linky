package main

import (
	"fmt"
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
	var concurrency int
	pflag.IntVarP(&concurrency, "concurrency", "c", 1, "Number of workers to use concurrently.")
	pflag.Parse()

	startURL := pflag.Arg(0)

	if len(startURL) == 0 {
		fmt.Println("Usage: linky [options] URL\n\nOptions:")
		pflag.PrintDefaults()
		return
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
