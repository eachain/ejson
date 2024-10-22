package ejson

import "strconv"

// IsBool判断是否是bool值，当值是true或false时返回true
func (g *JSON) IsBool() bool {
	raw := g.getRaw()
	return (len(raw) == 4 && unsafeString(raw) == "true") ||
		(len(raw) == 5 && unsafeString(raw) == "false")
}

// Bool尝试将值转为bool值，以下情况返回true:
//   - true
//   - "1", "t", "T", "TRUE", "true", "True"
//   - number != 0
//   - array or object len > 0
func (g *JSON) Bool() bool {
	b, _ := g.TryBool()
	return b
}

// Bool尝试将值转为bool值，并返回是否是bool值，以下情况返回true, true:
//   - true
//   - "1", "t", "T", "TRUE", "true", "True"
//   - number != 0
//
// 以下情况返回false, true:
//   - false
//   - "0", "f", "F", "FALSE", "false", "False"
//   - number == 0
//
// 以下情况返回true, false:
//   - string != ""
//   - array or object len > 0
//
// 以下情况返回false, false:
//   - null
//   - ""(empty string)
//   - array or object len == 0
func (g *JSON) TryBool() (bool, bool) {
	raw := g.getRaw()
	if len(raw) == 4 && unsafeString(raw) == "true" {
		return true, true
	}
	if len(raw) == 5 && unsafeString(raw) == "false" {
		return false, true
	}
	if g.IsNull() {
		return false, false
	}
	if g.IsStr() {
		s := g.Str()
		b, err := strconv.ParseBool(s)
		if err == nil {
			return b, true
		}
		return s != "", false
	}
	if g.IsNumber() {
		if val, ok := g.TryFloat(); ok {
			return val != 0, true
		}
		if val, ok := g.TryInt(); ok {
			return val != 0, true
		}
		if val, ok := g.TryUint(); ok {
			return val != 0, true
		}
	}
	return g.Len() != 0, false
}
