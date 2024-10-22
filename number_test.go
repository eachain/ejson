package ejson

import "testing"

func TestIsNumber(t *testing.T) {
	if new(JSON).IsNumber() {
		t.Fatalf("empty json should not be a number")
	}

	if !FromString("0").IsNumber() {
		t.Fatal("0 should be a number")
	}
	if !FromString("-0").IsNumber() {
		t.Fatal("-0 should be a number")
	}
	if FromString("a").IsNumber() {
		t.Fatal("a should not be a number")
	}

	if !FromString(`"0"`).IsNumber() {
		t.Fatal(`"0" should be a number`)
	}
	if !FromString(`"-0"`).IsNumber() {
		t.Fatal(`"-0" should be a number`)
	}
	if !FromString(`"0.123"`).IsNumber() {
		t.Fatal(`"0.123" should be a number`)
	}
	if !FromString(`"-0.123"`).IsNumber() {
		t.Fatal(`"-0.123" should be a number`)
	}
}

func TestIsInt(t *testing.T) {
	if !FromString(`123`).IsInt() {
		t.Fatal("123 should be an int")
	}

	if !FromString(`-123`).IsInt() {
		t.Fatal("-123 should be an int")
	}

	if FromString(`-1.23`).IsInt() {
		t.Fatal("-1.23 should not be an int")
	}

	if FromString(`1e3`).IsInt() {
		t.Fatal("1e3 should not be an int")
	}
}

func TestTryInt(t *testing.T) {
	if v, ok := FromString(`123`).TryInt(); !ok {
		t.Fatal("123 should be an int")
	} else if v != 123 {
		t.Fatalf("123 returns %v", v)
	}

	if v, ok := FromString(`-123`).TryInt(); !ok {
		t.Fatal("-123 should be an int")
	} else if v != -123 {
		t.Fatalf("-123 returns %v", v)
	}

	if v, ok := FromString(`1e3`).TryInt(); !ok {
		t.Fatal("1e3 should be an int")
	} else if v != 1e3 {
		t.Fatalf("1e3 returns %v", v)
	}

	if v, ok := FromString(`-1e3`).TryInt(); !ok {
		t.Fatal("-1e3 should be an int")
	} else if v != -1e3 {
		t.Fatalf("-1e3 returns %v", v)
	}
}

func TestIsUint(t *testing.T) {
	if !FromString(`123`).IsUint() {
		t.Fatal("123 should be an uint")
	}

	if FromString(`-123`).IsUint() {
		t.Fatal("-123 should not be an uint")
	}

	if FromString(`1.23`).IsUint() {
		t.Fatal("1.23 should not be an uint")
	}

	if FromString(`1e3`).IsUint() {
		t.Fatal("1e3 should not be an uint")
	}
}

func TestTryUint(t *testing.T) {
	if v, ok := FromString(`123`).TryUint(); !ok {
		t.Fatal("123 should be an uint")
	} else if v != 123 {
		t.Fatalf("123 returns %v", v)
	}

	if v, ok := FromString(`-123`).TryUint(); ok {
		t.Fatalf("-123 should not be an uint: %v", v)
	}

	if v, ok := FromString(`1e3`).TryUint(); !ok {
		t.Fatal("1e3 should be an uint")
	} else if v != 1e3 {
		t.Fatalf("1e3 returns %v", v)
	}

	if v, ok := FromString(`-1e3`).TryUint(); ok {
		t.Fatalf("-1e3 should not be an uint: %v", v)
	}
}

func TestTryFloat(t *testing.T) {
	if v, ok := FromString(`123`).TryFloat(); !ok {
		t.Fatal("123 should be a float")
	} else if v != 123 {
		t.Fatalf("123 returns %v", v)
	}

	if v, ok := FromString(`-123`).TryFloat(); !ok {
		t.Fatal("-123 should be a float")
	} else if v != -123 {
		t.Fatalf("-123 returns %v", v)
	}

	if v, ok := FromString(`1.23`).TryFloat(); !ok {
		t.Fatal("1.23 should be a float")
	} else if v != 1.23 {
		t.Fatalf("-1.23 returns %v", v)
	}
	if v, ok := FromString(`-1.23`).TryFloat(); !ok {
		t.Fatal("-1.23 should be a float")
	} else if v != -1.23 {
		t.Fatalf("-1.23 returns %v", v)
	}

	if v, ok := FromString(`1e3`).TryFloat(); !ok {
		t.Fatal("1e3 should be a float")
	} else if v != 1e3 {
		t.Fatalf("1e3 returns %v", v)
	}
}
