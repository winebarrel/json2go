package json2go_test

import (
	"os"
	"path/filepath"
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
			out, err := json2go.Convert([]byte(tt.input), false)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(out))
		})
	}
}

func TestConvert_Err(t *testing.T) {
	_, err := json2go.Convert([]byte("invalid"), false)
	require.ErrorContains(t, err, "failed to parse json:")
}

type testCase struct {
	Name     string
	Input    string
	Expected string
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

			if name == "" {
				name = input + "->" + expected
			}

			t.Run(f+"/"+name, func(t *testing.T) {
				out, err := json2go.Convert([]byte(input), true)
				require.NoError(t, err)
				assert.Equal(t, expected, string(out))
			})
		}
	}
}
