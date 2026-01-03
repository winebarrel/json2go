package parser

import "iter"

type OrderedMap struct {
	keys       []string
	valueByKey map[string]*JsonValue
}

func newOrderedMap() *OrderedMap {
	return &OrderedMap{
		keys:       []string{},
		valueByKey: map[string]*JsonValue{},
	}
}

func (m *OrderedMap) Set(key string, value *JsonValue) {
	if _, ok := m.valueByKey[key]; !ok {
		m.keys = append(m.keys, key)
	}

	m.valueByKey[key] = value
}

func (m *OrderedMap) Get(key string) (*JsonValue, bool) {
	v, ok := m.valueByKey[key]
	return v, ok
}

func (m *OrderedMap) Entries() iter.Seq2[string, *JsonValue] {
	return func(yield func(string, *JsonValue) bool) {
		for _, k := range m.keys {
			if !yield(k, m.valueByKey[k]) {
				return
			}
		}
	}
}

func (m *OrderedMap) Merge(other *OrderedMap) {
	for k, v := range other.Entries() {
		m.Set(k, v)
	}
}

func (m *OrderedMap) Object() *JsonObject {
	members := []*JsonObjectMember{}

	for k, v := range m.Entries() {
		members = append(members, &JsonObjectMember{Key: k, Value: v})
	}

	return &JsonObject{Members: members}
}
