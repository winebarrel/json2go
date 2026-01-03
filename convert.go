package json2go

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/winebarrel/json2go/parser"
)

func Convert(src []byte) ([]byte, error) {
	return ConvertWithFilename("", src)
}

func ConvertWithFilename(filename string, src []byte) ([]byte, error) {
	if len(bytes.TrimSpace(src)) == 0 {
		return []byte{}, nil
	}

	v, err := parser.ParseJSON(filename, src)

	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	var buf bytes.Buffer
	convert0(v, &buf)
	out, err := format.Source(buf.Bytes())

	if err != nil {
		return nil, fmt.Errorf("failed to format golang: %w", err)
	}

	return out, nil
}

func convert0(v *parser.JsonValue, w io.Writer) {
	if v.String != nil {
		w.Write([]byte("string"))
	} else if v.False != nil || v.True != nil {
		w.Write([]byte("bool"))
	} else if v.Number != nil {
		if strings.Contains(*v.Number, ".") {
			w.Write([]byte("float64"))
		} else {
			w.Write([]byte("int"))
		}
	} else if v.Array != nil {
		convertArray(v.Array, w)
	} else if v.Object != nil {
		convertObject(v.Object, w)
	} else {
		w.Write([]byte("any"))
	}
}

func convertArray(a *parser.JsonArray, w io.Writer) {
	if len(a.Elements) == 0 {
		w.Write([]byte("[]any"))
		return
	}

	objs := []*parser.JsonObject{}

	for _, x := range a.Elements {
		if x.Object != nil {
			objs = append(objs, x.Object)
		}
	}

	if len(a.Elements) == len(objs) {
		convertObjectArray(objs, w)
		return
	}

	w.Write([]byte("[]"))
	base := ""

	for _, x := range a.Elements {
		var buf bytes.Buffer
		convert0(x, &buf)
		t := buf.String()

		if t == "any" || (base != "" && base != t) {
			base = "any"
			break
		}

		base = t
	}

	w.Write([]byte(base))
}

func convertObjectArray(a []*parser.JsonObject, w io.Writer) {
	union := a[0].Map()
	a = a[1:]

	for _, obj := range a {
		m := obj.Map()

		for k, v := range m.Entries() {
			uv, ok := union.Get(k)

			if !ok {
				continue
			}

			var buf1, buf2 bytes.Buffer
			convert0(uv, &buf1)
			convert0(v, &buf2)
			t1 := buf1.String()
			t2 := buf2.String()

			if t1 == t2 {
				continue
			}

			if strings.HasPrefix(t1, "[]") && strings.HasPrefix(t2, "[]") {
				union.Set(k, &parser.JsonValue{Array: &parser.JsonArray{}}) // []any
			} else {
				null := "null"
				union.Set(k, &parser.JsonValue{Null: &null}) // []any
			}
		}

		m.Merge(union)
		union = m
	}

	w.Write([]byte("[]"))
	convertObject(union.Object(), w)
}

func convertObject(obj *parser.JsonObject, w io.Writer) {
	w.Write([]byte("struct {\n"))
	fields := map[string]int{}

	for _, m := range obj.Members {
		f := convertKey(m.Key)

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
		convert0(m.Value, w)
		w.Write([]byte(" `json:\""))
		w.Write([]byte(m.Key))
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
				buf.WriteString("Num")
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
