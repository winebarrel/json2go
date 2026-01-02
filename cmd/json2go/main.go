package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/json2go"
)

var version string

func init() {
	log.SetFlags(0)
}

func parseArgs() ([]byte, *json2go.Options) {
	var cli struct {
		json2go.Options
		BodyFile kong.FileContentFlag `arg:"" optional:"" type:"filecontent" xor:"stdin" help:"JSON file. If not specified, read from stdin."`
		Version  kong.VersionFlag
	}

	parser := kong.Must(&cli, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	args := os.Args[1:]

	if _, err := parser.Parse(args); err != nil {
		parser.FatalIfErrorf(err)
	}

	if len(args) == 0 {
		if stdin, err := io.ReadAll(os.Stdin); err != nil {
			parser.FatalIfErrorf(err)
		} else {
			cli.BodyFile = stdin
		}
	}

	return cli.BodyFile, &cli.Options
}

func main() {
	src, options := parseArgs()
	out, err := json2go.Convert(src, options)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))
}
