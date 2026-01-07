package json2go_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/json2go/v2"
)

func TestBuffers(t *testing.T) {
	bs := json2go.NewBuffers()
	bs.Push(&bytes.Buffer{})
	bs.Write([]byte("foo"))
	bs.Push(&bytes.Buffer{})
	bs.Write([]byte("bar"))
	bs.Push(&bytes.Buffer{})
	bs.Write([]byte("zoo"))

	for _, expected := range []string{"zoo", "bar", "foo"} {
		b := bs.Pop()
		assert.Equal(t, expected, b.String())
	}

	assert.Nil(t, bs.Pop())
}
