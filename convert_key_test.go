package json2go_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/json2go"
)

func TestNameToField(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "foo", expected: "Foo"},
		{input: "Foo", expected: "Foo"},
		{input: "fFoo", expected: "FFoo"},
		{input: "foo_bar-zoo", expected: "FooBarZoo"},
		{input: "f00", expected: "F00"},
		{input: "0foo", expected: "X_0Foo"},
		{input: "_0foo", expected: "X_0Foo"},
		{input: "_foo", expected: "Foo"},
		{input: "aaBBcc", expected: "AaBBcc"},
		{input: "_-", expected: ""},
		{input: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.input+"->"+tt.expected, func(t *testing.T) {
			assert := assert.New(t)
			field := json2go.ConvertKey(tt.input)
			assert.Equal(tt.expected, field)
		})
	}
}
