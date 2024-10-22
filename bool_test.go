package ejson

import "testing"

func TestIsBool(t *testing.T) {
	if !FromString(`true`).IsBool() {
		t.Fatal("true should be bool")
	}

	if !FromString(`false`).IsBool() {
		t.Fatal("false should be bool")
	}
}

func TestBool(t *testing.T) {
	if !FromString(`true`).Bool() {
		t.Fatal("true returns false")
	}

	if FromString(`false`).Bool() {
		t.Fatal("false returns true")
	}

	if FromString(`null`).Bool() {
		t.Fatal("null returns true")
	}

	if FromString(`""`).Bool() {
		t.Fatal("empty str returns true")
	}

	if !FromString(`"abc"`).Bool() {
		t.Fatal("not empty str returns false")
	}

	if FromString(`0`).Bool() {
		t.Fatal("0 returns true")
	}

	if !FromString(`1`).Bool() {
		t.Fatal("1 returns false")
	}

	if FromString(`0.0`).Bool() {
		t.Fatal("0.0 returns true")
	}

	if !FromString(`0.1`).Bool() {
		t.Fatal("0.1 returns false")
	}

	if FromString(`-0`).Bool() {
		t.Fatal("-0 returns true")
	}

	if !FromString(`-1`).Bool() {
		t.Fatal("-1 returns false")
	}

	if FromString(`-0.0`).Bool() {
		t.Fatal("-0.0 returns true")
	}

	if !FromString(`-0.1`).Bool() {
		t.Fatal("0.1 returns false")
	}

	if FromString(`{}`).Bool() {
		t.Fatal("empty object returns true")
	}

	if FromString(`[]`).Bool() {
		t.Fatal("empty array returns true")
	}

	if !FromString(`{"a":1}`).Bool() {
		t.Fatal("not empty object returns false")
	}

	if !FromString(`[1]`).Bool() {
		t.Fatal("not empty array returns false")
	}

	if new(JSON).Bool() {
		t.Fatal("empty json returns true")
	}
}
