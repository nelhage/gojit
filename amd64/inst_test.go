package amd64

import (
	"github.com/nelhage/gojit"
	"testing"
)

func TestMov(t *testing.T) {
	buf, e := gojit.Alloc(4096)
	if e != nil {
		t.Fatalf(e.Error())
	}
	defer gojit.Release(buf)

	cases := []struct {
		f     func(*Assembler)
		inout []uintptr
	}{
		{
			func(a *Assembler) {
				a.Mov(Imm(0xdeadbeef), Rax)
				a.Ret()
			},
			[]uintptr{0, 0xdeadbeef},
		},
		{
			func(a *Assembler) {
				a.Mov(Rdi, Rax)
				a.Ret()
			},
			[]uintptr{0, 0, 1, 1, 0xdeadbeef, 0xdeadbeef, 0xffffffffffffffff, 0xffffffffffffffff},
		},
	}

	for i, tc := range cases {
		asm := &Assembler{Buf: buf}

		tc.f(asm)

		f := gojit.Build(buf)

		for j := 0; j < len(tc.inout); j += 2 {
			in := tc.inout[j]
			out := tc.inout[j+1]
			got := f(in)
			if out != got {
				t.Errorf("f[%d](%x) = %x, expect %x",
					i, in, got, out)
			}
		}
	}
}
