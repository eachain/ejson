package ejson

import "testing"

func TestIsObject(t *testing.T) {
	if !FromString(`{}`).IsObject() {
		t.Fatal("{} should be an object")
	}
}

func TestKeys(t *testing.T) {
	results := []string{"a", "b", "c"}
	keys := FromString(`{"a":1, "b":2, "c":3}`).Keys()
	if len(results) != len(keys) {
		t.Fatalf("keys returns %v, should be %v", len(keys), len(results))
	}
	for i := range results {
		if keys[i] != results[i] {
			t.Fatalf("keys[%v] returns %q, should be %q", i, keys[i], results[i])
		}
	}
}

func TestAny(t *testing.T) {
	g := FromString(`{"err_code":1, "errcode":"2"}`)
	if c := g.Any("err_code", "errcode").Int(); c != 1 {
		t.Fatalf("err_code: %v", c)
	}
	if c := g.Any("errcode", "err_code").Int(); c != 2 {
		t.Fatalf("errcode: %v", c)
	}
}
