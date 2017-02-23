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
		return fmt.Sprintf("[ERR] %s (%s)", u.URL, u.Error)
	}

	if u.Skipped {
		return fmt.Sprintf("[SKP] %s", u.URL)
	}

	return fmt.Sprintf("[%03d] %s (%s; %d links)", u.Status, u.URL, u.ResponseTime, len(u.Links))
}
