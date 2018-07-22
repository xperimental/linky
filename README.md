[![Docker Build Status](https://img.shields.io/docker/build/xperimental/linky.svg?style=flat-square)](https://hub.docker.com/r/xperimental/linky/)

# linky

`linky` is a simple link-checker for websites. Given a start URL which returns HTML it will recursively fetch all referenced URLs and check if they return a valid result. It does not leave the domain it was started on so it is suitable for testing a complete website without leaving to external sites.

## Installation

If you have a current (>=1.8) go installation you can simply do:

```
go get github.com/xperimental/linky
```

Binary release can be found on the [Releases page](https://github.com/xperimental/linky/releases).

There's also a Docker image `xperimental/linky`. The `latest` tag points to the master-version and releases are also tagged.

## Usage

```
Usage: linky [options] URL

Options:
  -i, --ignore-referrer   Ignore referrer when checking for duplicate URLs.
  -q, --quiet             Only show errors.
  -v, --verbose           Show all requests including skipped.
```
