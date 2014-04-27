package amd64

import "fmt"

type Operand interface {
	// isOperand is unexported prevents external packages from
	// implementing Operand.
	isOperand()

	Rex(asm *Assembler, reg Register)
	ModRM(asm *Assembler, reg Register)
}

type Imm struct {
	Val int32
}

func U32(u uint32) int32 {
	return int32(u)
}

func (i Imm) isOperand() {}
func (i Imm) Rex(asm *Assembler, reg Register) {
	panic("Imm.Rex")
}
func (i Imm) ModRM(asm *Assembler, reg Register) {
	panic("Imm.ModRM")
}

type Register struct {
	Val  byte
	Bits byte
}

func (r Register) isOperand() {}
func (i Register) Rex(asm *Assembler, reg Register) {
	asm.rexBits(i.Bits, reg.Bits, reg.Val > 7, false, i.Val > 7)
}

func (r Register) ModRM(asm *Assembler, reg Register) {
	if reg.Bits != r.Bits {
		panic(fmt.Sprintf("mismatched Bits %d!=%d", r.Bits, reg.Bits))
	}
	asm.modrm(MOD_REG, reg.Val&7, r.Val&7)
}

var (
	Eax = Register{0, 32}
	Rax = Register{0, 64}
	Ecx = Register{1, 32}
	Rcx = Register{1, 64}
	Edx = Register{2, 32}
	Rdx = Register{2, 64}
	Ebx = Register{3, 32}
	Rbx = Register{3, 64}
	Esp = Register{4, 32}
	Rsp = Register{4, 64}
	Ebp = Register{5, 32}
	Rbp = Register{5, 64}
	Esi = Register{6, 32}
	Rsi = Register{6, 64}
	Edi = Register{7, 32}
	Rdi = Register{7, 64}

	R8d  = Register{8, 32}
	R8   = Register{8, 64}
	R9d  = Register{9, 32}
	R9   = Register{9, 64}
	R10d = Register{10, 32}
	R10  = Register{10, 64}
	R11d = Register{11, 32}
	R11  = Register{11, 64}
	R12d = Register{12, 32}
	R12  = Register{12, 64}
	R13d = Register{13, 32}
	R13  = Register{13, 64}
	R14d = Register{14, 32}
	R14  = Register{14, 64}
	R15d = Register{15, 32}
	R15  = Register{15, 64}
)

const (
	REG_DISP32 = 5
	REG_SIB    = 4
)

type Indirect struct {
	Base   Register
	Offset int32
	Bits   byte
}

func (i Indirect) short() bool {
	return int32(int8(i.Offset)) == i.Offset
}

func (i Indirect) isOperand() {}
func (i Indirect) Rex(asm *Assembler, reg Register) {
	asm.rexBits(reg.Bits, i.Bits, reg.Val > 7, false, i.Base.Val > 7)
}

func (i Indirect) ModRM(asm *Assembler, reg Register) {
	if i.Offset == 0 {
		asm.modrm(MOD_INDIR, reg.Val&7, i.Base.Val&7)
	} else if i.short() {
		asm.modrm(MOD_INDIR_DISP8, reg.Val&7, i.Base.Val&7)
		asm.byte(byte(i.Offset))
	} else {
		asm.modrm(MOD_INDIR_DISP32, reg.Val&7, i.Base.Val&7)
		asm.int32(uint32(i.Offset))
	}
}

type Absolute uint64

func (i Absolute) isOperand() {}
func (i Absolute) Rex(asm *Assembler, reg Register) {
	asm.rex(reg.Bits == 64, reg.Val > 7, false, false)
}
func (i Absolute) ModRM(asm *Assembler, reg Register) {
	asm.modrm(MOD_INDIR, reg.Val&7, REG_DISP32)
	asm.int64(uint64(i))
}

type Scale struct {
	scale byte
}

var (
	Scale1 = Scale{SCALE_1}
	Scale2 = Scale{SCALE_2}
	Scale4 = Scale{SCALE_4}
	Scale8 = Scale{SCALE_8}
)

type SIB struct {
	Offset      uint32
	Base, Index Register
	Scale       Scale
}

func (s SIB) isOperand() {}
func (s SIB) Rex(asm *Assembler, reg Register) {
	asm.rex(reg.Bits == 64, reg.Val > 7, s.Index.Val > 7, s.Base.Val > 7)
}

func (s SIB) ModRM(asm *Assembler, reg Register) {
	if s.Offset != 0 {
		asm.modrm(MOD_INDIR_DISP32, reg.Val&7, REG_SIB)
		asm.sib(s.Scale.scale, s.Index.Val&7, s.Base.Val&7)
		asm.int32(s.Offset)
	} else {
		asm.modrm(MOD_INDIR, reg.Val&7, REG_SIB)
		asm.sib(s.Scale.scale, s.Index.Val&7, s.Base.Val&7)
	}
}
