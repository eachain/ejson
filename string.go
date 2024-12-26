package ejson

import (
	"bytes"
	"encoding/json"
	"unicode/utf8"
	"unsafe"
)

// IsStr判断是否是字符串
func (g *JSON) IsStr() bool {
	raw := g.getRaw()
	return len(raw) > 0 && raw[0] == '"'
}

// Str返回字符串的值，如果*JSON不是字符串值，将返回空字符串。
// 如果需要返回原json串，请用String()方法。
func (g *JSON) Str() string {
	raw := g.getRaw()
	if t, ok := unquoteBytes(raw); ok {
		return unsafeString(t)
	}

	var s string
	json.Unmarshal(raw, &s)
	return s
}

// StrIsJSON判断字符串本身是不是个json，比如常见的数字字符串"123"等。
func (g *JSON) StrIsJSON() bool {
	return json.Valid(unsafeBytes(g.Str()))
}

type strjson struct {
	value  *JSON
	parent *JSON
	update updateFuncs
}

func (g *JSON) asStr() bool {
	if g.str != nil {
		return true
	}
	raw := unsafeBytes(g.Str())
	if json.Valid(raw) {
		g.str = &strjson{
			value:  FromBytes(raw),
			parent: g,
		}
		g.str.value.parent = g
		return true
	}
	return false
}

func (str *strjson) unmount() {
	if str == nil {
		return
	}
	if str.value != nil {
		str.value.parent = nil
	}
	str.parent = nil
}

func (str *strjson) Elem() *JSON {
	if g := str.value; g != nil {
		return g
	}

	g := &JSON{parent: str.parent}

	parent := str.parent
	g.update = str.update.append(func() {
		if str.parent != parent {
			return
		}
		str.value = g
	})

	return g
}

func (str *strjson) Remove(g *JSON) {
	if str.value == g {
		str.value = nil
	}
}

func (str *strjson) MarshalJSON() ([]byte, error) {
	raw, err := str.value.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return marshal(unsafeString(raw))
}

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func unsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func marshal(x any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(x)
	if err != nil {
		return nil, err
	}
	return bytes.TrimSpace(buf.Bytes()), nil
}

// copied from encoding/json/decode.go
func unquoteBytes(s []byte) (t []byte, ok bool) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return
	}
	s = s[1 : len(s)-1]

	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < len(s) {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == len(s) {
		return s, true
	}
	return nil, false
}
