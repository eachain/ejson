package ejson

import "strconv"

// IsNumber判断是否是数字类型
func (g *JSON) IsNumber() bool {
	raw := g.getRaw()
	if len(raw) == 0 {
		return false
	}

	if ('0' <= raw[0] && raw[0] <= '9') || raw[0] == '-' {
		return true
	}

	if raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw := raw[1 : len(raw)-1]
		if len(raw) == 0 {
			return false
		}
		_, err := strconv.ParseFloat(unsafeString(raw), 64)
		if err == nil {
			return true
		}
		_, err = strconv.ParseInt(unsafeString(raw), 10, 64)
		if err == nil {
			return true
		}
		_, err = strconv.ParseUint(unsafeString(raw), 10, 64)
		if err == nil {
			return true
		}
	}
	return false
}

// TryInt尝试将*JSON解析为int64值
func (g *JSON) TryInt() (int64, bool) {
	raw := g.getRaw()
	if len(raw) == 0 {
		return 0, false
	}
	if raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
		if len(raw) == 0 {
			return 0, false
		}
	}
	val, err := strconv.ParseInt(unsafeString(raw), 10, 64)
	if err == nil {
		return val, true
	}
	f, err := strconv.ParseFloat(unsafeString(raw), 64)
	return int64(f), err == nil
}

// IsInt判断是否是int值。
func (g *JSON) IsInt() bool {
	raw := g.getRaw()
	if len(raw) == 0 {
		return false
	}
	if raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
		if len(raw) == 0 {
			return false
		}
	}
	_, err := strconv.ParseInt(unsafeString(raw), 10, 64)
	return err == nil
}

// Int尝试将*JSON解析为int64值
func (g *JSON) Int() int64 {
	val, _ := g.TryInt()
	return val
}

// TryUint尝试将*JSON解析为uint64值
func (g *JSON) TryUint() (uint64, bool) {
	raw := g.getRaw()
	if len(raw) == 0 {
		return 0, false
	}
	if raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
		if len(raw) == 0 {
			return 0, false
		}
	}
	if raw[0] == '-' {
		return 0, false
	}
	val, err := strconv.ParseUint(unsafeString(raw), 10, 64)
	if err == nil {
		return val, true
	}
	f, err := strconv.ParseFloat(unsafeString(raw), 64)
	return uint64(f), err == nil
}

// IsUint判断是否是uint值。
func (g *JSON) IsUint() bool {
	raw := g.getRaw()
	if len(raw) == 0 {
		return false
	}
	if raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
		if len(raw) == 0 {
			return false
		}
	}
	if raw[0] == '-' {
		return false
	}
	_, err := strconv.ParseUint(unsafeString(raw), 10, 64)
	return err == nil
}

// Uint尝试将*JSON解析为uint64值
func (g *JSON) Uint() uint64 {
	val, _ := g.TryUint()
	return val
}

// TryFloat尝试将*JSON解析为float64值
func (g *JSON) TryFloat() (float64, bool) {
	raw := g.getRaw()
	if len(raw) == 0 {
		return 0, false
	}
	if raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
		if len(raw) == 0 {
			return 0, false
		}
	}
	val, err := strconv.ParseFloat(unsafeString(raw), 64)
	return val, err == nil
}

// IsFloat判断是否是float值。
func (g *JSON) IsFloat() bool {
	_, ok := g.TryFloat()
	return ok
}

// Float尝试将*JSON解析为float64值
func (g *JSON) Float() float64 {
	val, _ := g.TryFloat()
	return val
}
