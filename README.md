# linky

`linky` is a simple link-checker for websites. Given a start URL which returns HTML it will recursively fetch all referenced URLs and check if they return a valid result. It does not leave the domain it was started on so it is suitable for testing a complete website without leaving to external sites.

## Installation

If you have a current (>=1.8) go installation you can simply do:

```
go get github.com/xperimental/linky
```

Binary releases will be coming soon.

## Usage

```
Usage: linky [options] URL

Options:
      --show-skipped   Show skipped URLs.
```
