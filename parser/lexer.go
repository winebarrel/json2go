package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/alecthomas/participle/v2/lexer"
)

const (
	TokenTypeDelim  lexer.TokenType = iota // '[',']','{','}'
	TokenTypeFalse                         // false
	TokenTypeNull                          // null
	TokenTypeTrue                          // true
	TokenTypeNumber                        // number
	TokenTypeString                        // string
)

var jsonSymbols = map[string]lexer.TokenType{
	"[":      TokenTypeDelim,
	"]":      TokenTypeDelim,
	"{":      TokenTypeDelim,
	"}":      TokenTypeDelim,
	"false":  TokenTypeFalse,
	"null":   TokenTypeNull,
	"true":   TokenTypeTrue,
	"number": TokenTypeNumber,
	"string": TokenTypeString,
}

type JsonDefinition struct {
}

func (l *JsonDefinition) Symbols() map[string]lexer.TokenType {
	return jsonSymbols
}

func (l *JsonDefinition) Lex(filename string, r io.Reader) (lexer.Lexer, error) {
	buf := &bytes.Buffer{}
	decoder := json.NewDecoder(io.TeeReader(r, buf))
	decoder.UseNumber()

	lex := &JsonLexer{
		decoder: decoder,
		buf:     buf,
		pos: lexer.Position{
			Filename: filename,
			Line:     1,
			Column:   1,
		},
	}

	return lex, nil
}

type JsonLexer struct {
	decoder *json.Decoder
	buf     *bytes.Buffer
	pos     lexer.Position
}

func (l *JsonLexer) Next() (lexer.Token, error) {
	startOffset := l.decoder.InputOffset()
	rawTok, err := l.decoder.Token()
	span := make([]byte, l.decoder.InputOffset()-startOffset)
	tok := lexer.Token{}

	if _, err := l.buf.Read(span); err != nil {
		return tok, err
	}

	tok.Pos = l.pos
	l.pos.Advance(string(span))

	if err == io.EOF {
		tok.Type = lexer.EOF
		return tok, nil
	} else if err != nil {
		return tok, fmt.Errorf("%d:%d: %w", tok.Pos.Line, tok.Pos.Column, err)
	}

	switch v := rawTok.(type) {
	case json.Delim:
		tok.Type = TokenTypeDelim
		tok.Value = v.String()
	case bool:
		if v {
			tok.Type = TokenTypeTrue
			tok.Value = "true"
		} else {
			tok.Type = TokenTypeFalse
			tok.Value = "false"
		}
	case nil:
		tok.Type = TokenTypeNull
		tok.Value = "null"
	case json.Number:
		tok.Type = TokenTypeNumber
		tok.Value = v.String()
	case string:
		tok.Type = TokenTypeString
		tok.Value = v
	}

	return tok, nil
}
