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
	Verbose bool `flag:"verbose,enable verbose logs"`

	// Nested fields
	Options struct {
		Port     uint16 `flag:"port,the port number (1024-65535)"`
		LogLevel string `flag:"log-level,the logging level (debug|info|warn|error)"`
	} `flag:""`

	// Nested fields with prefix
	Config struct {
		Timeout   time.Duration `flag:"timeout,the request timeout"`
		Endpoints []url.URL     `flag:"endpoints,the replica endpoints"`
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
