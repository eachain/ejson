package ejson

import "testing"

func TestIsStr(t *testing.T) {
	if !FromString(`""`).IsStr() {
		t.Fatalf(`"" should be a string`)
	}
}

func TestStr(t *testing.T) {
	if s := FromString(`"abc"`).Str(); s != "abc" {
		t.Fatalf("str returns: %q", s)
	}
}

func TestEscapeHTML(t *testing.T) {
	html := "<&>"
	g := new(JSON)
	g.Set(html)
	if p, _ := g.MarshalJSON(); string(p) != `"<&>"` {
		t.Fatalf("escape html: %s", p)
	}

	g = new(JSON)
	g.ObjectIndex("a").Set(html)
	if p, _ := g.MarshalJSON(); string(p) != `{"a":"<&>"}` {
		t.Fatalf("escape html: %s", p)
	}

	g = new(JSON)
	g.ArrayIndex(0).Set(html)
	if p, _ := g.MarshalJSON(); string(p) != `["<&>"]` {
		t.Fatalf("escape html: %s", p)
	}
}
