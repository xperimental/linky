package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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

	if res.StatusCode < http.StatusOK || res.StatusCode >= 300 {
		result.Error = fmt.Errorf("non-ok status code: %d", res.StatusCode)
		return
	}

	result.ContentType = res.Header.Get("Content-Type")

	if strings.HasPrefix(result.ContentType, "text/html") {
		result.Links = w.parseLinks(res.Body)
	}

	return result
}

var attributeMap = map[string]string{
	"a":      "href",
	"link":   "href",
	"img":    "src",
	"script": "src",
}

func (w *worker) parseLinks(body io.Reader) []string {
	links := []string{}
	page := html.NewTokenizer(body)
	for {
		tt := page.Next()
		if tt == html.ErrorToken {
			return links
		}

		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			token := page.Token()
			attribute, ok := attributeMap[token.DataAtom.String()]
			if !ok {
				continue
			}

			link := extractAttribute(token, attribute)
			if len(link) > 0 {
				links = append(links, link)
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
