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

func asByteInsn(op *Instruction) *Instruction {
	out := *op
	out.Mnemonic += "b"
	if out.imm_r.ok() {
		// XXX: Is this right in general? MOV is right now the
		// only entry in our table with am imm,r
		// specialization, and it's right there. Hell: *are*
		// there even other instructions in amd64 with imm,r
		// specializations?
		out.imm_r = j{out.imm_r.value() & ^byte(8)}
	}
	if out.imm_rm.op.ok() {
		out.imm_rm.op = j{out.imm_rm.op.value() & ^byte(1)}
	}
	if out.r_rm.ok() {
		out.r_rm = j{out.r_rm.value() & ^byte(1)}
	}
	if out.rm_r.ok() {
		out.r_rm = j{out.rm_r.value() & ^byte(1)}
	}

	out.bits = 8

	return &out
}

var (
	InstAdd   = &Instruction{"add", no{}, ImmRm{j{0x81}, 0}, j{0x01}, j{0x03}, 64}
	InstAddb  = asByteInsn(InstAdd)
	InstAnd   = &Instruction{"and", no{}, ImmRm{j{0x81}, 4}, j{0x21}, j{0x23}, 64}
	InstAndb  = asByteInsn(InstAnd)
	InstCmp   = &Instruction{"cmp", no{}, ImmRm{j{0x81}, 7}, j{0x39}, j{0x3B}, 64}
	InstCmpb  = asByteInsn(InstCmp)
	InstOr    = &Instruction{"or", no{}, ImmRm{j{0x81}, 1}, j{0x09}, j{0x0B}, 64}
	InstOrb   = asByteInsn(InstOr)
	InstSub   = &Instruction{"sub", no{}, ImmRm{j{0x81}, 5}, j{0x29}, j{0x2B}, 64}
	InstSubb  = asByteInsn(InstSub)
	InstTest  = &Instruction{"test", no{}, ImmRm{j{0xF7}, 0}, j{0x85}, no{}, 64}
	InstTestb = asByteInsn(InstTest)
	InstXor   = &Instruction{"xor", no{}, ImmRm{j{0x81}, 6}, j{0x31}, j{0x33}, 64}
	InstXorb  = asByteInsn(InstXor)

	InstLea  = &Instruction{"lea", no{}, ImmRm{no{}, 0}, no{}, j{0x8D}, 64}
	InstMov  = &Instruction{"mov", j{0xB8}, ImmRm{j{0xc7}, 0}, j{0x89}, j{0x8b}, 64}
	InstMovb = asByteInsn(InstMov)
)
