package ejson

import "testing"

func TestIsArray(t *testing.T) {
	if !FromString(`[]`).IsArray() {
		t.Fatalf("[] should be an array")
	}

	if !FromString(`[`).IsArray() {
		t.Fatalf("[ should be an empty array")
	}
}
