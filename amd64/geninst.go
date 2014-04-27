package amd64

func (a *Assembler) Inc(o Operand) {
	o.Rex(a, Register{})
	a.byte(0xff)
	o.ModRM(a, Register{})
}

func (asm *Assembler) arithmeticImmReg(insn *Instruction, src Imm, dst Register) {
	if insn.imm_r != 0 {
		// v-- the below would generate MOVABS, which is not what we want
		// asm.rex(dst.Bits == 64, false, false, dst.Val > 7)
		asm.byte(insn.imm_r + dst.Val)
	} else {
		asm.rex(dst.Bits == 64, false, dst.Val > 7, false)
		asm.byte(insn.imm_rm.op)
		asm.modrm(MOD_REG, insn.imm_rm.sub, dst.Val&7)
	}
}

func (asm *Assembler) arithmeticRegReg(insn *Instruction, src Register, dst Register) {
	if insn.r_rm != 0 {
		dst.Rex(asm, src)
		asm.byte(insn.r_rm)
		dst.ModRM(asm, src)
	} else {
		src.Rex(asm, dst)
		asm.byte(insn.rm_r)
		src.ModRM(asm, dst)
	}
}

func (asm *Assembler) arithmetic(insn *Instruction, src, dst Operand) {
	switch s := src.(type) {
	case Imm:
		if dr, ok := dst.(Register); ok {
			asm.arithmeticImmReg(insn, s, dr)
		} else {
			dst.Rex(asm, Register{insn.imm_rm.sub, 0})
			asm.byte(insn.imm_rm.op)
			dst.ModRM(asm, Register{insn.imm_rm.sub, 0})
		}
		asm.int32(uint32(s))
		return
	case Register:
		if dr, ok := dst.(Register); ok {
			asm.arithmeticRegReg(insn, s, dr)
			return
		}
	}

}

func (a *Assembler) Mov(src, dst Operand) {
	a.arithmetic(InstMov, src, dst)
}

func (a *Assembler) Ret() {
	a.byte(0xc3)
}
