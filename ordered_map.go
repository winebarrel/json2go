package json2go

import (
	"iter"

	"github.com/winebarrel/jsonast"
)

type orderedMap struct {
	keys       []string
	valueByKey map[string]*jsonast.JsonValue
}

func orderedMapFrom(obj *jsonast.JsonObject) *orderedMap {
	om := &orderedMap{
		keys:       []string{},
		valueByKey: map[string]*jsonast.JsonValue{},
	}

	for _, m := range obj.Members {
		om.Set(m.Key, m.Value)
	}

	return om
}

func (m *orderedMap) Set(key string, value *jsonast.JsonValue) {
	if _, ok := m.valueByKey[key]; !ok {
		m.keys = append(m.keys, key)
	}

	m.valueByKey[key] = value
}

func (m *orderedMap) Get(key string) (*jsonast.JsonValue, bool) {
	v, ok := m.valueByKey[key]
	return v, ok
}

func (m *orderedMap) Entries() iter.Seq2[string, *jsonast.JsonValue] {
	return func(yield func(string, *jsonast.JsonValue) bool) {
		for _, k := range m.keys {
			if !yield(k, m.valueByKey[k]) {
				return
			}
		}
	}
}

func (m *orderedMap) Merge(other *orderedMap) {
	for k, v := range other.Entries() {
		m.Set(k, v)
	}
}

func (m *orderedMap) Object() *jsonast.JsonObject {
	members := []*jsonast.JsonObjectMember{}

	for k, v := range m.Entries() {
		members = append(members, &jsonast.JsonObjectMember{Key: k, Value: v})
	}

	return &jsonast.JsonObject{Members: members}
}
