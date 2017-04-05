package main

import (
	"fmt"
	"time"
)

type location struct {
	URL      string
	Referrer string
}

func (l location) String() string {
	if len(l.Referrer) == 0 {
		return l.URL
	}

	return fmt.Sprintf("%s (ref %s)", l.URL, l.Referrer)
}

type update struct {
	Location     location
	Skipped      bool
	Error        error
	Status       int
	ResponseTime time.Duration
	ContentType  string
	Links        []string
}

func (u update) String() string {
	if u.Error != nil {
		return fmt.Sprintf("[ERR] %s (%s; %s)", u.Location, u.ResponseTime, u.Error)
	}

	if u.Skipped {
		return fmt.Sprintf("[SKP] %s", u.Location)
	}

	return fmt.Sprintf("[%03d] %s (%s; %d links)", u.Status, u.Location, u.ResponseTime, len(u.Links))
}

func (u update) IsOK() bool {
	return u.Status >= 200 && u.Status <= 299
}
