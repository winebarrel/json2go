package json2go_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/json2go/v2"
	"github.com/winebarrel/jsonast"
)

func TestOrderedMapGet(t *testing.T) {
	json := `{"str":"s","num":1,"t":true,"f":false,"null":null}`
	obj, err := jsonast.ParseBytes("", []byte(json))
	require.NoError(t, err)
	om := json2go.OrderedMapFrom(obj.Object)

	m := map[string]*jsonast.JsonValue{
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

func TestOrderedMapEntries(t *testing.T) {
	json := `{"str":"s","num":1,"t":true,"f":false,"null":null}`
	obj, err := jsonast.ParseBytes("", []byte(json))
	require.NoError(t, err)
	om := json2go.OrderedMapFrom(obj.Object)
	members := []*jsonast.JsonObjectMember{}

	for k, v := range om.Entries() {
		members = append(members, &jsonast.JsonObjectMember{Key: k, Value: v})
	}

	expected := []*jsonast.JsonObjectMember{
		{Key: "str", Value: &jsonast.JsonValue{String: ptr("s")}},
		{Key: "num", Value: &jsonast.JsonValue{Number: ptr("1")}},
		{Key: "t", Value: &jsonast.JsonValue{True: ptr("true")}},
		{Key: "f", Value: &jsonast.JsonValue{False: ptr("false")}},
		{Key: "null", Value: &jsonast.JsonValue{Null: ptr("null")}},
	}

	assert.Equal(t, expected, members)
}

func TestOrderedMapWeakMerge(t *testing.T) {
	var om1, om2 *json2go.OrderedMap

	{
		json := `{"str":"s1","num":1,"f":false,"x":null}`
		obj, err := jsonast.ParseBytes("", []byte(json))
		require.NoError(t, err)
		om1 = json2go.OrderedMapFrom(obj.Object)
	}

	{
		json := `{"str":"s2","t":true,"null":null,"x":100}`
		obj, err := jsonast.ParseBytes("", []byte(json))
		require.NoError(t, err)
		om2 = json2go.OrderedMapFrom(obj.Object)
	}

	om1.WeakMerge(om2)
	members := []*jsonast.JsonObjectMember{}

	for k, v := range om1.Entries() {
		members = append(members, &jsonast.JsonObjectMember{Key: k, Value: v})
	}

	expected := []*jsonast.JsonObjectMember{
		{Key: "str", Value: &jsonast.JsonValue{String: ptr("s1")}},
		{Key: "num", Value: &jsonast.JsonValue{Number: ptr("1")}},
		{Key: "f", Value: &jsonast.JsonValue{False: ptr("false")}},
		{Key: "x", Value: &jsonast.JsonValue{Null: ptr("null")}},
		{Key: "t", Value: &jsonast.JsonValue{True: ptr("true")}},
		{Key: "null", Value: &jsonast.JsonValue{Null: ptr("null")}},
	}
	assert.Equal(t, expected, members)
}

func TestOrderedMapObject(t *testing.T) {
	json := `{"str":"s","num":1,"t":true,"f":false,"null":null}`
	obj, err := jsonast.ParseBytes("", []byte(json))
	require.NoError(t, err)
	om := json2go.OrderedMapFrom(obj.Object)
	expected := &jsonast.JsonObject{Members: []*jsonast.JsonObjectMember{
		{Key: "str", Value: &jsonast.JsonValue{String: ptr("s")}},
		{Key: "num", Value: &jsonast.JsonValue{Number: ptr("1")}},
		{Key: "t", Value: &jsonast.JsonValue{True: ptr("true")}},
		{Key: "f", Value: &jsonast.JsonValue{False: ptr("false")}},
		{Key: "null", Value: &jsonast.JsonValue{Null: ptr("null")}},
	}}
	assert.Equal(t, expected, om.Object())
}

func TestOrderedMapXorKeys(t *testing.T) {
	var om1, om2 *json2go.OrderedMap

	{
		json := `{"str":"s1","num":1,"f":false,"x":null}`
		obj, err := jsonast.ParseBytes("", []byte(json))
		require.NoError(t, err)
		om1 = json2go.OrderedMapFrom(obj.Object)
	}

	{
		json := `{"str":"s2","t":true,"null":null,"x":100}`
		obj, err := jsonast.ParseBytes("", []byte(json))
		require.NoError(t, err)
		om2 = json2go.OrderedMapFrom(obj.Object)
	}

	expected := map[string]struct{}{"f": {}, "null": {}, "num": {}, "t": {}}
	assert.Equal(t, expected, om1.XorKeys(om2))
	assert.Equal(t, expected, om2.XorKeys(om1))
}
