package amd64

import (
	"github.com/nelhage/gojit"
	"testing"
)

//   48 89 fe             	mov    %rdi,%rsi
//   48 8b 3f             	mov    (%rdi),%rdi
var Preamble = []byte{0x48, 0x89, 0xfe, 0x48, 0x8b, 0x3f}

//   48 89 46 08          	mov    %rax,0x8(%rsi)
var Post = []byte{0x48, 0x89, 0x46, 0x08}

func begin(a *Assembler) {
	copy(a.Buf[a.Off:], Preamble)
	a.Off += len(Preamble)
}

func finish(a *Assembler) func(uintptr) uintptr {
	copy(a.Buf[a.Off:], Post)
	a.Off += len(Post)
	a.Ret()
	var f1 func(uintptr) uintptr
	gojit.BuildTo(a.Buf, &f1)
	a.Buf = a.Buf[a.Off:]
	a.Off = 0
	return f1
}

func newAsm(t testing.TB) *Assembler {
	buf, e := gojit.Alloc(4096)
	if e != nil {
		t.Fatalf("alloc: ", e.Error())
	}
	return &Assembler{buf, 0}
}
