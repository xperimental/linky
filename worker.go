package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/xperimental/linky/html"
)

var clientTimeout = 60 * time.Second

type worker struct {
	client    *http.Client
	locations <-chan location
	updates   chan<- update
	userAgent string
}

func newWorker(locations <-chan location, updates chan<- update, userAgent string) *worker {
	w := &worker{
		client: &http.Client{
			Timeout: clientTimeout,
		},
		locations: locations,
		updates:   updates,
		userAgent: userAgent,
	}

	go w.loop()

	return w
}

func (w *worker) loop() {
	for l := range w.locations {
		result := w.fetchURL(l)

		go func() {
			w.updates <- result
		}()
	}
}

func (w *worker) fetchURL(location location) (result update) {
	result.Location = location
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, location.URL, nil)
	if err != nil {
		result.Error = fmt.Errorf("can not create request: %s", err)
		return
	}

	if w.userAgent != "" {
		req.Header.Set("User-Agent", w.userAgent)
	}

	res, err := w.client.Do(req)
	result.ResponseTime = time.Since(start)

	if err != nil {
		result.Error = err
		return
	}
	defer res.Body.Close()

	result.Status = res.StatusCode

	if res.StatusCode < http.StatusOK || res.StatusCode >= 300 {
		return
	}

	result.ContentType = res.Header.Get("Content-Type")

	if strings.HasPrefix(result.ContentType, "text/html") {
		result.Links = html.ParseLinks(res.Body)
	}

	return result
}
