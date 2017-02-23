package main

import (
	"fmt"
	"time"
)

type update struct {
	URL          string
	Skipped      bool
	Error        error
	Status       int
	ResponseTime time.Duration
	ContentType  string
	Links        []string
}

func (u update) String() string {
	if u.Error != nil {
		return fmt.Sprintf("[ERR] %s (%s; %s)", u.URL, u.ResponseTime, u.Error)
	}

	if u.Skipped {
		return fmt.Sprintf("[SKP] %s", u.URL)
	}

	return fmt.Sprintf("[%03d] %s (%s; %d links)", u.Status, u.URL, u.ResponseTime, len(u.Links))
}

func (u update) IsValid() bool {
	if u.Error != nil {
		return false
	}

	if u.Status < 200 {
		return false
	}

	if u.Status > 299 {
		return false
	}

	return true
}
