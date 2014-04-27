package amd64

type Assembler struct {
	Buf []byte
	Off int
}

func (a *Assembler) byte(b byte) {
	a.Buf[a.Off] = b
	a.Off++
}

func (a *Assembler) int16(i uint16) {
	a.Buf[a.Off] = byte(i & 0xFF)
	a.Buf[a.Off+1] = byte(i >> 8)
	a.Off += 2
}

func (a *Assembler) int32(i uint32) {
	a.Buf[a.Off] = byte(i & 0xFF)
	a.Buf[a.Off+1] = byte(i >> 8)
	a.Buf[a.Off+2] = byte(i >> 16)
	a.Buf[a.Off+3] = byte(i >> 24)
	a.Off += 4
}

func (a *Assembler) int64(i uint64) {
	a.Buf[a.Off] = byte(i & 0xFF)
	a.Buf[a.Off+1] = byte(i >> 8)
	a.Buf[a.Off+2] = byte(i >> 16)
	a.Buf[a.Off+3] = byte(i >> 24)
	a.Buf[a.Off+4] = byte(i >> 32)
	a.Buf[a.Off+5] = byte(i >> 40)
	a.Buf[a.Off+6] = byte(i >> 48)
	a.Buf[a.Off+7] = byte(i >> 56)
	a.Off += 8
}

func (a *Assembler) rex(w, r, x, b bool) {
	var bits byte
	if w {
		bits |= REXW
	}
	if r {
		bits |= REXR
	}
	if x {
		bits |= REXX
	}
	if b {
		bits |= REXB
	}
	if bits != 0 {
		a.byte(PFX_REX | bits)
	}
}

func (a *Assembler) rexBits(lsize, rsize byte, r, x, b bool) {
	if lsize != 0 && rsize != 0 && lsize != rsize {
		panic("mismatched instruction sizes")
	}
	lsize = lsize | rsize
	if lsize == 0 {
		lsize = 64
	}
	a.rex(lsize == 64, r, x, b)
}

func (a *Assembler) modrm(mod, reg, rm byte) {
	a.byte((mod << 6) | (reg << 3) | rm)
}

func (a *Assembler) sib(s, i, b byte) {
	a.byte((s << 6) | (i << 3) | b)
}
