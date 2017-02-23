package html

import "net/url"

// CanonicalizeURL parses an URL relative to a base and removes the fragment.
func CanonicalizeURL(base *url.URL, rawurl string) (string, error) {
	parsed, err := base.Parse(rawurl)
	if err != nil {
		return "", err
	}

	parsed.Fragment = ""
	return parsed.String(), nil
}
