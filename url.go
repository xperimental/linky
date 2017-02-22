package main

import "net/url"

func canonicalizeURL(base *url.URL, rawurl string) (string, error) {
	parsed, err := base.Parse(rawurl)
	if err != nil {
		return "", nil
	}

	parsed.Fragment = ""
	return parsed.String(), nil
}
