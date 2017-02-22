package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/pflag"
)

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

	results := s.Results()
	successful := 0
	skipped := 0
	errors := 0
	var totalTime time.Duration
	for _, v := range results {
		if v.Error != nil {
			errors++
			continue
		}

		if v.Skipped {
			skipped++
			continue
		}

		successful++
		totalTime += v.ResponseTime
	}

	fmt.Println("Results:")
	fmt.Printf(" %5d total\n", len(results))
	fmt.Printf(" %5d successful\n", successful)
	fmt.Printf(" %5d skipped\n", skipped)
	fmt.Printf(" %5d errors\n", errors)
	fmt.Printf("Total time: %s\n", totalTime)
}
