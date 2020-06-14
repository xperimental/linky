package main

import "fmt"

var (
	// Version contains the version as set during the build.
	Version = ""

	// GitCommit contains the git commit hash set during the build.
	GitCommit = ""
)

func defaultUserAgent() string {
	if Version == "" {
		return "linky/unknown"
	}

	return fmt.Sprintf("linky/%s", Version)
}
