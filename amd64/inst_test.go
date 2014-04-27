package amd64

import (
	"fmt"
	"github.com/nelhage/gojit"
	"testing"
)

var mem []byte = make([]byte, 64)

func TestMov(t *testing.T) {
	cases := []struct {
		f     func(*Assembler)
		inout []uintptr
	}{
		{
			func(a *Assembler) {
				a.Mov(Imm{U32(0xdeadbeef)}, Rax)
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
		{
			func(a *Assembler) {
				a.Mov(Imm{U32(0xcafebabe)}, Indirect{Rdi, 0, 64})
				a.Mov(Indirect{Rdi, 0, 64}, Rax)
				a.Ret()
			},
			[]uintptr{gojit.Addr(mem), 0xffffffffcafebabe},
		},
		{
			func(a *Assembler) {
				a.Mov(Imm{U32(0xf00dface)}, R10)
				a.Mov(R10, Rax)
				a.Ret()
			},
			[]uintptr{0, 0xf00dface},
		},
	}

	buf, e := gojit.Alloc(4096)
	if e != nil {
		t.Fatalf(e.Error())
	}
	defer gojit.Release(buf)

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

func TestArith(t *testing.T) {
	cases := []struct {
		insn     *Instruction
		lhs, rhs int32
		out      uintptr
	}{
		{InstAdd, 20, 30, 50},
		{InstAdd, 0x7fffffff, 0x70000001, 0xf0000000},
		{InstAnd, 0x77777777, U32(0xffffffff), 0x77777777},
		{InstAnd, 0x77777777, U32(0x88888888), 0},
		{InstOr, 0x77777777, U32(0x88888888), 0xffffffff},
		{InstOr, 1, 0, 1},
		{InstSub, 5, 10, 5},
		{InstSub, 10, 5, 0xfffffffffffffffb},
	}

	buf, e := gojit.Alloc(4096)
	if e != nil {
		t.Fatalf(e.Error())
	}
	defer gojit.Release(buf)

	for _, tc := range cases {
		asm := &Assembler{buf, 0}
		var funcs []func(uintptr) uintptr
		if tc.insn.imm_r != 0 {
			funcs = append(funcs, gojit.Build(asm.Buf[asm.Off:]))
			asm.Mov(Imm{tc.rhs}, Rax)
			asm.Arithmetic(tc.insn, Imm{tc.lhs}, Rax)
			asm.Ret()
		}
		if tc.insn.imm_rm.op != 0 {
			funcs = append(funcs, gojit.Build(asm.Buf[asm.Off:]))
			asm.Mov(Imm{0}, Indirect{Rdi, 0, 0})
			asm.Mov(Imm{tc.rhs}, Indirect{Rdi, 0, 32})
			asm.Arithmetic(tc.insn, Imm{tc.lhs}, Indirect{Rdi, 0, 64})
			asm.Mov(Indirect{Rdi, 0, 64}, Rax)
			asm.Ret()
		}
		if tc.insn.r_rm != 0 {
			funcs = append(funcs, gojit.Build(asm.Buf[asm.Off:]))
			asm.Mov(Imm{tc.lhs}, R10)
			asm.Mov(Imm{0}, Indirect{Rdi, 0, 0})
			asm.Mov(Imm{tc.rhs}, Indirect{Rdi, 0, 32})
			asm.Arithmetic(tc.insn, R10, Indirect{Rdi, 0, 64})
			asm.Mov(Indirect{Rdi, 0, 64}, Rax)
			asm.Ret()
		}
		if tc.insn.rm_r != 0 {
			funcs = append(funcs, gojit.Build(asm.Buf[asm.Off:]))
			asm.Mov(Imm{0}, Indirect{Rdi, 0, 0})
			asm.Mov(Imm{tc.lhs}, Indirect{Rdi, 0, 32})
			asm.Mov(Imm{tc.rhs}, R10)
			asm.Arithmetic(tc.insn, Indirect{Rdi, 0, 64}, R10)
			asm.Mov(R10, Rax)
			asm.Ret()
		}

		for i, f := range funcs {
			got := f(gojit.Addr(mem))
			if got != tc.out {
				t.Errorf("%s(0x%x,0x%x) [%d] = 0x%x (expect 0x%x)",
					tc.insn.Mnemonic, tc.lhs, tc.rhs, i, got, tc.out)
			} else if testing.Verbose() {
				// We don't use `testing.Logf` because
				// if we panic inside JIT'd code, the
				// runtime dies horrible (rightfully
				// so!), and so the `testing` cleanup
				// code never runs, and we never see
				// log messages. We want to get these
				// out as soon as possible, so we
				// write them directly.
				fmt.Printf("OK %d %s(0x%x,0x%x) = 0x%x\n",
					i, tc.insn.Mnemonic, tc.lhs, tc.rhs, got)
			}
		}
	}
}
