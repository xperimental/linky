package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

func showResults(results []update) int {
	successful := 0
	nonOk := 0
	skipped := 0
	errors := 0
	var totalTime time.Duration
	codes := make(map[int]int)
	codes[200] = 0
	types := make(map[string]int)
	for _, v := range results {
		if v.Skipped {
			skipped++
			continue
		}

		totalTime += v.ResponseTime

		if v.Error != nil {
			errors++
			continue
		}

		successful++

		codes[v.Status]++
		if !v.IsOK() {
			nonOk++
		}

		if len(v.ContentType) > 0 {
			types[v.ContentType]++
		} else {
			types["Unknown"]++
		}
	}

	fmt.Println("\nResults:")
	fmt.Printf(" %5d total\n", len(results))
	fmt.Printf(" %5d successful\n", successful)
	fmt.Printf(" %5d errors\n", errors)
	fmt.Printf(" %5d skipped\n", skipped)
	fmt.Printf("Total time: %s\n", totalTime)

	if len(codes) > 1 {
		showStatusCodes(codes)
	}

	if len(types) > 0 {
		showContentTypes(types)
	}

	return nonOk + errors
}

func showStatusCodes(codes map[int]int) {
	fmt.Println("\nStatus Codes:")
	sortCodes := []struct {
		status int
		count  int
	}{}
	for s, c := range codes {
		sortCodes = append(sortCodes, struct {
			status int
			count  int
		}{
			status: s,
			count:  c,
		})
	}
	sort.Slice(sortCodes, func(i int, j int) bool {
		a := sortCodes[i].count
		b := sortCodes[j].count
		if a == b {
			return sortCodes[i].status < sortCodes[j].status
		}

		return a > b
	})
	for _, v := range sortCodes {
		if v.count > 0 {
			fmt.Printf(" %5d %d\n", v.count, v.status)
		}
	}
}

func showContentTypes(types map[string]int) {
	fmt.Println("\nContent Types:")
	sortTypes := []struct {
		contentType string
		count       int
	}{}
	for t, c := range types {
		sortTypes = append(sortTypes, struct {
			contentType string
			count       int
		}{
			contentType: t,
			count:       c,
		})
	}
	sort.Slice(sortTypes, func(i int, j int) bool {
		a := sortTypes[i].count
		b := sortTypes[j].count
		if a == b {
			return strings.Compare(sortTypes[i].contentType, sortTypes[j].contentType) < 0
		}

		return a > b
	})
	for _, v := range sortTypes {
		fmt.Printf(" %5d %s\n", v.count, v.contentType)
	}
}

func main() {
	var showSkipped bool
	pflag.BoolVar(&showSkipped, "show-skipped", false, "Show skipped URLs.")
	pflag.Parse()

	startURL := pflag.Arg(0)

	if len(startURL) == 0 {
		fmt.Println("Usage: linky [options] URL\n\nOptions:")
		pflag.PrintDefaults()
		return
	}

	fmt.Printf("URL: %s\n", startURL)

	s, err := newSupervisor(startURL, showSkipped)
	if err != nil {
		log.Fatalf("Error creating supervisor: %s", err)
	}

	newWorker(s.WorkerChan(), s.UpdateChan())

	<-s.Done()

	if showResults(s.Results()) > 0 {
		os.Exit(1)
	}
}
