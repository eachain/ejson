package ejson

import (
	"encoding/json"
	"testing"
)

func TestValid(t *testing.T) {
	if !Valid([]byte(`{}`)) {
		t.Fatal("{} should be a valid json")
	}
}

func TestLen(t *testing.T) {
	if FromString(`null`).Len() != 0 {
		t.Fatal("null len should be 0")
	}

	if FromString(`"abc"`).Len() != 3 {
		t.Fatal("str len should be 3")
	}

	if FromString(`{"a":1, "b":2}`).Len() != 2 {
		t.Fatal("object len should be 2")
	}

	if FromString(`["a", 2, null, {}]`).Len() != 4 {
		t.Fatal("array len should be 4")
	}

	if FromString(`true`).Len() != 0 {
		t.Fatal("bool len should be 0")
	}
}

func TestArrayIndex(t *testing.T) {
	if FromString(`[123, 456]`).ArrayIndex(0).Int() != 123 {
		t.Fatal("array first value should be 123")
	}

	if FromString(`[123, 456]`).ArrayIndex(1).Int() != 456 {
		t.Fatal("array second value should be 456")
	}

	if FromString(`[123, 456]`).ArrayIndex(-1).Int() != 456 {
		t.Fatal("array last value should be 456")
	}

	if FromString(`[123, 456]`).ArrayIndex(3).Int() != 0 {
		t.Fatal("array out of bound value should be 0")
	}

	{
		g := new(JSON)
		g.ArrayIndex(0).Set(123)
		if p, _ := g.MarshalJSON(); string(p) != "[123]" {
			t.Fatalf("set: %s", p)
		}
	}

	{
		g := new(JSON)
		g.ArrayIndex(0).ArrayIndex(-1).ArrayIndex(2).Set(123)
		if p, _ := g.MarshalJSON(); string(p) != "[[[null,null,123]]]" {
			t.Fatalf("set: %s", p)
		}
	}
}

func TestObjectIndex(t *testing.T) {
	if FromString(`{"a": 123, "b": 456}`).ObjectIndex("a").Int() != 123 {
		t.Fatal("object.a should be 123")
	}
	if FromString(`{"a": 123, "b": 456}`).ObjectIndex("b").Int() != 456 {
		t.Fatal("object.b should be 456")
	}
	if FromString(`{"a": 123, "b": 456}`).ObjectIndex("c").Int() != 0 {
		t.Fatal("object.c (not exists) should be 0")
	}

	g := new(JSON)
	g.ObjectIndex("a").ArrayIndex(0).ObjectIndex("b").ArrayIndex(0).ObjectIndex("c").ArrayIndex(0).Set(123)
	if p, _ := g.MarshalJSON(); string(p) != `{"a":[{"b":[{"c":[123]}]}]}` {
		t.Fatalf("set: %s", p)
	}
}

func TestStrJSON(t *testing.T) {
	if a := FromString(`"[{\"a\":\"123\"}]"`).StrJSON().ArrayIndex(0).ObjectIndex("a").StrJSON().Int(); a != 123 {
		t.Fatalf("object.a should be 123 but %v", a)
	}

	g := new(JSON)
	g.StrJSON().ObjectIndex("a").ArrayIndex(0).ObjectIndex("b").ArrayIndex(0).ObjectIndex("c").ArrayIndex(0).StrJSON().Set(123)
	if p, _ := g.MarshalJSON(); string(p) != `"{\"a\":[{\"b\":[{\"c\":[\"123\"]}]}]}"` {
		t.Fatalf("set: %s", p)
	}

	g = FromString(`"[{\"a\":\"123\"}]"`)
	g.Get("[0].a").Set(456)
	if p, _ := g.MarshalJSON(); string(p) != `"[{\"a\":456}]"` {
		t.Fatalf("set: %s", p)
	}
	g.Get("[0].a").StrJSON().Set(789)
	if p, _ := g.MarshalJSON(); string(p) != `"[{\"a\":\"789\"}]"` {
		t.Fatalf("set: %s", p)
	}

	g = FromString(`"[{\"a\":\"123\", \"b\":\"456\"}]"`)
	g.Get("[0].b").Remove()
	if p, _ := g.MarshalJSON(); string(p) != `"[{\"a\":\"123\"}]"` {
		t.Fatalf("set: %s", p)
	}

	g = FromString(`"[{\"a\":\"123\", \"b\":\"456\"}]"`)
	g.Get("[0].b").StrJSON().Remove()
	if p, _ := g.MarshalJSON(); string(p) != `"[{\"a\":\"123\",\"b\":\"null\"}]"` {
		t.Fatalf("set: %s", p)
	}

	g = FromString(`"[{\"a\":\"123\", \"b\":\"456\"}]"`)
	g.Get("[0].b.x").Set(789)
	if p, _ := g.MarshalJSON(); string(p) != `"[{\"a\":\"123\",\"b\":\"{\\\"x\\\":789}\"}]"` {
		t.Fatalf("set: %s", p)
	}
	g.Get("[0].b").Set(true)
	if p, _ := g.MarshalJSON(); string(p) != `"[{\"a\":\"123\",\"b\":true}]"` {
		t.Fatalf("set: %s", p)
	}
}

func TestRemove(t *testing.T) {
	{
		g := FromString(`123`)
		g.Remove()
		if g.UnsafeString() != "123" {
			t.Fatalf("remove self result: %s", g.UnsafeString())
		}
	}

	{
		g := FromString(`{"a":123}`)
		g.ObjectIndex("a").Remove()
		if g.UnsafeString() != "{}" {
			t.Fatalf("object remove result: %s", g.UnsafeString())
		}
	}

	{
		g := FromString(`[123]`)
		g.ArrayIndex(0).Remove()
		if g.UnsafeString() != "[]" {
			t.Fatalf("array remove result: %s", g.UnsafeString())
		}
	}

	{
		g := FromString(`[{"a":[{"x":true}]}]`)
		if !g.Get("[0].a[0].x").Bool() {
			t.Fatalf("[0].a[0].x result should be true")
		}
		g.ArrayIndex(0).Remove()
		if g.UnsafeString() != "[]" {
			t.Fatalf("array remove result: %s", g.UnsafeString())
		}
	}
}

func TestString(t *testing.T) {
	if s := FromString(`123`).String(); s != "123" {
		t.Fatalf("string: %q", s)
	}
	if s := FromString(`"123"`).String(); s != "123" {
		t.Fatalf("string: %q", s)
	}
	if s := FromString(`"123"`).UnsafeString(); s != `"123"` {
		t.Fatalf("unsafe string: %q", s)
	}
}

func TestExists(t *testing.T) {
	if new(JSON).Exists() {
		t.Fatal("empty json should not exists")
	}

	if FromString(`{}`).Get("a.b.c").Exists() {
		t.Fatal("a.b.c should not exists")
	}

	if !FromString(`[{"a":[{"x": 123}]}]`).Get("0.a.0.x").Exists() {
		t.Fatal("0.a.0.x should exists")
	}
}

func TestValue(t *testing.T) {
	var a int
	FromString(`123`).Value(&a)
	if a != 123 {
		t.Fatalf("value: %v", a)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	g := new(JSON)
	json.Unmarshal([]byte(`{"a":123}`), g)
	if a := g.ObjectIndex("a").Int(); a != 123 {
		t.Fatalf("a value: %v", a)
	}
}

func TestSet(t *testing.T) {
	g := new(JSON)
	g.Set("123")
	if a := g.Int(); a != 123 {
		t.Fatalf("a value: %v", a)
	}
	if g.UnsafeString() != `"123"` {
		t.Fatalf("string: %q", g.UnsafeString())
	}
}

func TestGenJSON(t *testing.T) {
	resp := new(JSON)
	resp.Get("code").Set(0)
	resp.Get("msg").Set("OK")
	data := resp.Get("data")
	data.Get("object.a").Set(123)
	data.Get("array[0]").Set(456)

	result := `{"code":0,"msg":"OK","data":{"object":{"a":123},"array":[456]}}`
	if resp.UnsafeString() != result {
		t.Fatalf("string: %s", resp.UnsafeString())
	}
}

func TestSetChild(t *testing.T) {
	origin := `{"code":0,"msg":"OK","data":[{"object":{"a":123},"array":[456]}]}`
	result := `{"code":0,"msg":"OK","data":{"object":{"a":123},"array":[456]}}`

	{
		resp := FromString(origin)
		resp.Get("data").Set(resp.Get("data[0]"))
		if resp.UnsafeString() != result {
			t.Fatalf("string: %s", resp.UnsafeString())
		}
	}

	{
		resp := FromString(origin)
		elem := resp.Get("data[0]")
		elem.Remove()
		resp.Get("data").Set(elem)
		if resp.UnsafeString() != result {
			t.Fatalf("string: %s", resp.UnsafeString())
		}
	}

	{
		resp := FromString(origin)
		resp.Set(resp.Get("data[0].object"))
		if a := resp.Get("a").Int(); a != 123 {
			t.Fatalf("a value: %v", a)
		}
	}
}

func TestItoa(t *testing.T) {
	origin := `{"code":0,"msg":"OK","data":[{"object":{"a":123},"array":[456]}]}`
	resp := FromString(origin)
	resp.Get("data").Set(resp.Get("data[0]"))
	a := resp.Get("data.object.a")
	a.Set(a.String())
	e := resp.Get("data.array[0]")
	e.Set(e.UnsafeString())

	result := `{"code":0,"msg":"OK","data":{"object":{"a":"123"},"array":["456"]}}`
	if resp.UnsafeString() != result {
		t.Fatalf("string: %s", resp.UnsafeString())
	}
}
