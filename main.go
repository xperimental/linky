package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/pflag"
)

func showResults(results []update) {
	successful := 0
	skipped := 0
	errors := 0
	var totalTime time.Duration
	types := make(map[string]int)
	for _, v := range results {
		if v.Error != nil {
			errors++
			continue
		}

		if v.Skipped {
			skipped++
			continue
		}

		if len(v.ContentType) > 0 {
			types[v.ContentType]++
		} else {
			types["Unknown"]++
		}

		successful++
		totalTime += v.ResponseTime
	}

	fmt.Println("\nResults:")
	fmt.Printf(" %5d total\n", len(results))
	fmt.Printf(" %5d successful\n", successful)
	fmt.Printf(" %5d skipped\n", skipped)
	fmt.Printf(" %5d errors\n", errors)
	fmt.Printf("Total time: %s\n", totalTime)

	fmt.Println("\nContent Types:")
	for t, c := range types {
		fmt.Printf(" %5d %s\n", c, t)
	}
}

func main() {
	var concurrency int
	var showSkipped bool
	pflag.IntVarP(&concurrency, "concurrency", "c", 1, "Number of workers to use concurrently.")
	pflag.BoolVar(&showSkipped, "show-skipped", false, "Show skipped URLs.")
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

	fmt.Printf("URL: %s\n", startURL)

	s, err := newSupervisor(startURL, showSkipped)
	if err != nil {
		log.Fatalf("Error creating supervisor: %s", err)
	}

	for i := 0; i < concurrency; i++ {
		newWorker(s.WorkerChan(), s.UpdateChan())
	}

	<-s.Done()

	showResults(s.Results())
}
