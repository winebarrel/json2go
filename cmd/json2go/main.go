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
	File    string `arg:"" optional:"" help:"JSON file. If not specified, read from stdin."`
	Version kong.VersionFlag
}

func parseArgs() *options {
	opts := &options{}
	parser := kong.Must(opts, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."

	if _, err := parser.Parse(os.Args[1:]); err != nil {
		parser.FatalIfErrorf(err)
	}

	return opts
}

func main() {
	opts := parseArgs()
	filename := opts.File
	var data []byte
	var err error

	if opts.File == "" || opts.File == "-" {
		filename = "<stdin>"
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(opts.File)
	}

	if err != nil {
		log.Fatal(err)
	}

	out, err := json2go.ConvertWithFilename(filename, data)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))
}
