package ejson

import (
	"bytes"
	"encoding/json"
	"unsafe"
)

type JSON struct {
	raw    json.RawMessage
	object *object
	array  *array
	parent *JSON
	update updateFuncs
}

type updateFuncs []func()

func (fs updateFuncs) run() {
	for _, f := range fs {
		f()
	}
}

func (fs updateFuncs) append(f func()) updateFuncs {
	return append(fs[:len(fs):len(fs)], f)
}

// Valid reports whether data is a valid JSON encoding.
func Valid(b []byte) bool {
	return json.Valid(b)
}

func FromBytes(b []byte) *JSON {
	return &JSON{raw: bytes.TrimSpace(b)}
}

func FromString(s string) *JSON {
	return FromBytes(unsafe.Slice(unsafe.StringData(s), len(s)))
}

// Len返回string、array、object长度，其它类型值返回0
func (g *JSON) Len() int {
	if g.IsNull() {
		return 0
	}
	if g.IsStr() {
		return len(g.Str())
	}
	if g.IsObject() {
		g.asObject()
		return g.object.Len()
	}
	if g.IsArray() {
		g.asArray()
		return g.array.Len()
	}
	return 0
}

// ArrayIndex返回array第i个值。如果不存在，返回*JSON不为空。
// 当返回的*JSON被写入数据后，返回的*JSON将写入当前array。
func (g *JSON) ArrayIndex(i int) *JSON {
	g.asArray()
	a := g.array
	if a == nil {
		a = &array{parent: g}
		a.update = g.update.append(func() {
			if g.array != a {
				g.array.unmount()
				g.array = a
				g.object.unmount()
				g.object = nil
			}
		})
	}
	return a.Index(i)
}

// ObjectIndex返回object[key]的值。如果不存在，返回*JSON不为空。
// 当返回的*JSON被写入数据后，返回的*JSON将写入当前object。
func (g *JSON) ObjectIndex(key string) *JSON {
	g.asObject()
	obj := g.object
	if obj == nil {
		obj = &object{parent: g}
		obj.update = g.update.append(func() {
			if g.object != obj {
				g.object.unmount()
				g.object = obj
				g.array.unmount()
				g.array = nil
			}
		})
	}
	return obj.Index(key)
}

// Remove将当前*JSON从json树中移除
func (g *JSON) Remove() {
	if g.parent == nil {
		return
	}
	if g.parent.object != nil {
		g.parent.object.Remove(g)
	}
	if g.parent.array != nil {
		g.parent.array.Remove(g)
	}

	for p := g.parent; p != nil; p = p.parent {
		p.raw = nil
	}
	g.parent = nil
}

// String实现fmt.Stringer接口。
// 当前*JSON是字符串时，返回结果和(*JSON).Str()相同；
// 否则返回当前*JSON原始json字符串。
func (g *JSON) String() string {
	if g.IsStr() {
		return g.Str()
	}
	return string(g.getRaw())
}

// UnsafeString返回当前*JSON原始json不安全的字符串。
func (g *JSON) UnsafeString() string {
	return unsafeString(g.getRaw())
}

// Exists返回当前*JSON是否存在。
func (g *JSON) Exists() bool {
	return len(g.getRaw()) > 0
}

// Value将当前*JSON解析到v
func (g *JSON) Value(v any) error {
	return json.Unmarshal(g.getRaw(), v)
}

// UnmarshalJSON实现json.Unmarshaler接口，可用于写入*JSON。
func (g *JSON) UnmarshalJSON(raw []byte) error {
	g.reset(raw)
	return nil
}

// MarshalJSON实现json.Marshaler接口。
func (g *JSON) MarshalJSON() ([]byte, error) {
	if g == nil {
		return []byte("null"), nil
	}

	raw := g.getRaw()
	if len(raw) == 0 {
		return []byte("null"), nil
	}
	return raw, nil
}

func (g *JSON) getRaw() []byte {
	if len(g.raw) == 0 {
		if g.array != nil {
			g.raw, _ = g.array.MarshalJSON()
		} else if g.object != nil {
			g.raw, _ = g.object.MarshalJSON()
		}
	}

	return g.raw
}

// Set设置当前*JSON值。
func (g *JSON) Set(v any) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	g.reset(raw)
	return nil
}

func (g *JSON) reset(raw json.RawMessage) {
	if bytes.Equal(g.getRaw(), raw) {
		return
	}

	g.raw = raw

	if g.object != nil {
		g.object.unmount()
		g.object = nil
	}
	if g.array != nil {
		g.array.unmount()
		g.array = nil
	}

	for p := g.parent; p != nil; p = p.parent {
		p.raw = nil
	}

	g.update.run()
}

// Get支持类型'[0].object.key.array[0][1].key'式取值。
func (g *JSON) Get(smartKey string) *JSON {
	es, err := parseSmartKey(smartKey)
	if err != nil {
		return new(JSON)
	}
	for _, e := range es {
		g = e(g)
	}
	return g
}
