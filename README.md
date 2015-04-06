# clidemo
[![License MIT](https://img.shields.io/npm/l/express.svg)](http://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/composer22/clidemo.svg?branch=master)](http://travis-ci.org/composer22/clidemo)
[![Current Release](https://img.shields.io/badge/release-v0.1.0--alpha-brightgreen.svg)](https://github.com/composer22/clidemo/releases/tag/v0.1.0-alpha)
[![Coverage Status](https://coveralls.io/repos/composer22/clidemo/badge.svg?branch=master)](https://coveralls.io/r/composer22/clidemo?branch=master)

A text parser for counting words and sentence locations written in [Go.](http://golang.org)

## About

This is a small demonstration app of some Golang functions and features. Initially intended
to act as a bed for testing Golang development, it could form the framework for a larger application.

Some demonstration objectives:

* CLI flag parameter passing.
* CLI submission of data by file or by pipe.
* Concurrent worker jobqueue (here with a simple text parser/scanner).
* Webserver API:
    + Custom connection listener for optional throttling of connections.
    + Middleware for globally validating incoming requests.
    + Status information from the system.
    + Integration of pprof for performance profiling.
    + JSON encoding of the response body.
    + Standard RESTful request and response header usage.

For TODOs, please see TODO.md

## Usage

```
Description: Parse text counting words and sentence locations, this command can be
evoked as either a command line utility or as a stand alone server process.

Usage: clidemo [options...] [input_filename]

Server options:
    -N, --name NAME                  NAME of the server (default: empty field).
    -h, --hostname HOSTNAME          HOSTNAME of the server (default: localhost).
    -p, --port PORT                  PORT to listen on (default: 49152).
	-L, --profiler_port PORT         *PORT the profiler is listening on (default: off).
    -n, --connections MAX            *MAX server connections allowed (default: unlimited).
    -W, --workers MAX                MAX running workers allowed (default: 1000).
    -X, --procs MAX                  *MAX processor cores to use from the machine.

    -d, --debug                      Enable debugging output (default: false)

     *  Anything <= 0 is no change to the environment (default: 0).

File input options:
    -f, --file FILE                  Process input FILE

Common options:
    -h, --help                       Show this message
    -V, --version                    Show version

Examples:

    # Server mode activated as "San Francisco" on localhost port 8080;
	# 10 conns; 30 workers; 2 processors
    clidemo -N "San Francisco" -p 8080 -n 10 -W 30 -X 2

	# File input using -f flag with debug option
	clidemo -f /tmp/inputfiles/foo/bar.txt -d > out.txt

	# Implicit file input (no -f flag needed)
	clidemo /tmp/inputfiles/foo/bar.txt > out.txt

	# Piping input
	cat /tmp/inputfiles/foo/bar.txt | clidemo > out.txt

```

## Configuration

```
command line flags only

```

## Building

This code currently requires version 1.42 or higher of Go, but we encourage the use of the latest stable release.

Information on Golang installation, including pre-built binaries, is available at
<http://golang.org/doc/install>.  Stable branches of operating system packagers provided by
your OS vendor may not be sufficient.

Run `go version` to see the version of Go which you have installed.

Run `go build` inside the directory to build.

Run `go test ./...` to run the unit regression tests.

A successful build run produces no messages and creates an executable called `clidemo` in this
directory.  You can invoke that binary, with no options to start a server with acceptable standalone defaults.

Run `go help` for more guidance, and visit <http://golang.org/> for tutorials, presentations, references and more.

## API calls

Header should contain:

* Accept: application/json
* Authorization: Bearer with token
* Content-Type: application/json

Example cURL:

```

$ curl -i -H "Accept: application/json" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer 3A3E6C4C51F12DF2415682CCF9D18" \
-X GET "http://0.0.0.0:8080/v1.0/alive"

HTTP/1.1 200 OK
Content-Type: application/json;charset=utf-8
Date: Fri, 03 Apr 2015 17:29:17 +0000
Server: San Francisco
X-Request-Id: DC8D9C2E-8161-4FC0-937F-4CA7037970D5
Content-Length: 0


```

URL Endpoints:

http://localhost:49152/v1.0/alive - GET Is the server alive?

http://localhost:49152/v1.0/parse - POST Submit a parse request to the server.
                                    Body should contain {"text":"Your text to parse. More text."}

http://localhost:49152/v1.0/status - GET Returns information about the server state.

## License

(The MIT License)

Copyright (c) 2015 Pyxxel Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to
deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
IN THE SOFTWARE.
