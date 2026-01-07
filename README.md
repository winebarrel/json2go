# json2go

[![CI](https://github.com/winebarrel/json2go/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/json2go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/winebarrel/json2go/v2)](https://goreportcard.com/report/github.com/winebarrel/json2go/v2)

json2go is a tool to convert JSON to Go struct.

## Installation

```
brew install winebarrel/json2go/json2go
```

## Usage

```
Usage: json2go [<file>] [flags]

Arguments:
  [<file>]    JSON file. If not specified, read from stdin.

Flags:
  -h, --help              Show help.
      --[no-]flat         Flattening structs.
      --[no-]omitempty    Add 'omitempty' to optional fields. (default: true)
      --[no-]pointer      Make nullable fields pointer types. (default: true)
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
		Title    string `json:"title"`
		GlossDiv struct {
			Title     string `json:"title"`
			GlossList struct {
				GlossEntry struct {
					ID        string `json:"ID"`
					SortAs    string `json:"SortAs"`
					GlossTerm string `json:"GlossTerm"`
					Acronym   string `json:"Acronym"`
					Abbrev    string `json:"Abbrev"`
					GlossDef  struct {
						Para         string   `json:"para"`
						GlossSeeAlso []string `json:"GlossSeeAlso"`
					} `json:"GlossDef"`
					GlossSee string `json:"GlossSee"`
				} `json:"GlossEntry"`
			} `json:"GlossList"`
		} `json:"GlossDiv"`
	} `json:"glossary"`
}
```

### `--flat` option

```go
// json2go --flat example.json
type Root struct {
	Glossary Glossary `json:"glossary"`
}
type Glossary struct {
	Title    string   `json:"title"`
	GlossDiv GlossDiv `json:"GlossDiv"`
}
type GlossDiv struct {
	Title     string    `json:"title"`
	GlossList GlossList `json:"GlossList"`
}
type GlossList struct {
	GlossEntry GlossEntry `json:"GlossEntry"`
}
type GlossEntry struct {
	ID        string   `json:"ID"`
	SortAs    string   `json:"SortAs"`
	GlossTerm string   `json:"GlossTerm"`
	Acronym   string   `json:"Acronym"`
	Abbrev    string   `json:"Abbrev"`
	GlossDef  GlossDef `json:"GlossDef"`
	GlossSee  string   `json:"GlossSee"`
}
type GlossDef struct {
	Para         string   `json:"para"`
	GlossSeeAlso []string `json:"GlossSeeAlso"`
}
```

### Use as a library

```go
package main

import (
	"fmt"
	"log"

	"github.com/winebarrel/json2go/v2"
)

func main() {
	json := `{"foo":"bar","zoo":[100,200],"baz":{"hoge":"piyo"}}`
	gosrc, err := json2go.ConvertBytes([]byte(json))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(gosrc))
	//=> struct {
	//     Foo string `json:"foo"`
	//     Zoo []int  `json:"zoo"`
	//     Baz struct {
	//       Hoge string `json:"hoge"`
	//     } `json:"baz"`
	//   }
}
```
