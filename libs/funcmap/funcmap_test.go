package funcmap

import (
	"testing"
)

var (
	m = map[string]interface{}{
		"hello":      func() { print("hello\n") },
		"foobar":     func(a, b, c int) int { return a + b + c },
		"errstring":  "Can not call this as a function",
		"errnumeric": 123456789,
	}
	funcs = New()
)

func TestBind(t *testing.T) {
	for k, v := range m {
		err := funcs.Bind(k, v)
		t.Log("bind ", k, " return ", err)
		if k[:3] == "err" {
			if err == nil {
				t.Error("Bind %s: %s", k, "an error should be paniced.")
			}
		} else {
			if err != nil {
				t.Error("Bind Failed, %s: %s", k, err)
			}
		}
	}
}

func TestCall(t *testing.T) {
	if _, err := funcs.Call("hello"); err != nil {
		t.Error("Call %s: %s", "hello", "an error should be paniced.")
	} else {
		t.Log("call hello return ", err)
	}
	if val, err := funcs.Call("foobar", 0, 1, 2); err != nil {
		t.Error("Call %s: %s", "foobar", err)
	} else {
		t.Logf("call foobar return %d %v", val[0].Int(), err)
	}
	if _, err := funcs.Call("errstring", 0, 1, 2); err == nil {
		t.Error("Call %s: %s", "errstring", "an error should be paniced.")
	} else {
		t.Log("call errstring return ", err)
	}
	if _, err := funcs.Call("errnumeric"); err == nil {
		t.Error("Call %s: %s", "errnumeric", "an error should be paniced.")
	} else {
		t.Log("call errnumeric return ", err)
	}
}
