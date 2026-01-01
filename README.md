# json2go

[![CI](https://github.com/winebarrel/json2go/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/json2go/actions/workflows/ci.yml)

json2go is a tool to convert JSON to Go Struct.

## Usage

```
Usage: json2go [<body-file>] [flags]

Arguments:
  [<body-file>]    JSON file. If not specified, read from stdin.

Flags:
  -h, --help       Show help.
  -s, --sort       Sort struct fields.
      --version
```

```
$ cat example.json
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

$ json2go example.json # or `cat example.json | json2go`
struct {
	Glossary struct {
		Title    string `json:"title"`
		GlossDiv struct {
			Title     string `json:"title"`
			GlossList struct {
				GlossEntry struct {
					Abbrev   string `json:"Abbrev"`
					GlossDef struct {
						Para         string   `json:"para"`
						GlossSeeAlso []string `json:"GlossSeeAlso"`
					} `json:"GlossDef"`
					GlossSee  string `json:"GlossSee"`
					ID        string `json:"ID"`
					SortAs    string `json:"SortAs"`
					GlossTerm string `json:"GlossTerm"`
					Acronym   string `json:"Acronym"`
				} `json:"GlossEntry"`
			} `json:"GlossList"`
		} `json:"GlossDiv"`
	} `json:"glossary"`
}
```
