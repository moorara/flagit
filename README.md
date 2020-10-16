[![Go Doc][godoc-image]][godoc-url]
[![Build Status][workflow-image]][workflow-url]
[![Go Report Card][goreport-image]][goreport-url]
[![Test Coverage][coverage-image]][coverage-url]
[![Maintainability][maintainability-image]][maintainability-url]

# flagit

flagit allows you to use the `flag` struct tag on your Go struct fields.
You can then read the values for those fields from the command-line arguments and parse them to the corresponding types of struct fields.

## Quick Start

```go
package main

import (
  "fmt"
  "net/url"
  "time"

  "github.com/moorara/flagit"
)

// Spec is a struct for mapping its fields to command-line flags.
type Spec struct {
  // Flag fields
  Verbose bool `flag:"verbose"`

  // Nested fields
  Options struct {
    Port     uint16 `flag:"port"`
    LogLevel string `flag:"log-level"`
  }

  // Nested fields with prefix
  Config struct {
    Timeout   time.Duration `flag:"timeout"`
    Endpoints []url.URL     `flag:"endpoints"`
  } `flag:"config-"`
}

func main() {
  spec := new(Spec)

  if err := flagit.Populate(spec, false); err != nil {
    panic(err)
  }

  fmt.Printf("%+v\n", spec)
}
```

## Examples

You can find more examples of using `flagit` [here](./examples).

## Documentation

### Supported Types

  - `bool`, `[]bool`
  - `string`, `[]string`
  - `float32`, `float64`
  - `[]float32`, `[]float64`
  - `int`, `int8`, `int16`, `int32`, `int64`
  - `[]int`, `[]int8`, `[]int16`, `[]int32`, `[]int64`
  - `uint`, `uint8`, `uint16`, `uint32`, `uint64`
  - `[]uint`, `[]uint8`, `[]uint16`, `[]uint32`, `[]uint64`
  - `time.Duration`, `[]time.Duration`
  - `url.URL`, `[]url.URL`
  - `regexp.Regexp`, `[]regexp.Regexp`

Nested structs are also supported.


[godoc-url]: https://pkg.go.dev/github.com/moorara/flagit
[godoc-image]: https://godoc.org/github.com/moorara/flagit?status.svg
[workflow-url]: https://github.com/moorara/flagit/actions
[workflow-image]: https://github.com/moorara/flagit/workflows/Main/badge.svg
[goreport-url]: https://goreportcard.com/report/github.com/moorara/flagit
[goreport-image]: https://goreportcard.com/badge/github.com/moorara/flagit
[coverage-url]: https://codeclimate.com/github/moorara/flagit/test_coverage
[coverage-image]: https://api.codeclimate.com/v1/badges/f441152938de958f7ebe/test_coverage
[maintainability-url]: https://codeclimate.com/github/moorara/flagit/maintainability
[maintainability-image]: https://api.codeclimate.com/v1/badges/f441152938de958f7ebe/maintainability
