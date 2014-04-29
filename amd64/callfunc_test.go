package amd64

import (
	"github.com/nelhage/gojit"
	"testing"
)

func TestCallFunc(t *testing.T) {
	asm := newAsm(t)
	defer gojit.Release(asm.Buf)

	called := false

	asm.CallFunc(func() { called = true })
	asm.Ret()

	gojit.Build(asm.Buf)()

	if !called {
		t.Error("CallFunc did not call the function")
	}
}

func BenchmarkGoCall(b *testing.B) {
	asm := newAsm(b)
	defer gojit.Release(asm.Buf)

	f := func() {}
	asm.CallFunc(f)
	asm.Ret()

	jit := gojit.Build(asm.Buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jit()
	}

}
