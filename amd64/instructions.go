package amd64

type maybeByte interface {
	ok() bool
	value() byte
}

type j struct {
	val byte
}

func (j j) ok() bool    { return true }
func (j j) value() byte { return j.val }

type no struct{}

func (n no) ok() bool    { return false }
func (n no) value() byte { panic("no{}.value()!") }

type ImmRm struct {
	op  maybeByte
	sub byte
}

type Instruction struct {
	Mnemonic string
	imm_r    maybeByte
	imm_rm   ImmRm
	r_rm     maybeByte
	rm_r     maybeByte
	bits     byte
}

var (
	InstAdd   = &Instruction{"add", no{}, ImmRm{j{0x81}, 0}, j{0x01}, j{0x03}, 64}
	InstAddb  = &Instruction{"addb", no{}, ImmRm{j{0x80}, 0}, j{0x00}, j{0x03}, 64}
	InstAnd   = &Instruction{"and", no{}, ImmRm{j{0x81}, 4}, j{0x21}, j{0x23}, 64}
	InstCmp   = &Instruction{"cmp", no{}, ImmRm{j{0x81}, 7}, j{0x39}, j{0x3B}, 64}
	InstLea   = &Instruction{"lea", no{}, ImmRm{no{}, 0}, no{}, j{0x8D}, 64}
	InstMov   = &Instruction{"mov", j{0xB8}, ImmRm{j{0xc7}, 0}, j{0x89}, j{0x8b}, 64}
	InstMovb  = &Instruction{"movb", j{0XB0}, ImmRm{j{0xc6}, 0}, j{0x88}, j{0x8a}, 8}
	InstOr    = &Instruction{"or", no{}, ImmRm{j{0x81}, 1}, j{0x09}, j{0x0B}, 64}
	InstSub   = &Instruction{"sub", no{}, ImmRm{j{0x81}, 5}, j{0x29}, j{0x2B}, 64}
	InstSubb  = &Instruction{"subb", no{}, ImmRm{j{0x80}, 5}, j{0x28}, j{0x2a}, 8}
	InstTest  = &Instruction{"test", no{}, ImmRm{j{0xF7}, 0}, j{0x85}, no{}, 64}
	InstTestb = &Instruction{"testb", no{}, ImmRm{j{0xF6}, 0}, j{0x84}, no{}, 8}
	InstXor   = &Instruction{"xor", no{}, ImmRm{j{0x81}, 6}, j{0x31}, j{0x33}, 64}
)
