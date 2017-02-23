package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/xperimental/linky/html"
)

var clientTimeout = 60 * time.Second

type worker struct {
	client  *http.Client
	urls    <-chan string
	updates chan<- update
}

func newWorker(urls <-chan string, updates chan<- update) *worker {
	w := &worker{
		client: &http.Client{
			Timeout: clientTimeout,
		},
		urls:    urls,
		updates: updates,
	}

	go w.loop()

	return w
}

func (w *worker) loop() {
	for u := range w.urls {
		result := w.fetchURL(u)

		go func() {
			w.updates <- result
		}()
	}
}

func (w *worker) fetchURL(url string) (result update) {
	result.URL = url
	start := time.Now()

	res, err := w.client.Get(url)
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
