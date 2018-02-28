package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/xperimental/linky/html"
)

type supervisor struct {
	baseURL        *url.URL
	workers        chan location
	updates        chan update
	queue          []location
	visited        map[location]bool
	results        []update
	done           chan struct{}
	showSkipped    bool
	ignoreReferrer bool
	hideOK         bool
}

func newSupervisor(baseURL string, showSkipped bool, ignoreReferrer bool, hideOK bool) (*supervisor, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if !base.IsAbs() {
		return nil, fmt.Errorf("not an absolute URL: %s", baseURL)
	}

	start := location{
		URL:      baseURL,
		Referrer: "",
	}

	s := &supervisor{
		baseURL:        base,
		workers:        make(chan location),
		updates:        make(chan update),
		queue:          []location{start},
		visited:        make(map[location]bool),
		results:        []update{},
		done:           make(chan struct{}),
		showSkipped:    showSkipped,
		ignoreReferrer: ignoreReferrer,
		hideOK:         hideOK,
	}

	go s.loop()

	return s, nil
}

func (s *supervisor) UpdateChan() chan<- update {
	return s.updates
}

func (s *supervisor) WorkerChan() <-chan location {
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

		if s.checkVisited(next) {
			continue
		}

		result := s.checkAndVisit(next)

		s.markVisit(result.Location)
		s.results = append(s.results, result)

		if (!result.Skipped || s.showSkipped) && !(result.IsOK() && s.hideOK) {
			fmt.Println(result)
		}

		unvisited := s.filterLinks(result.Location.URL, result.Links)
		s.queue = append(s.queue, unvisited...)
	}
	s.done <- struct{}{}
}

func (s *supervisor) checkAndVisit(loc location) update {
	parsed, err := url.Parse(loc.URL)
	if err != nil {
		return update{
			Location: loc,
			Error:    err,
		}
	}

	if parsed.Host != s.baseURL.Host {
		return update{
			Location: loc,
			Skipped:  true,
		}
	}

	s.workers <- loc
	return <-s.updates
}

func (s *supervisor) filterLinks(referrer string, links []string) []location {
	refererURL, err := url.Parse(referrer)
	if err != nil {
		refererURL = s.baseURL
	}

	unvisited := []location{}

	for _, u := range links {
		canonical, err := html.CanonicalizeURL(refererURL, u)
		if err != nil {
			log.Printf("[s] Error parsing link %s: %s", u, err)
		}

		canonicalLocation := location{
			URL:      canonical,
			Referrer: referrer,
		}

		if s.checkVisited(canonicalLocation) {
			continue
		}

		unvisited = append(unvisited, canonicalLocation)
	}
	return unvisited
}

func (s *supervisor) markVisit(loc location) {
	if s.ignoreReferrer {
		loc.Referrer = ""
	}

	s.visited[loc] = true
}

func (s *supervisor) checkVisited(loc location) bool {
	if s.ignoreReferrer {
		loc.Referrer = ""
	}

	return s.visited[loc]
}
