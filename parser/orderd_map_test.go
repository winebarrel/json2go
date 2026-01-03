package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/json2go/parser"
)

func TestOrderdMapGet(t *testing.T) {
	json := `{"str":"s","num":1,"t":true,"f":false,"null":null}`
	obj, err := parser.ParseJSON("", []byte(json))
	require.NoError(t, err)
	om := obj.Object.Map()

	m := map[string]*parser.JsonValue{
		"str":  {String: ptr("s")},
		"num":  {Number: ptr("1")},
		"t":    {True: ptr("true")},
		"f":    {False: ptr("false")},
		"null": {Null: ptr("null")},
	}

	for k, expected := range m {
		v, ok := om.Get(k)
		assert.True(t, ok)
		assert.Equal(t, expected, v)
	}

	_, ok := om.Get("invalid")
	assert.False(t, ok)
}

func TestOrderdMapEntries(t *testing.T) {
	json := `{"str":"s","num":1,"t":true,"f":false,"null":null}`
	obj, err := parser.ParseJSON("", []byte(json))
	require.NoError(t, err)
	om := obj.Object.Map()
	members := []*parser.JsonObjectMember{}

	for k, v := range om.Entries() {
		members = append(members, &parser.JsonObjectMember{Key: k, Value: v})
	}

	expected := []*parser.JsonObjectMember{
		{Key: "str", Value: &parser.JsonValue{String: ptr("s")}},
		{Key: "num", Value: &parser.JsonValue{Number: ptr("1")}},
		{Key: "t", Value: &parser.JsonValue{True: ptr("true")}},
		{Key: "f", Value: &parser.JsonValue{False: ptr("false")}},
		{Key: "null", Value: &parser.JsonValue{Null: ptr("null")}},
	}

	assert.Equal(t, expected, members)
}

func TestOrderedMapMerge(t *testing.T) {
	var om1, om2 *parser.OrderedMap

	{
		json := `{"str":"s1","num":1,"f":false}`
		obj, err := parser.ParseJSON("", []byte(json))
		require.NoError(t, err)
		om1 = obj.Object.Map()
	}

	{
		json := `{"str":"s2","t":true,"null":null}`
		obj, err := parser.ParseJSON("", []byte(json))
		require.NoError(t, err)
		om2 = obj.Object.Map()
	}

	om1.Merge(om2)
	members := []*parser.JsonObjectMember{}

	for k, v := range om1.Entries() {
		members = append(members, &parser.JsonObjectMember{Key: k, Value: v})
	}

	expected := []*parser.JsonObjectMember{
		{Key: "str", Value: &parser.JsonValue{String: ptr("s2")}},
		{Key: "num", Value: &parser.JsonValue{Number: ptr("1")}},
		{Key: "f", Value: &parser.JsonValue{False: ptr("false")}},
		{Key: "t", Value: &parser.JsonValue{True: ptr("true")}},
		{Key: "null", Value: &parser.JsonValue{Null: ptr("null")}},
	}

	assert.Equal(t, expected, members)
}

func TestOrderdMapObject(t *testing.T) {
	json := `{"str":"s","num":1,"t":true,"f":false,"null":null}`
	obj, err := parser.ParseJSON("", []byte(json))
	require.NoError(t, err)
	om := obj.Object.Map()
	expected := &parser.JsonObject{Members: []*parser.JsonObjectMember{
		{Key: "str", Value: &parser.JsonValue{String: ptr("s")}},
		{Key: "num", Value: &parser.JsonValue{Number: ptr("1")}},
		{Key: "t", Value: &parser.JsonValue{True: ptr("true")}},
		{Key: "f", Value: &parser.JsonValue{False: ptr("false")}},
		{Key: "null", Value: &parser.JsonValue{Null: ptr("null")}},
	}}
	assert.Equal(t, expected, om.Object())
}
