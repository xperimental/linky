package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	var (
		verbose        bool
		ignoreReferrer bool
		quiet          bool
		userAgent      = defaultUserAgent()
	)
	pflag.BoolVarP(&verbose, "verbose", "v", false, "Show all requests including skipped.")
	pflag.BoolVarP(&quiet, "quiet", "q", false, "Only show errors.")
	pflag.BoolVarP(&ignoreReferrer, "ignore-referrer", "i", false, "Ignore referrer when checking for duplicate URLs.")
	pflag.StringVar(&userAgent, "user-agent", userAgent, "HTTP User-Agent header to send. If empty, the default Go User-Agent will be used.")
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

	newWorker(s.WorkerChan(), s.UpdateChan(), userAgent)

	<-s.Done()

	if showResults(s.Results()) > 0 {
		os.Exit(1)
	}
}
