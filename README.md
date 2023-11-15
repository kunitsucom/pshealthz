# [pshealthz](https://github.com/kunitsucom/pshealthz)

[![license](https://img.shields.io/github/license/kunitsucom/pshealthz)](LICENSE)
[![pkg](https://pkg.go.dev/badge/github.com/kunitsucom/pshealthz)](https://pkg.go.dev/github.com/kunitsucom/pshealthz)
[![goreportcard](https://goreportcard.com/badge/github.com/kunitsucom/pshealthz)](https://goreportcard.com/report/github.com/kunitsucom/pshealthz)
[![workflow](https://github.com/kunitsucom/pshealthz/workflows/go-lint/badge.svg)](https://github.com/kunitsucom/pshealthz/tree/main)
[![workflow](https://github.com/kunitsucom/pshealthz/workflows/go-test/badge.svg)](https://github.com/kunitsucom/pshealthz/tree/main)
[![workflow](https://github.com/kunitsucom/pshealthz/workflows/go-vuln/badge.svg)](https://github.com/kunitsucom/pshealthz/tree/main)
[![codecov](https://codecov.io/gh/kunitsucom/pshealthz/graph/badge.svg?token=8Jtk2bpTe2)](https://codecov.io/gh/kunitsucom/pshealthz)
[![sourcegraph](https://sourcegraph.com/github.com/kunitsucom/pshealthz/-/badge.svg)](https://sourcegraph.com/github.com/kunitsucom/pshealthz)

## Overview

`pshealthz` is a tool to check the health of a process.

## Example

```console
$ pshealthz
```

other terminal:

```console
$ curl -i http://localhost:8080 -d '{"regex":"^/usr/sbin/chronyd"}'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 15 Nov 2023 09:38:24 GMT
Content-Length: 74

{"ok":true,"processes":[{"pid":"10","cmdline":"/usr/sbin/chronyd -F -1"}]}
```

## Installation

### pre-built binary

```bash
VERSION=v0.0.1

# download
curl -fLROSs https://github.com/kunitsucom/pshealthz/releases/download/${VERSION}/pshealthz_${VERSION}_linux_amd64.zip

# unzip
unzip -j pshealthz_${VERSION}_linux_amd64.zip '*/pshealthz'
```

### go install

```bash
go install github.com/kunitsucom/pshealthz/cmd/pshealthz@latest
```

## Usage

```console
$ pshealthz --help
Usage:
    pshealthz [options]

Description:
    check process health via http

options:
    --version (default: false)
        show version information and exit
    --addr (env: PSHEALTHZ_ADDR, default: localhost:8888)
        listen address for http server
    --help (default: false)
        show usage
```
