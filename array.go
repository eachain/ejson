package ejson

import "encoding/json"

type array struct {
	values []*JSON
	parent *JSON
	update updateFuncs
}

// IsArray判断是否是数组
func (g *JSON) IsArray() bool {
	if g.array != nil {
		return true
	}
	raw := g.getRaw()
	return len(raw) > 0 && raw[0] == '['
}

func (g *JSON) asArray() bool {
	if g.array != nil {
		return true
	}

	var elems []json.RawMessage
	err := json.Unmarshal(g.getRaw(), &elems)
	if err != nil {
		return false
	}
	g.array = &array{parent: g}
	if len(elems) > 0 {
		g.array.values = make([]*JSON, len(elems))
		for i := range g.array.values {
			g.array.values[i] = &JSON{raw: elems[i], parent: g}
		}
	}
	return true
}

func (a *array) MarshalJSON() ([]byte, error) {
	var err error
	raws := make([]json.RawMessage, len(a.values))
	for i, g := range a.values {
		raws[i], err = g.MarshalJSON()
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(raws)
}

func (a *array) Len() int {
	if a == nil {
		return 0
	}
	return len(a.values)
}

func (a *array) unmount() {
	if a == nil {
		return
	}
	for _, g := range a.values {
		if g != nil {
			g.parent = nil
		}
	}
	a.parent = nil
}

func (a *array) Index(i int) *JSON {
	if 0 <= i && i < len(a.values) {
		return a.values[i]
	}
	if i < 0 && len(a.values) > 0 {
		return a.values[i%len(a.values)+len(a.values)]
	}

	g := &JSON{parent: a.parent}

	parent := a.parent
	g.update = a.update.append(func() {
		if a.parent != parent {
			return
		}

		if i < 0 {
			if len(a.values) > 0 {
				a.values[i%len(a.values)+len(a.values)] = g
			} else {
				a.values = append(a.values, g)
			}
		} else {
			if len(a.values) < i+1 {
				values := make([]*JSON, i+1)
				copy(values, a.values)
				a.values = values
			}
			a.values[i] = g
		}

		for p := a.parent; p != nil; p = p.parent {
			p.raw = nil
		}
	})

	return g
}

func (a *array) Remove(g *JSON) {
	for i := 0; i < len(a.values); i++ {
		if a.values[i] == g {
			a.values = append(a.values[:i], a.values[i+1:]...)
			break
		}
	}
}
