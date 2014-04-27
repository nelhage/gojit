package gojit

import (
	"testing"
)

func TestBuild(t *testing.T) {
	b, e := Alloc(4096)
	if e != nil {
		t.Fatalf("Alloc: %s", e.Error())
	}
	defer Release(b)
	b[0] = 0xc3
	f := Build(b)
	f()
}

func TestBuildTo(t *testing.T) {
	b, e := Alloc(4096)
	if e != nil {
		t.Fatalf("Alloc: %s", e.Error())
	}
	defer Release(b)
	// 0:	48 8b 44 24 08       	mov    0x8(%rsp),%rax
	// 5:	48 ff 00             	incq   (%rax)
	// 8:	c3                   	retq
	copy(b, []byte{0x48, 0x8b, 0x44, 0x24, 0x08, 0x48, 0xff, 0x00, 0xc3})

	var f1 func(*uint64)
	BuildTo(b, &f1)

	x := uint64(128)
	f1(&x)
	if x != 129 {
		t.Errorf("expected f(128) = 129, got %d", x)
	}
}
