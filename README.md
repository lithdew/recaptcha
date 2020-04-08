# recaptcha

[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](LICENSE)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/lithdew/recaptcha)
[![Discord Chat](https://img.shields.io/discord/697002823123992617)](https://discord.gg/58dJzS)

**recaptcha** is a package that handles verifying reCAPTCHA v2/v3 submissions in Go.

- Validating a verification request, sending a verification request, and parsing a verification response are separated into individual functions.
- Validates and double-checks all outgoing requests and incoming responses for errors.
- Interoperable and made easy to work with both reCAPTCHA v2 and v3.
- Uses [valyala/fasthttp](https://github.com/valyala/fasthttp) for sending reCAPTCHA request with an optional timeout.
- Uses [valyala/fastjson](https://github.com/valyala/fastjson) for parsing responses from the reCAPTCHA API.

## Inspiration

Someone told me they were looking through reCAPTCHA packages online in Go and couldn't find a simple, idiomatic one.

This one's a bit overly optimized and uses two popular 3rd party libraries over the standard library, but here you go ¯\_(ツ)_/.

## Usage

```
go get github.com/lithdew/recaptcha
```

```go
package main

import (
    "github.com/lithdew/recaptcha"
    "time"
)

func main() {
    req := recaptcha.Request{
        Secret: "", // Your reCAPTCHA secret.
        Response: "", // The reCAPTCHA response sent by the reCAPTCHA API.
        RemoteIP : "", // (optional) The remote IP of the user submitting the reCAPTCHA response.
    }

    // Verify the reCAPTCHA request.
    
    res, err := recaptcha.Do(req) 
    if err != nil {
    	panic(err)
    }

    if res.Success {
        println("reCAPTCHA attempt successfully verified!")
    } else {
        println("reCAPTCHA attempt failed!")
    }

    // Verify the reCAPTCHA request, and timeout after 3 seconds.

    res, err = recaptcha.DoTimeout(req, 3 * time.Second) 
    if err != nil {
    	panic(err)
    }

    if res.Success {
        println("reCAPTCHA attempt successfully verified!")
    } else {
        println("reCAPTCHA attempt failed!")
    }
}
```

## Benchmarks

Take these with a grain of salt; network latency should sum up the majority of the benchmark results.

```
go test -bench=. -benchmem -benchtime=10s

goos: linux
goarch: amd64
pkg: github.com/lithdew/recaptcha
BenchmarkDo-8                        187          55273288 ns/op            1513 B/op         17 allocs/op
BenchmarkDoTimeout-8                 205          55503923 ns/op            1482 B/op         19 allocs/op
BenchmarkParallelDo-8               1500           7060534 ns/op            1386 B/op         17 allocs/op
BenchmarkParallelDoTimeout-8        1740           6752978 ns/op            1405 B/op         18 allocs/op
```