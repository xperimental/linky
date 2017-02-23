package html

import (
	"io"

	"golang.org/x/net/html"
)

var attributeMap = map[string]string{
	"a":      "href",
	"link":   "href",
	"img":    "src",
	"script": "src",
}

// ParseLinks parses an HTML file and returns all outgoing links.
func ParseLinks(body io.Reader) []string {
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
