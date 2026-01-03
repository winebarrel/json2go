package parser

import "github.com/alecthomas/participle/v2"

var (
	jsonParser = participle.MustBuild[JsonValue](
		participle.Lexer(&JsonDefinition{}),
	)
)

type JsonValue struct {
	False  *string     `parser:"@false |"`
	Null   *string     `parser:"@null |"`
	True   *string     `parser:"@true |"`
	Object *JsonObject `parser:"@@ |"`
	Array  *JsonArray  `parser:"@@ |"`
	Number *string     `parser:"@number |"`
	String *string     `parser:"@string"`
}

type JsonObject struct {
	Members []*JsonObjectMember `parser:"'{' @@* '}'"`
}

type JsonObjectMember struct {
	Key   string     `parser:"@string"`
	Value *JsonValue `parser:"@@"`
}

type JsonArray struct {
	Elements []*JsonValue `parser:"'[' @@* ']'"`
}

func ParseJSON(filename string, src []byte) (*JsonValue, error) {
	v, err := jsonParser.ParseBytes(filename, src)

	if err != nil {
		return nil, err
	}
	return v, nil
}
