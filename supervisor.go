package main

import (
	"fmt"
	"log"
	"net/url"
)

type supervisor struct {
	baseURL *url.URL
	workers chan string
	updates chan update
	queue   []string
	visited map[string]bool
	results []update
	done    chan struct{}
}

func newSupervisor(baseURL string) (*supervisor, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if !base.IsAbs() {
		return nil, fmt.Errorf("not an absolute URL: %s", baseURL)
	}

	s := &supervisor{
		baseURL: base,
		workers: make(chan string),
		updates: make(chan update),
		queue:   []string{baseURL},
		visited: make(map[string]bool),
		results: []update{},
		done:    make(chan struct{}),
	}

	go s.loop()

	return s, nil
}

func (s *supervisor) UpdateChan() chan<- update {
	return s.updates
}

func (s *supervisor) WorkerChan() <-chan string {
	return s.workers
}

func (s *supervisor) Done() <-chan struct{} {
	return s.done
}

func (s *supervisor) Results() []update {
	return s.results
}

func (s *supervisor) loop() {
	for {
		if len(s.queue) == 0 {
			break
		}

		next := s.queue[0]
		s.queue = s.queue[1:]

		if v := s.visited[next]; v {
			continue
		}

		log.Printf("Queue: %d Visited: %d", len(s.queue), len(s.visited))
		result := s.checkAndVisit(next)

		s.visited[result.URL] = true
		s.results = append(s.results, result)

		fmt.Println(result)

		unvisited := s.filterLinks(result.Links)
		s.queue = append(s.queue, unvisited...)
	}
	s.done <- struct{}{}
}

func (s *supervisor) checkAndVisit(rawurl string) update {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return update{
			URL:   rawurl,
			Error: err,
		}
	}

	if parsed.Host != s.baseURL.Host {
		return update{
			URL:     rawurl,
			Skipped: true,
		}
	}

	s.workers <- rawurl
	return <-s.updates
}

func (s *supervisor) filterLinks(links []string) []string {
	unvisited := []string{}

	for _, u := range links {
		canonical, err := canonicalizeURL(s.baseURL, u)
		if err != nil {
			log.Printf("[s] Error parsing link %s: %s", u, err)
		}

		if _, ok := s.visited[canonical]; ok {
			continue
		}

		unvisited = append(unvisited, canonical)
	}
	return unvisited
}
