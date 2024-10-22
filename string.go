package ejson

import (
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

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
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
