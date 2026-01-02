package json2go_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/json2go"
	"go.yaml.in/yaml/v4"
)

func TestConvert_Empty(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: strings.Repeat(" ", 100), expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.input+"->"+tt.expected, func(t *testing.T) {
			raw, err := json2go.Convert([]byte(tt.input), false)
			out := string(raw)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, out)
		})
	}
}

func TestConvert_Err(t *testing.T) {
	_, err := json2go.Convert([]byte("invalid"), false)
	require.ErrorContains(t, err, "failed to parse json:")
}

type testCase struct {
	Name      string
	Input     string
	Expected  string
	Unmarshal string
}

func TestConvert_OK(t *testing.T) {
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
			unmarshal := strings.TrimSpace(tt.Unmarshal)

			if name == "" {
				name = input + "->" + expected
			}

			t.Run(f+"/"+name, func(t *testing.T) {
				out, err := json2go.Convert([]byte(input), true)
				require.NoError(t, err)
				assert.Equal(t, expected, string(out))

				if testAcc && f != "testdata/primitive.yml" {
					x := compile(t, out)
					err := json.Unmarshal([]byte(input), x)
					require.NoError(t, err)
					assert.Equal(t, unmarshal, fmt.Sprintf("%+v", x))
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
