package ejson

import "testing"

func TestIsNull(t *testing.T) {
	if !FromString(`null`).IsNull() {
		t.Fatal("null should be a null value")
	}
}
