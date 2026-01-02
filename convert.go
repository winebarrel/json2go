package json2go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

func Convert(src []byte, sort bool) ([]byte, error) {
	if len(bytes.TrimSpace(src)) == 0 {
		return []byte{}, nil
	}

	decoder := json.NewDecoder(bytes.NewReader(src))
	decoder.UseNumber()
	var x any
	err := decoder.Decode(&x)

	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	var buf bytes.Buffer
	convert0(x, &buf, sort)
	out, err := format.Source(buf.Bytes())

	if err != nil {
		return nil, fmt.Errorf("failed to format golang: %w", err)
	}

	return out, nil
}

func convert0(x any, w io.Writer, sort bool) {
	switch v := x.(type) {
	case string:
		w.Write([]byte("string"))
	case bool:
		w.Write([]byte("bool"))
	case json.Number:
		if strings.Contains(v.String(), ".") {
			w.Write([]byte("float64"))
		} else {
			w.Write([]byte("int"))
		}
	case []any:
		convertArray(v, w)
	case map[string]any:
		convertObject(v, w, sort)
	default:
		w.Write([]byte("any"))
	}
}

func convertArray(a []any, w io.Writer) {
	w.Write([]byte("[]"))
	base := ""

	for _, x := range a {
		var buf bytes.Buffer
		convert0(x, &buf, true)
		t := buf.String()

		if t == "any" || (base != "" && base != t) {
			base = "any"
			break
		}

		base = t
	}

	w.Write([]byte(base))
}

func convertObject(m map[string]any, w io.Writer, sort bool) {
	w.Write([]byte("struct {\n"))
	fields := map[string]int{}
	keys := []string{}

	for k := range m {
		keys = append(keys, k)
	}

	if sort {
		slices.Sort(keys)
	}

	for _, n := range keys {
		f := convertKey(n)

		if f == "" {
			f = "NAMING_FIELD"
		}

		w.Write([]byte(f))
		num, ok := fields[f]

		if ok {
			w.Write([]byte(strconv.Itoa(num)))
			fields[f] = num + 1
		} else {
			fields[f] = 0
		}

		w.Write([]byte(" "))
		convert0(m[n], w, sort)
		w.Write([]byte(" `json:\""))
		w.Write([]byte(n))
		w.Write([]byte("\"`\n"))
	}

	w.Write([]byte("}"))
}

func convertKey(key string) string {
	var buf strings.Builder
	boundary := true

	for _, r := range key {
		if rune('0') <= r && r <= rune('9') {
			boundary = true

			if buf.Len() == 0 {
				buf.WriteString("X_")
			}

			buf.WriteRune(r)
		} else if rune('a') <= r && r <= rune('z') {
			if boundary {
				r = unicode.ToUpper(r)
			}

			boundary = false
			buf.WriteRune(r)
		} else if rune('A') <= r && r <= rune('Z') {
			boundary = false
			buf.WriteRune(r)
		} else {
			boundary = true
		}
	}

	return buf.String()
}
