// Package amd64 implements a simple amd64 assembler.
package amd64

import "github.com/nelhage/gojit"

type ABI int

const (
	CgoABI ABI = iota
	GoABI
)

// Assembler implements a simple amd64 assembler. All methods on
// Assembler will emit code to Buf[Off:] and advances Off. Buf will
// never be reallocated, and attempts to assemble off the end of Buf
// will panic.
type Assembler struct {
	Buf []byte
	Off int
	ABI ABI
}

func (a *Assembler) BuildTo(out interface{}) {
	switch a.ABI {
	case CgoABI:
		gojit.BuildToCgo(a.Buf, out)
	case GoABI:
		gojit.BuildTo(a.Buf, out)
	default:
		panic("bad ABI")
	}
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

func (a *Assembler) rel32(addr uintptr) {
	off := uintptr(addr) - gojit.Addr(a.Buf[a.Off:]) - 4
	if uintptr(int32(off)) != off {
		panic("call rel: target out of range")
	}
	a.int32(uint32(off))
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
