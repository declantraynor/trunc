# Trunc

A simple URL shortening service written in [Go](https://golang.org).

## API

The service provides a HTTP API which supports a couple of basic operations.

_Note: these examples assume the service is accessible at http://localhost:5000_

### Shorten a URL

```
POST /shorten
Host: http://localhost:5000
Content-Type: application/json
{
    "url": "http://example.com/really/long/url?foo=bar123456789"
}

---

HTTP/1.1 200 OK
Content-Type: application/json
{
    "long_url": "http://example.com/really/long/url?foo=bar123456789",
    "short_url": "http://localhost:5000/cTRD2S4d"
}
```

### Redirect from a short URL

```
GET /cTRD2S4d
Host: http://localhost:5000

---

HTTP/1.1 301 Moved Permanently
Location: http://example.com/really/long/url?foo=bar123456789
```

## Running locally

Running the service locally requires [Docker](https://docs.docker.com/engine/) and [Docker Compose](https://docs.docker.com/compose/).

A Makefile is provided. This collects commonly-used development commands into easy-to-remember aliases.

To run the service:

```sh
$ make run
```

This will start the service running in the background. By default, it is is configured to listen for new connections at http://localhost:5000.

To stop the service:

```sh
$ make stop
```

## Running the tests

The service has a suite of unit tests, which can be run as follows:

```sh
$ make test
```

## Using the CLI

Trunc comes with a command line tool that you can use to quickly generate short URLs:

```sh
$ go install ./cmd/trunc
$ trunc http://example.com/really/long/url?foo=bar123456789 

------------------------------
http://localhost:5000/ko9RVl4o
------------------------------
```

For configuration info, see

```sh
$ trunc help
```
