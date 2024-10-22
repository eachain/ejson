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
