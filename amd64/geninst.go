package amd64

import (
	"fmt"
	"reflect"
	"unsafe"
)

func (a *Assembler) Inc(o Operand) {
	o.Rex(a, Register{})
	a.byte(0xff)
	o.ModRM(a, Register{})
}

func (a *Assembler) Dec(o Operand) {
	o.Rex(a, Register{})
	a.byte(0xff)
	o.ModRM(a, Register{1, 0})
}

func (a *Assembler) Incb(o Operand) {
	o.Rex(a, Register{})
	a.byte(0xfe)
	o.ModRM(a, Register{})
}

func (a *Assembler) Decb(o Operand) {
	o.Rex(a, Register{})
	a.byte(0xfe)
	o.ModRM(a, Register{1, 0})
}

func (asm *Assembler) arithmeticImmReg(insn *Instruction, src Imm, dst Register) {
	if insn.imm_r.ok() {
		asm.rex(false, false, false, dst.Val > 7)
		asm.byte(insn.imm_r.value() | (dst.Val & 7))
	} else {
		asm.rex(dst.Bits == 64, false, dst.Val > 7, false)
		asm.byte(insn.imm_rm.op.value())
		asm.modrm(MOD_REG, insn.imm_rm.sub, dst.Val&7)
	}
}

func (asm *Assembler) arithmeticRegReg(insn *Instruction, src Register, dst Register) {
	if insn.r_rm.ok() {
		dst.Rex(asm, src)
		asm.byte(insn.r_rm.value())
		dst.ModRM(asm, src)
	} else {
		src.Rex(asm, dst)
		asm.byte(insn.rm_r.value())
		src.ModRM(asm, dst)
	}
}

func (asm *Assembler) Arithmetic(insn *Instruction, src, dst Operand) {
	switch s := src.(type) {
	case Imm:
		if dr, ok := dst.(Register); ok {
			asm.arithmeticImmReg(insn, s, dr)
		} else {
			dst.Rex(asm, Register{insn.imm_rm.sub, 0})
			asm.byte(insn.imm_rm.op.value())
			dst.ModRM(asm, Register{insn.imm_rm.sub, 0})
		}
		if insn.bits == 8 {
			asm.byte(byte(s.Val))
		} else {
			asm.int32(uint32(s.Val))
		}
		return
	case Register:
		if dr, ok := dst.(Register); ok {
			asm.arithmeticRegReg(insn, s, dr)
		} else {
			dst.Rex(asm, s)
			asm.byte(insn.r_rm.value())
			dst.ModRM(asm, s)
		}
		return
	}
	// if the LHS is neither an immediate nor a register, the rhs
	// must be a register
	dr, ok := dst.(Register)
	if !ok {
		panic(fmt.Sprintf("arithmetic: %#v/%#v not supported!", src, dst))
	}

	src.Rex(asm, dr)
	asm.byte(insn.rm_r.value())
	src.ModRM(asm, dr)
}

func (a *Assembler) Add(src, dst Operand) {
	a.Arithmetic(InstAdd, src, dst)
}

func (a *Assembler) And(src, dst Operand) {
	a.Arithmetic(InstAnd, src, dst)
}

func (a *Assembler) Cmp(src, dst Operand) {
	a.Arithmetic(InstCmp, src, dst)
}

func (a *Assembler) Mov(src, dst Operand) {
	a.Arithmetic(InstMov, src, dst)
}

func (a *Assembler) Movb(src, dst Operand) {
	a.Arithmetic(InstMovb, src, dst)
}

func (a *Assembler) MovAbs(src uint64, dst Register) {
	a.rex(true, false, false, dst.Val > 7)
	a.byte(InstMov.imm_r.value() | (dst.Val & 7))
	a.int64(src)
}

func (a *Assembler) Or(src, dst Operand) {
	a.Arithmetic(InstOr, src, dst)
}

func (a *Assembler) Lea(src, dst Operand) {
	a.Arithmetic(InstLea, src, dst)
}

func (a *Assembler) Sub(src, dst Operand) {
	a.Arithmetic(InstSub, src, dst)
}

func (a *Assembler) Test(src, dst Operand) {
	a.Arithmetic(InstTest, src, dst)
}

func (a *Assembler) Testb(src, dst Operand) {
	a.Arithmetic(InstTestb, src, dst)
}

func (a *Assembler) Xor(src, dst Operand) {
	a.Arithmetic(InstXor, src, dst)
}

func (a *Assembler) Ret() {
	a.byte(0xc3)
}

func (a *Assembler) Call(dst Operand) {
	if _, ok := dst.(Imm); ok {
		panic("can't call(Imm); use CallRel instead.")
	} else {
		a.byte(0xff)
		dst.ModRM(a, Register{0x2, 64})
	}
}

func (a *Assembler) CallRel(dst uintptr) {
	a.byte(0xe8)
	a.rel32(dst)
}

// Clobbers RDX
func (a *Assembler) CallFunc(f interface{}) {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		panic("CallFunc: Can't call non-func")
	}
	ival := *(*struct {
		typ uintptr
		fun uintptr
	})(unsafe.Pointer(&f))

	a.MovAbs(uint64(ival.fun), Rdx)
	a.Call(Indirect{Rdx, 0, 64})
}

func (a *Assembler) Push(src Operand) {
	if imm, ok := src.(Imm); ok {
		a.byte(0x68)
		a.int32(uint32(imm.Val))
	} else {
		a.byte(0xff)
		src.ModRM(a, Register{0x6, 64})
	}
}

func (a *Assembler) Pop(dst Operand) {
	switch d := dst.(type) {
	case Imm:
		panic("can't pop imm")
	case Register:
		a.rex(false, false, false, d.Val > 7)
		a.byte(0x58 | (d.Val & 7))
	default:
		dst.Rex(a, Register{0x0, 64})
		a.byte(0x8f)
		dst.ModRM(a, Register{0x0, 64})
	}
}

func (a *Assembler) JmpRel(dst uintptr) {
	a.byte(0xe9)
	a.rel32(dst)
}

func (a *Assembler) JccShort(cc byte, off int8) {
	a.byte(0x70 | cc)
	a.byte(byte(off))
}

func (a *Assembler) JccRel(cc byte, dst uintptr) {
	a.byte(0x0f)
	a.byte(0x80 | cc)
	a.rel32(dst)
}
