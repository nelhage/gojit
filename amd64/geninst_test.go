package amd64

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/nelhage/gojit"
)

var mem []byte = make([]byte, 64)

type simple struct {
	f     func(*Assembler)
	inout []uintptr
}

func TestMov(t *testing.T) {
	cases := []simple{
		{
			func(a *Assembler) {
				a.Mov(Imm{U32(0xdeadbeef)}, Rax)
			},
			[]uintptr{0, 0xdeadbeef},
		},
		{
			func(a *Assembler) {
				a.Mov(Rdi, Rax)
			},
			[]uintptr{0, 0, 1, 1, 0xdeadbeef, 0xdeadbeef, 0xffffffffffffffff, 0xffffffffffffffff},
		},
		{
			func(a *Assembler) {
				a.Mov(Imm{U32(0xcafebabe)}, Indirect{Rdi, 0, 64})
				a.Mov(Indirect{Rdi, 0, 64}, Rax)
			},
			[]uintptr{gojit.Addr(mem), 0xffffffffcafebabe},
		},
		{
			func(a *Assembler) {
				a.Mov(Imm{U32(0xf00dface)}, R10)
				a.Mov(R10, Rax)
			},
			[]uintptr{0, 0xf00dface},
		},
	}

	testSimple("mov", t, cases)
}

func TestIncDec(t *testing.T) {
	cases := []simple{
		{
			func(a *Assembler) {
				a.Mov(Rdi, Rax)
				a.Inc(Rax)
			},
			[]uintptr{0, 1, 10, 11},
		},
		{
			func(a *Assembler) {
				a.Mov(Rdi, Rax)
				a.Dec(Rax)
			},
			[]uintptr{1, 0, 11, 10},
		},
		{
			func(a *Assembler) {
				a.Mov(Imm{0x11223344}, Indirect{Rdi, 0, 32})
				a.Incb(Indirect{Rdi, 1, 8})
				a.Mov(Indirect{Rdi, 0, 32}, Eax)
			},
			[]uintptr{gojit.Addr(mem), 0x11223444},
		},
		{
			func(a *Assembler) {
				a.Mov(Imm{0x11223344}, Indirect{Rdi, 0, 32})
				a.Decb(Indirect{Rdi, 1, 8})
				a.Mov(Indirect{Rdi, 0, 32}, Eax)
			},
			[]uintptr{gojit.Addr(mem), 0x11223244},
		},
	}
	testSimple("inc/dec", t, cases)
}

func testSimple(name string, t *testing.T, cases []simple) {
	buf, e := gojit.Alloc(gojit.PageSize)
	if e != nil {
		t.Fatalf(e.Error())
	}
	defer gojit.Release(buf)

	for i, tc := range cases {
		asm := &Assembler{Buf: buf}
		begin(asm)
		tc.f(asm)
		f := finish(asm)

		runtime.GC()

		for j := 0; j < len(tc.inout); j += 2 {
			in := tc.inout[j]
			out := tc.inout[j+1]
			got := f(in)
			if out != got {
				t.Errorf("f(%s)[%d](%x) = %x, expect %x",
					name, i, in, got, out)
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

	buf, e := gojit.Alloc(gojit.PageSize)
	if e != nil {
		t.Fatalf(e.Error())
	}
	defer gojit.Release(buf)

	for _, tc := range cases {
		asm := &Assembler{buf, 0, CgoABI}
		var funcs []func(uintptr) uintptr
		if tc.insn.imm_r.ok() {
			begin(asm)
			asm.Mov(Imm{tc.rhs}, Rax)
			asm.Arithmetic(tc.insn, Imm{tc.lhs}, Rax)
			funcs = append(funcs, finish(asm))
		}
		if tc.insn.imm_rm.op.ok() {
			begin(asm)
			asm.Mov(Imm{0}, Indirect{Rdi, 0, 0})
			asm.Mov(Imm{tc.rhs}, Indirect{Rdi, 0, 32})
			asm.Arithmetic(tc.insn, Imm{tc.lhs}, Indirect{Rdi, 0, 64})
			asm.Mov(Indirect{Rdi, 0, 64}, Rax)
			funcs = append(funcs, finish(asm))
		}
		if tc.insn.r_rm.ok() {
			begin(asm)
			asm.Mov(Imm{tc.lhs}, R10)
			asm.Mov(Imm{0}, Indirect{Rdi, 0, 0})
			asm.Mov(Imm{tc.rhs}, Indirect{Rdi, 0, 32})
			asm.Arithmetic(tc.insn, R10, Indirect{Rdi, 0, 64})
			asm.Mov(Indirect{Rdi, 0, 64}, Rax)
			funcs = append(funcs, finish(asm))
		}
		if tc.insn.rm_r.ok() {
			begin(asm)
			asm.Mov(Imm{0}, Indirect{Rdi, 0, 0})
			asm.Mov(Imm{tc.lhs}, Indirect{Rdi, 0, 32})
			asm.Mov(Imm{tc.rhs}, R10)
			asm.Arithmetic(tc.insn, Indirect{Rdi, 0, 64}, R10)
			asm.Mov(R10, Rax)
			funcs = append(funcs, finish(asm))
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

func TestMovEsp(t *testing.T) {
	asm := newAsm(t)
	defer gojit.Release(asm.Buf)

	begin(asm)

	// 48 c7 44 24 f8 69 7a 	movq   $0x7a69,-0x8(%rsp)
	// 00 00
	mov_rsp := []byte{0x48, 0xc7, 0x44, 0x24, 0xf8, 0x69, 0x7a, 0x00, 0x00}

	copy(asm.Buf[asm.Off:], mov_rsp)
	asm.Off += len(mov_rsp)
	asm.Mov(Indirect{Rsp, -8, 64}, Rax)
	f := finish(asm)

	got := f(0)
	if got != 31337 {
		t.Errorf("Fatal: mov from esp: got %d != %d", got, 31337)
	}
}
