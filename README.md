# json2go

[![CI](https://github.com/winebarrel/json2go/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/json2go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/winebarrel/json2go)](https://goreportcard.com/report/github.com/winebarrel/json2go)

json2go is a tool to convert JSON to Go struct.

## Installation

```
go install github.com/winebarrel/json2go/cmd/json2go@latest
```

## Usage

```
Usage: json2go [<body-file>] [flags]

Arguments:
  [<body-file>]    JSON file. If not specified, read from stdin.

Flags:
  -h, --help       Show help.
      --version
```

```js
// example.json
{
  "glossary": {
    "title": "example glossary",
    "GlossDiv": {
      "title": "S",
      "GlossList": {
        "GlossEntry": {
          "ID": "SGML",
          "SortAs": "SGML",
          "GlossTerm": "Standard Generalized Markup Language",
          "Acronym": "SGML",
          "Abbrev": "ISO 8879:1986",
          "GlossDef": {
            "para": "A meta-markup language, used to create markup languages such as DocBook.",
            "GlossSeeAlso": [
              "GML",
              "XML"
            ]
          },
          "GlossSee": "markup"
        }
      }
    }
  }
}
```
```go
// json2go example.json # or `cat example.json | json2go`
struct {
	Glossary struct {
		GlossDiv struct {
			GlossList struct {
				GlossEntry struct {
					Abbrev   string `json:"Abbrev"`
					Acronym  string `json:"Acronym"`
					GlossDef struct {
						GlossSeeAlso []string `json:"GlossSeeAlso"`
						Para         string   `json:"para"`
					} `json:"GlossDef"`
					GlossSee  string `json:"GlossSee"`
					GlossTerm string `json:"GlossTerm"`
					ID        string `json:"ID"`
					SortAs    string `json:"SortAs"`
				} `json:"GlossEntry"`
			} `json:"GlossList"`
			Title string `json:"title"`
		} `json:"GlossDiv"`
		Title string `json:"title"`
	} `json:"glossary"`
}
```

### Use as a library

```go
package main

import (
	"fmt"
	"log"

	"github.com/winebarrel/json2go"
)

func main() {
	json := `{"foo":"bar","zoo":[100,200],"baz":{"hoge":"piyo"}}`
	gosrc, err := json2go.Convert([]byte(json), false)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(gosrc))
	//=> struct {
	//     Baz struct {
	//       Hoge string `json:"hoge"`
	//     } `json:"baz"`
	//     Foo string `json:"foo"`
	//     Zoo []int  `json:"zoo"`
	//   }
}
```
