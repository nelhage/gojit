package gojit

import (
	"testing"
)

func TestRET(t *testing.T) {
	b, e := Alloc(4096)
	if e != nil {
		t.Fatalf("Alloc: %s", e.Error())
	}
	defer Release(b)
	b[0] = 0xc3
	Call(b)
}

func TestBuild(t *testing.T) {
	b, e := Alloc(4096)
	if e != nil {
		t.Fatalf("Alloc: %s", e.Error())
	}
	defer Release(b)
	b[0] = 0xc3
	f := Build(b)
	f(0)
}
