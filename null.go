package ejson

// IsNull判断是否是null值
func (g *JSON) IsNull() bool {
	raw := g.getRaw()
	return len(raw) == 4 && unsafeString(raw) == "null"
}
