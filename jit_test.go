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

func TestInOut(t *testing.T) {
	b, e := Alloc(4096)
	if e != nil {
		t.Fatalf("Alloc: %s", e.Error())
	}
	defer Release(b)
	/*
	* 0000000000000000 <inc>:
	*    0:	48 89 f8             	mov    %rdi,%rax
	*    3:	48 ff c0             	inc    %rax
	*    6:	c3                   	retq
	 */
	copy(b, []byte{0x48, 0x89, 0xf8, 0x48, 0xff, 0xc0, 0xc3})

	f := Build(b)
	out := f(128)
	if out != 129 {
		t.Errorf("expected f(128) = 128, got %d", out)
	}
}
