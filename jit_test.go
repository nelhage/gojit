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
	for _, build := range []func([]byte) func(){Build, BuildCgo} {
		f := build(b)
		f()
	}
}

func TestBuildTo(t *testing.T) {
	b, e := Alloc(4096)
	if e != nil {
		t.Fatalf("Alloc: %s", e.Error())
	}
	defer Release(b)
	// 0000000000000000 <inc>:
	//    0:	48 8b 07             	mov    (%rdi),%rax
	//    3:	48 ff c0             	inc    %rax
	//    6:	48 89 47 08          	mov    %rax,0x8(%rdi)
	//    a:	c3                   	retq
	copy(b, []byte{
		0x48, 0x8b, 0x07,
		0x48, 0xff, 0xc0,
		0x48, 0x89, 0x47, 0x08,
		0xc3,
	})

	for _, buildTo := range []func([]byte, interface{}){BuildTo, BuildToCgo} {
		var f1 func(uintptr) uintptr
		buildTo(b, &f1)

		got := f1(128)
		if got != 129 {
			t.Errorf("expected f(128) = 129, got %d", got)
		}
	}
}

func BenchmarkEmptyCall(b *testing.B) {
	benchmarkEmptyCall(b, Build)
}

func BenchmarkEmptyCgoCall(b *testing.B) {
	benchmarkEmptyCall(b, BuildCgo)
}

func benchmarkEmptyCall(b *testing.B, build func([]byte) func()) {
	buf, e := Alloc(PageSize)
	if e != nil {
		b.Fatalf("alloc: %s", e.Error())
	}
	defer Release(buf)

	buf[0] = 0xc3

	f := build(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f()
	}
}
