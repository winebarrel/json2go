package json2go_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/winebarrel/json2go/v2"
)

func TestLargeJSON(t *testing.T) {
	if !testAcc {
		t.Skip()
	}

	urls := []string{
		"https://raw.githubusercontent.com/json-iterator/test-data/refs/heads/master/large-file.json",
		// NOTE: The test passes but is very slow, so comment it out
		// "https://raw.githubusercontent.com/seductiveapps/largeJSON/refs/heads/master/100mb.json",
	}

	for _, u := range urls {
		t.Run(u, func(t *testing.T) {
			resp, err := http.Get(u)
			require.NoError(t, err)
			t.Cleanup(func() {
				if resp != nil && resp.Body != nil {
					resp.Body.Close()
				}
			})
			var buf bytes.Buffer
			r := io.TeeReader(resp.Body, &buf)
			out, err := json2go.Convert(r)
			require.NoError(t, err)
			x := compile(t, out)
			err = json.Unmarshal(buf.Bytes(), x)
			require.NoError(t, err)
		})
	}
}
