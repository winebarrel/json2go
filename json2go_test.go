package json2go_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert/yaml"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/json2go"
)

func TestJsonToGo_Empty(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: strings.Repeat(" ", 100), expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.input+"->"+tt.expected, func(t *testing.T) {
			out, err := json2go.JsonToGo([]byte(tt.input), &json2go.Options{Sort: false})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(out))
		})
	}
}

type testCase struct {
	Name     string
	Input    string
	Expected string
}

func TestJsonToGo_OK(t *testing.T) {
	files, err := filepath.Glob("testdata/*.yml")

	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		yml, err := os.ReadFile(f)

		if err != nil {
			t.Fatal(err)
		}

		var tests []testCase
		err = yaml.Unmarshal(yml, &tests)

		if err != nil {
			t.Fatal(err)
		}

		for _, tt := range tests {
			name := tt.Name
			input := strings.TrimSpace(tt.Input)
			expected := strings.TrimSpace(tt.Expected)

			if name == "" {
				name = input + "->" + expected
			}

			t.Run(name, func(t *testing.T) {
				options := &json2go.Options{Sort: true}
				out, err := json2go.JsonToGo([]byte(input), options)
				require.NoError(t, err)
				assert.Equal(t, expected, string(out))
			})
		}
	}

	// 	tests := []struct {
	// 		input    string
	// 		expected string
	// 	}{
	// 		{input: "", expected: ""},
	// 		{input: strings.Repeat(" ", 100), expected: ""},

	// 		{input: "100", expected: "int"},
	// 		{input: `"hello"`, expected: "string"},
	// 		{input: "true", expected: "bool"},
	// 		{input: "false", expected: "bool"},
	// 		{input: "null", expected: "any"},

	// 		{input: "[1,2,3]", expected: "[]int"},
	// 		{input: `["h","e","l","l","o"]`, expected: "[]string"},
	// 		{input: "[true,false,true]", expected: "[]bool"},
	// 		{input: "[null]", expected: "[]any"},
	// 		{input: `[1,"2",true]`, expected: "[]any"},

	// 		{input: "[[100],[200]]", expected: "[][]int"},
	// 		{input: "[[[100]],[[200]]]", expected: "[][][]int"},
	// 		{input: `[["100"],["200"]]`, expected: "[][]string"},
	// 		{input: "[[true],[false]]", expected: "[][]bool"},
	// 		{input: `[[1],["2"],[true]]`, expected: "[]any"}, // NOTE: not a bug.

	// 		{
	// 			input: `{
	//     "glossary": {
	//         "title": "example glossary",
	// 		"GlossDiv": {
	//             "title": "S",
	// 			"GlossList": {
	//                 "GlossEntry": {
	//                     "ID": "SGML",
	// 					"SortAs": "SGML",
	// 					"GlossTerm": "Standard Generalized Markup Language",
	// 					"Acronym": "SGML",
	// 					"Abbrev": "ISO 8879:1986",
	// 					"GlossDef": {
	//                         "para": "A meta-markup language, used to create markup languages such as DocBook.",
	// 						"GlossSeeAlso": ["GML", "XML"]
	//                     },
	// 					"GlossSee": "markup"
	//                 }
	//             }
	//         }
	//     }
	// }`,
	// 			expected: `struct {
	//             	            	+	Glossary struct {
	//             	            	+		Title    string `json:"title"`
	//             	            	+		GlossDiv struct {
	//             	            	+			Title     string `json:"title"`
	//             	            	+			GlossList struct {
	//             	            	+				GlossEntry struct {
	//             	            	+					Acronym  string `json:"Acronym"`
	//             	            	+					Abbrev   string `json:"Abbrev"`
	//             	            	+					GlossDef struct {
	//             	            	+						Para         string   `json:"para"`
	//             	            	+						GlossSeeAlso []string `json:"GlossSeeAlso"`
	//             	            	+					} `json:"GlossDef"`
	//             	            	+					GlossSee  string `json:"GlossSee"`
	//             	            	+					ID        string `json:"ID"`
	//             	            	+					SortAs    string `json:"SortAs"`
	//             	            	+					GlossTerm string `json:"GlossTerm"`
	//             	            	+				} `json:"GlossEntry"`
	//             	            	+			} `json:"GlossList"`
	//             	            	+		} `json:"GlossDiv"`
	//             	            	+	} `json:"glossary"`
	//             	            	+}
	// `,
	// 		},
	// 	}

	//	for _, tt := range tests {
	//		t.Run(tt.input+"->"+tt.expected, func(t *testing.T) {
	//			out, err := json2go.JsonToGo([]byte(tt.input))
	//			require.NoError(t, err)
	//			assert.Equal(t, tt.expected, string(out))
	//		})
	//	}
}
