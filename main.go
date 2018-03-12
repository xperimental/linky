package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	var verbose bool
	var ignoreReferrer bool
	var quiet bool
	pflag.BoolVarP(&verbose, "verbose", "v", false, "Show all requests including skipped.")
	pflag.BoolVarP(&quiet, "quiet", "q", false, "Only show errors.")
	pflag.BoolVarP(&ignoreReferrer, "ignore-referrer", "i", false, "Ignore referrer when checking for duplicate URLs.")
	pflag.Parse()

	startURL := pflag.Arg(0)

	if len(startURL) == 0 {
		fmt.Println("Usage: linky [options] URL\n\nOptions:")
		pflag.PrintDefaults()
		return
	}

	fmt.Printf("URL: %s\n", startURL)

	s, err := newSupervisor(startURL, verbose, quiet, ignoreReferrer)
	if err != nil {
		log.Fatalf("Error creating supervisor: %s", err)
	}

	newWorker(s.WorkerChan(), s.UpdateChan())

	<-s.Done()

	if showResults(s.Results()) > 0 {
		os.Exit(1)
	}
}
