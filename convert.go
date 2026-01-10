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

func ConvertBytes(src []byte, optfns ...OptFn) ([]byte, error) {
	buf := bytes.NewBuffer(src)
	return Convert(buf, optfns...)
}

func Convert(r io.Reader, optfns ...OptFn) ([]byte, error) {
	f := func(filename string) (*jsonast.JsonValue, error) {
		return jsonast.Parse(filename, r)
	}

	return convert(f, optfns...)
}

func convert(parse func(string) (*jsonast.JsonValue, error), optfns ...OptFn) ([]byte, error) {
	options := newOptions(optfns)
	v, err := parse(options.filename)

	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	var buf bytes.Buffer
	bufs := []*bytes.Buffer{&buf}

	c := &converter{
		opts: options,
		bufs: bufs,
	}

	c.convertAny(v, &buf)
	var out bytes.Buffer

	if options.typeName != "" {
		out.WriteString("type ")
		out.WriteString(options.typeName)
		out.WriteString(" ")
	}

	for i, b := range c.bufs {
		o, err := format.Source(b.Bytes())

		if err != nil {
			return nil, fmt.Errorf("failed to format golang: %w", err)
		}

		if i > 0 {
			out.WriteByte('\n')
		}

		out.Write(o)
	}

	return out.Bytes(), nil
}

type converter struct {
	opts *options
	bufs []*bytes.Buffer
}

func (c *converter) convertAny(v *jsonast.JsonValue, w *bytes.Buffer) {
	if c.opts.pointer && v.Nullable() {
		w.WriteString("*")
	}

	switch value := v.Value().(type) {
	case *jsonast.JsonString:
		w.WriteString("string")
	case *jsonast.JsonTrue, *jsonast.JsonFalse:
		w.WriteString("bool")
	case *jsonast.JsonNumber:
		if strings.Contains(value.Text, ".") {
			w.WriteString("float64")
		} else {
			w.WriteString("int")
		}
	case *jsonast.JsonArray:
		c.convertArray(value, w)
	case *jsonast.JsonObject:
		c.convertObject(value, w)
	default:
		w.WriteString("any")
	}
}

func (c *converter) convertArray(a *jsonast.JsonArray, w *bytes.Buffer) {
	if a.Len() == 0 {
		w.WriteString("[]any")
		return
	}

	elem := a.UnionType(nil).Array.Elements[0]
	var base bytes.Buffer
	c.convertAny(elem, &base)

	w.WriteString("[]")
	w.Write(base.Bytes())
}

func (c *converter) convertObject(obj *jsonast.JsonObject, w *bytes.Buffer) {
	w.WriteString("struct {\n")
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

		num, ok := fields[f]

		if ok {
			fields[f] = num + 1
			f += strconv.Itoa(num)
		} else {
			fields[f] = 2
		}

		w.WriteString(f)
		w.WriteString(" ")

		var worig *bytes.Buffer

		if c.opts.flat && m.Value.IsObject() {
			w.WriteString(f)
			worig = w
			b := &bytes.Buffer{}
			b.WriteString("type ")
			b.WriteString(f)
			b.WriteString(" ")
			c.bufs = append(c.bufs, b)
			w = b
		}

		c.convertAny(m.Value, w)

		if worig != nil {
			w = worig
		}

		w.WriteString(" `json:")
		tag := m.Key
		if c.opts.omitempty {
			if _, ok := omittableKeys[m.Key]; ok {
				tag += ",omitempty"
			}
		}
		w.WriteString(strconv.Quote(tag))
		w.WriteString("`\n")
	}

	w.WriteString("}")
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
