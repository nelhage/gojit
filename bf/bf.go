// Package bf implements a JIT compiler for the Brainfuck programming
// language.
package bf

import (
	"bytes"
	"fmt"
	"github.com/nelhage/gojit"
	"github.com/nelhage/gojit/amd64"
	"io"
)

type compiled struct {
	buf   []byte
	code  func([]byte)
	r     func([]byte) (int, error)
	w     func([]byte) (int, error)
	stack []int
}

func (c *compiled) run(b []byte) {
	c.code(b)
}

// %rax is the tape pointer

func jcc(a *amd64.Assembler, cc byte, over func(*amd64.Assembler)) {
	start := a.Off
	a.JccShort(cc, 0)
	base := a.Off
	over(a)
	end := a.Off
	a.Off = start
	if int(int8(end-base)) != end-base {
		panic("jcc: too far!")
	}
	a.JccShort(cc, int8(end-base))
	a.Off = end
}

var knownOpcodes = []byte("+-[]<>,.")

type opcode struct {
	op     byte
	repeat int
}

func optimize(prog []byte) []opcode {
	out := make([]opcode, 0, len(prog)/4)
	for _, b := range prog {
		if bytes.IndexByte(knownOpcodes, b) == -1 {
			continue
		}
		if len(out) > 0 && out[len(out)-1].op == b {
			out[len(out)-1].repeat += 1
		} else {
			out = append(out, opcode{b, 1})
		}
	}
	return out
}

func emitDot(asm *amd64.Assembler, cc *compiled) {
	asm.Push(amd64.Rax)
	asm.Sub(amd64.Imm{48}, amd64.Rsp)
	asm.Mov(amd64.Imm{1}, amd64.Indirect{amd64.Rsp, 16, 64})
	asm.Mov(amd64.Imm{1}, amd64.Indirect{amd64.Rsp, 8, 64})
	asm.Mov(amd64.Rax, amd64.Indirect{amd64.Rsp, 0, 64})
	asm.CallFunc(cc.w)
	asm.Add(amd64.Imm{48}, amd64.Rsp)
	asm.Pop(amd64.Rax)
}

func emitComma(asm *amd64.Assembler, cc *compiled) {
	asm.Push(amd64.Rax)
	asm.Sub(amd64.Imm{48}, amd64.Rsp)
	asm.Mov(amd64.Imm{1}, amd64.Indirect{amd64.Rsp, 16, 64})
	asm.Mov(amd64.Imm{1}, amd64.Indirect{amd64.Rsp, 8, 64})
	asm.Mov(amd64.Rax, amd64.Indirect{amd64.Rsp, 0, 64})
	asm.CallFunc(cc.r)
	asm.Add(amd64.Imm{48}, amd64.Rsp)
	asm.Pop(amd64.Rax)
	asm.Test(amd64.Imm{-1}, amd64.Indirect{amd64.Rsp, -24, 64})
	jcc(asm, amd64.CC_Z, func(asm *amd64.Assembler) {
		asm.Movb(amd64.Imm{0}, amd64.Indirect{amd64.Rax, 0, 8})
	})
}

func emitLbrac(asm *amd64.Assembler, cc *compiled) {
	cc.stack = append(cc.stack, asm.Off)
	asm.Testb(amd64.Imm{0xff}, amd64.Indirect{amd64.Rax, 0, 8})
	asm.JccRel(amd64.CC_Z, gojit.Addr(asm.Buf[asm.Off:]))
}

func emitRbrac(asm *amd64.Assembler, cc *compiled) error {
	if len(cc.stack) == 0 {
		return fmt.Errorf("mismatched []")
	}
	header := cc.stack[len(cc.stack)-1]
	cc.stack = cc.stack[:len(cc.stack)-1]
	asm.JmpRel(gojit.Addr(asm.Buf[header:]))
	end := asm.Off

	asm.Off = header
	asm.Testb(amd64.Imm{0xff}, amd64.Indirect{amd64.Rax, 0, 8})
	asm.JccRel(amd64.CC_Z, gojit.Addr(asm.Buf[end:]))
	asm.Off = end

	return nil
}

// Compile compiles a brainfuck program (represented as a byte slice)
// into a Go function. The function accepts as an argument the tape to
// operate on. The provided Reader and Writer are used to implement
// `,' and `.', respectively.
//
// The compiled code does no bounds-checking on the tape. On EOF or
// other read error, `,' clears the current cell.
func Compile(prog []byte, r io.Reader, w io.Writer) (func([]byte), error) {
	buf, e := gojit.Alloc(4096 * 4)
	if e != nil {
		return nil, e
	}

	cc := &compiled{buf: buf, r: r.Read, w: w.Write}

	asm := &amd64.Assembler{buf, 0}
	asm.Mov(amd64.Indirect{amd64.Rdi, 0, 64}, amd64.Rax)

	opcodes := optimize(prog)

	for _, op := range opcodes {
		switch op.op {
		case '+':
			asm.Addb(amd64.Imm{int32(op.repeat)},
				amd64.Indirect{amd64.Rax, 0, 8})
		case '-':
			asm.Subb(amd64.Imm{int32(op.repeat)},
				amd64.Indirect{amd64.Rax, 0, 8})
		case '<':
			asm.Sub(amd64.Imm{int32(op.repeat)}, amd64.Rax)
		case '>':
			asm.Add(amd64.Imm{int32(op.repeat)}, amd64.Rax)
		default:
			for i := 0; i < op.repeat; i++ {
				switch op.op {
				case '.':
					emitDot(asm, cc)
				case ',':
					emitComma(asm, cc)
				case '[':
					emitLbrac(asm, cc)
				case ']':
					if e := emitRbrac(asm, cc); e != nil {
						return nil, e
					}
				}
			}
		}
	}

	asm.Ret()

	gojit.BuildTo(buf, &cc.code)
	return cc.run, nil
}
