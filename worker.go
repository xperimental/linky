package main

import (
	"io"
	"net/http"
	"time"

	"golang.org/x/net/html"
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
	if err != nil {
		result.Error = err
		return
	}
	defer res.Body.Close()

	result.ResponseTime = time.Since(start)
	result.Status = res.StatusCode

	if res.StatusCode >= http.StatusOK && res.StatusCode < 300 {
		result.Links = w.parseLinks(res.Body)
	}

	return result
}

func (w *worker) parseLinks(body io.Reader) []string {
	links := []string{}
	page := html.NewTokenizer(body)
	for {
		tt := page.Next()
		if tt == html.ErrorToken {
			return links
		}

		if tt == html.StartTagToken {
			token := page.Token()
			switch token.DataAtom.String() {
			case "a":
				link := extractAttribute(token, "href")
				if len(link) > 0 {
					links = append(links, link)
				}
			}
		}
	}
}

func extractAttribute(token html.Token, attrName string) string {
	for _, attr := range token.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}

	return ""
}
