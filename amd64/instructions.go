package amd64

type ImmRm struct {
	op  byte
	sub byte
}

type Instruction struct {
	imm_r  byte
	imm_rm ImmRm
	r_rm   byte
	rm_r   byte
}

var (
	InstAdd  = &Instruction{0, ImmRm{0x81, 0}, 0x01, 0x03}
	InstAnd  = &Instruction{0, ImmRm{0x81, 4}, 0x21, 0x23}
	InstCmp  = &Instruction{0, ImmRm{0x81, 7}, 0x39, 0x3B}
	InstLea  = &Instruction{0, ImmRm{0, 0}, 0, 0x8D}
	InstMov  = &Instruction{0xB8, ImmRm{0xc7, 0}, 0x89, 0x8b}
	InstOr   = &Instruction{0, ImmRm{0x81, 1}, 0x09, 0x0B}
	InstSub  = &Instruction{0, ImmRm{0x81, 5}, 0x29, 0x2B}
	InstTest = &Instruction{0, ImmRm{0xF7, 0}, 0x85, 0}
	InstXor  = &Instruction{0, ImmRm{0x81, 6}, 0x31, 0x33}
)
