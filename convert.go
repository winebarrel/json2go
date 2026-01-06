package json2go

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/winebarrel/jsonast"
)

func ConvertBytes(src []byte, optfns ...optionFunc) ([]byte, error) {
	f := func(filename string) (*jsonast.JsonValue, error) {
		return jsonast.ParseBytes(filename, src)
	}

	return convert(f, optfns...)
}

func Convert(r io.Reader, optfns ...optionFunc) ([]byte, error) {
	f := func(filename string) (*jsonast.JsonValue, error) {
		return jsonast.Parse(filename, r)
	}

	return convert(f, optfns...)
}

func convert(parse func(string) (*jsonast.JsonValue, error), optfns ...optionFunc) ([]byte, error) {
	options := &options{}

	for _, f := range optfns {
		f(options)
	}

	v, err := parse(options.filename)

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

func convert0(v *jsonast.JsonValue, w io.Writer) {
	if v.Nullable() {
		w.Write([]byte("*"))
	}

	switch value := v.Value().(type) {
	case *jsonast.JsonString:
		w.Write([]byte("string"))
	case *jsonast.JsonTrue, *jsonast.JsonFalse:
		w.Write([]byte("bool"))
	case *jsonast.JsonNumber:
		if strings.Contains(value.Text, ".") {
			w.Write([]byte("float64"))
		} else {
			w.Write([]byte("int"))
		}
	case *jsonast.JsonArray:
		convertArray(value, w)
	case *jsonast.JsonObject:
		convertObject(value, w)
	default:
		w.Write([]byte("any"))
	}
}

func convertArray(a *jsonast.JsonArray, w io.Writer) {
	if a.Len() == 0 {
		w.Write([]byte("[]any"))
		return
	}

	elem := a.UnionType(nil).Array.Elements[0]
	var base bytes.Buffer
	convert0(elem, &base)

	w.Write([]byte("[]"))
	w.Write(base.Bytes())
}

func convertObject(obj *jsonast.JsonObject, w io.Writer) {
	w.Write([]byte("struct {\n"))
	fields := map[string]int{}
	omittableKeys := obj.OmittableKeys

	if omittableKeys == nil {
		omittableKeys = map[string]struct{}{}
	}

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
			fields[f] = 2
		}

		w.Write([]byte(" "))
		convert0(m.Value, w)
		w.Write([]byte(" `json:"))
		tag := m.Key
		if _, ok := omittableKeys[m.Key]; ok {
			tag += ",omitempty"
		}
		w.Write([]byte(strconv.Quote(tag)))
		w.Write([]byte("`\n"))
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
