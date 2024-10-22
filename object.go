package ejson

import (
	"bytes"
	"encoding/json"
)

type object struct {
	keys   []string
	entry  map[string]*JSON
	parent *JSON
	update updateFuncs
}

// IsObject判断是否是object
func (g *JSON) IsObject() bool {
	if g.object != nil {
		return true
	}
	raw := g.getRaw()
	return len(raw) > 0 && raw[0] == '{'
}

func (g *JSON) asObject() bool {
	if g.object != nil {
		return true
	}

	dec := json.NewDecoder(bytes.NewReader(g.getRaw()))
	dec.UseNumber()

	leftDelim, err := dec.Token()
	if err != nil {
		return false
	}
	if d, ok := leftDelim.(json.Delim); !ok || d != '{' {
		return false
	}

	obj := new(object)
	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			return false
		}
		key, ok := token.(string)
		if !ok {
			return false
		}
		var val json.RawMessage
		err = dec.Decode(&val)
		if err != nil {
			return false
		}
		obj.keys = append(obj.keys, key)
		if obj.entry == nil {
			obj.entry = make(map[string]*JSON)
		}
		obj.entry[key] = &JSON{raw: val, parent: g}
	}

	rightDelim, err := dec.Token()
	if err != nil {
		return false
	}
	if d, ok := rightDelim.(json.Delim); !ok || d != '}' {
		return false
	}

	obj.parent = g
	g.object = obj

	return true
}

func (obj *object) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	buf.WriteByte('{')
	for i, key := range obj.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		err := enc.Encode(key)
		if err != nil {
			return nil, err
		}
		if buf.Bytes()[buf.Len()-1] == '\n' {
			buf.Truncate(buf.Len() - 1)
		}
		buf.WriteByte(':')
		err = enc.Encode(obj.entry[key])
		if err != nil {
			return nil, err
		}
		if buf.Bytes()[buf.Len()-1] == '\n' {
			buf.Truncate(buf.Len() - 1)
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (obj *object) Len() int {
	if obj == nil {
		return 0
	}
	return len(obj.keys)
}

func (obj *object) unmount() {
	if obj == nil {
		return
	}
	for _, g := range obj.entry {
		g.parent = nil
	}
	obj.parent = nil
}

// Keys返回object所有key列表
func (g *JSON) Keys() []string {
	g.asObject()
	if g.object == nil {
		return nil
	}

	keys := make([]string, len(g.object.keys))
	copy(keys, g.object.keys)
	return keys
}

func (obj *object) Index(key string) *JSON {
	if g := obj.entry[key]; g != nil {
		return g
	}

	g := &JSON{parent: obj.parent}

	parent := obj.parent
	g.update = obj.update.append(func() {
		if obj.parent != parent {
			return
		}

		if obj.entry == nil {
			obj.entry = make(map[string]*JSON)
		}
		if _, ok := obj.entry[key]; !ok {
			obj.keys = append(obj.keys, key)
		}
		obj.entry[key] = g

		for p := obj.parent; p != nil; p = p.parent {
			p.raw = nil
		}
	})

	return g
}

func (obj *object) Remove(g *JSON) {
	for key, j := range obj.entry {
		if j == g {
			delete(obj.entry, key)
			for i := 0; i < len(obj.keys); i++ {
				if obj.keys[i] == key {
					obj.keys = append(obj.keys[:i], obj.keys[i+1:]...)
					break
				}
			}
			break
		}
	}
}

// Any查询到任意key即返回。可用于err_code, errcode兼容的情况
func (g *JSON) Any(keys ...string) *JSON {
	if len(keys) == 0 {
		return new(JSON)
	}
	g.asObject()
	if g.object == nil {
		return new(JSON)
	}

	for _, key := range keys {
		if v := g.object.entry[key]; v != nil {
			return v
		}
	}
	return new(JSON)
}
