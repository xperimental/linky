package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type worker struct {
	client   *http.Client
	linkChan chan string
	stopChan chan struct{}
}

func newWorker(linkChan chan string) *worker {
	stop := make(chan struct{})
	w := &worker{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		linkChan: linkChan,
		stopChan: stop,
	}
	go w.loop()

	return w
}

func (w *worker) Stop() {
	w.stopChan <- struct{}{}
}

func (w *worker) loop() {
	for {
		select {
		case url := <-w.linkChan:
			if err := w.checkURL(url); err != nil {
				log.Printf("Error checking %s: %s", url, err)
			}
		case <-w.stopChan:
			return
		}
	}
}

func (w *worker) checkURL(url string) error {
	res, err := w.client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("invalid status code %d: %s", res.StatusCode, res.Status)
	}

	contentType := res.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		links, err := w.parseLinks(res.Body)
		if err != nil {
			return err
		}

		log.Println("Links:")
		for _, l := range links {
			log.Println(" * " + l)
		}
	}

	return nil
}

func (w *worker) parseLinks(body io.Reader) ([]string, error) {
	links := []string{}
	page := html.NewTokenizer(body)
	for {
		tt := page.Next()
		if tt == html.ErrorToken {
			return links, nil
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
			return trimAnchor(attr.Val)
		}
	}

	return ""
}

func trimAnchor(input string) string {
	tokens := strings.SplitN(input, "#", 2)
	return tokens[0]
}
