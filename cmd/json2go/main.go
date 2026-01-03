package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/json2go/v2"
)

var version string

func init() {
	log.SetFlags(0)
}

type options struct {
	BodyFile kong.FileContentFlag `arg:"" optional:"" type:"filecontent" xor:"stdin" help:"JSON file. If not specified, read from stdin."`
	Version  kong.VersionFlag
}

func parseArgs() *options {
	opts := &options{}
	parser := kong.Must(opts, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	args := os.Args[1:]

	if _, err := parser.Parse(args); err != nil {
		parser.FatalIfErrorf(err)
	}

	if len(args) == 0 {
		if stdin, err := io.ReadAll(os.Stdin); err != nil {
			parser.FatalIfErrorf(err)
		} else {
			opts.BodyFile = stdin
		}
	}

	return opts
}

func main() {
	opts := parseArgs()
	out, err := json2go.Convert(opts.BodyFile)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))
}
