package json2go_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/json2go/v2"
	"go.yaml.in/yaml/v4"
)

func TestConvertBytes_Empty(t *testing.T) {
	tests := []struct {
		input string
	}{
		{input: ""},
		{input: strings.Repeat(" ", 100)},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("SPACEx%d", len(tt.input)), func(t *testing.T) {
			_, err := json2go.ConvertBytes([]byte(tt.input))
			require.ErrorContains(t, err, `failed to parse json: 1:1: unexpected token "<EOF>"`)
		})
	}
}

func TestConvertBytes_Err(t *testing.T) {
	_, err := json2go.ConvertBytes([]byte("invalid"))
	require.ErrorContains(t, err, "failed to parse json:")
}

func TestConvertBytes_ErrWithFilename(t *testing.T) {
	_, err := json2go.ConvertBytes([]byte("{"), json2go.OptionFilename("example.json"))
	require.ErrorContains(t, err, `failed to parse json: example.json:1:2: unexpected token "<EOF>" (expected "}")`)
}

func TestConvert_OK(t *testing.T) {
	v, err := json2go.Convert(strings.NewReader("null"))
	require.NoError(t, err)
	assert.Equal(t, []byte(`any`), v)
}

type testCase struct {
	Name        string
	Input       string
	Expected    string
	Unmarshal   string
	Flat        bool
	NoOmitempty bool
	NoPointer   bool
}

func TestConvertBytes_OK(t *testing.T) {
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

			if name == "" {
				name = tt.Input + "->" + tt.Expected
			}

			t.Run(f+"/"+name, func(t *testing.T) {
				optfns := []json2go.OptFn{
					json2go.OptionFlat(tt.Flat),
					json2go.OptionOmitempty(!tt.NoOmitempty),
					json2go.OptionPointer(!tt.NoPointer),
				}

				out, err := json2go.ConvertBytes([]byte(tt.Input), optfns...)
				require.NoError(t, err)
				assert.Equal(t, tt.Expected, string(out))

				if testAcc && tt.Unmarshal != "skip" && !tt.Flat {
					x := compile(t, out)
					err := json.Unmarshal([]byte(tt.Input), x)
					require.NoError(t, err)

					if f == "testdata/primitive.yml" {
						x = reflect.ValueOf(x).Elem().Interface()
					}

					assert.Equal(t, tt.Unmarshal, fmt.Sprintf("%+v", x))
				}
			})
		}
	}
}

func compile(t *testing.T, src []byte) any {
	t.Helper()
	tmpdir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(tmpdir)
	data := fmt.Sprintf("package main\nvar A = *new(%s)", src)
	os.WriteFile("a.go", []byte(data), 0400)
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", "a.so", "a.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	require.NoError(t, err)
	plug, _ := plugin.Open("a.so")
	a, _ := plug.Lookup("A")
	return a
}
