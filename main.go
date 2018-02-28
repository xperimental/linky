package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	var showSkipped bool
	var ignoreReferrer bool
	var hideOK bool
	pflag.BoolVarP(&showSkipped, "show-skipped", "s", false, "Show skipped URLs.")
	pflag.BoolVarP(&ignoreReferrer, "ignore-referrer", "i", false, "Ignore referrer when checking for duplicate URLs.")
	pflag.BoolVarP(&hideOK, "hide-ok", "o", false, "Do not show URLs with good response codes.")
	pflag.Parse()

	startURL := pflag.Arg(0)

	if len(startURL) == 0 {
		fmt.Println("Usage: linky [options] URL\n\nOptions:")
		pflag.PrintDefaults()
		return
	}

	fmt.Printf("URL: %s\n", startURL)

	s, err := newSupervisor(startURL, showSkipped, ignoreReferrer, hideOK)
	if err != nil {
		log.Fatalf("Error creating supervisor: %s", err)
	}

	newWorker(s.WorkerChan(), s.UpdateChan())

	<-s.Done()

	if showResults(s.Results()) > 0 {
		os.Exit(1)
	}
}
