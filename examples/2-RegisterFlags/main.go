package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
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
	} `flag:""`

	// Nested fields with prefix
	Config struct {
		Timeout   time.Duration `flag:"timeout"`
		Endpoints []url.URL     `flag:"endpoints"`
	} `flag:"config-"`
}

func main() {
	spec := new(Spec)
	fs := flag.NewFlagSet("app", flag.ContinueOnError)

	if err := flagit.RegisterFlags(fs, spec, false); err != nil {
		panic(err)
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", spec)
}
