package json2go

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"maps"
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
		convertObject(v.Object, w, nil)
	} else {
		w.Write([]byte("any"))
	}
}

func convertArray(a *jsonast.JsonArray, w io.Writer) {
	if len(a.Elements) == 0 {
		w.Write([]byte("[]any"))
		return
	}

	objs := []*jsonast.JsonObject{}

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

func convertObjectArray(a []*jsonast.JsonObject, w io.Writer) {
	union := orderedMapFrom(a[0])
	a = a[1:]
	omitempty := map[string]struct{}{}

	for _, obj := range a {
		m := orderedMapFrom(obj)

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
				union.Set(k, &jsonast.JsonValue{Array: &jsonast.JsonArray{}}) // []any
			} else {
				null := "null"
				union.Set(k, &jsonast.JsonValue{Null: &null}) // any
			}
		}

		maps.Copy(omitempty, union.XorKeys(m))
		union.WeakMerge(m)
	}

	w.Write([]byte("[]"))
	convertObject(union.Object(), w, omitempty)
}

func convertObject(obj *jsonast.JsonObject, w io.Writer, omitempty map[string]struct{}) {
	if omitempty == nil {
		omitempty = map[string]struct{}{}
	}

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
			fields[f] = 2
		}

		w.Write([]byte(" "))
		convert0(m.Value, w)
		w.Write([]byte(" `json:\""))
		w.Write([]byte(m.Key))
		if _, ok := omitempty[m.Key]; ok {
			w.Write([]byte(",omitempty"))
		}
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
