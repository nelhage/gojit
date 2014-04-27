package main

import (
	"fmt"
	"github.com/nelhage/gojit"
	"github.com/nelhage/gojit/amd64"
	"io"
	"io/ioutil"
	"log"
	"os"
)

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

func compile(prog []byte, r io.Reader, w io.Writer) (func([]byte), error) {
	buf, e := gojit.Alloc(4096 * 4)
	if e != nil {
		return nil, e
	}

	asm := &amd64.Assembler{buf, 0}
	asm.Mov(amd64.Indirect{amd64.Rsp, 0x8, 64}, amd64.Rax)

	var stack []int

	for _, b := range prog {
		switch b {
		case '+':
			asm.Incb(amd64.Indirect{amd64.Rax, 0, 8})
		case '-':
			asm.Decb(amd64.Indirect{amd64.Rax, 0, 8})
		case '<':
			asm.Dec(amd64.Rax)
		case '>':
			asm.Inc(amd64.Rax)
		case '.':
			asm.Push(amd64.Rax)
			asm.Sub(amd64.Imm{48}, amd64.Rsp)
			asm.Mov(amd64.Imm{1}, amd64.Indirect{amd64.Rsp, 16, 64})
			asm.Mov(amd64.Imm{1}, amd64.Indirect{amd64.Rsp, 8, 64})
			asm.Mov(amd64.Rax, amd64.Indirect{amd64.Rsp, 0, 64})
			asm.CallFunc(w.Write)
			asm.Add(amd64.Imm{48}, amd64.Rsp)
			asm.Pop(amd64.Rax)
		case ',':
			asm.Push(amd64.Rax)
			asm.Sub(amd64.Imm{48}, amd64.Rsp)
			asm.Mov(amd64.Imm{1}, amd64.Indirect{amd64.Rsp, 16, 64})
			asm.Mov(amd64.Imm{1}, amd64.Indirect{amd64.Rsp, 8, 64})
			asm.Mov(amd64.Rax, amd64.Indirect{amd64.Rsp, 0, 64})
			asm.CallFunc(r.Read)
			asm.Add(amd64.Imm{48}, amd64.Rsp)
			asm.Pop(amd64.Rax)
			asm.Test(amd64.Imm{-1}, amd64.Indirect{amd64.Rsp, -24, 64})
			jcc(asm, amd64.CC_Z, func(asm *amd64.Assembler) {
				asm.Movb(amd64.Imm{0}, amd64.Indirect{amd64.Rax, 0, 8})
			})
		case '[':
			stack = append(stack, asm.Off)
			asm.Testb(amd64.Imm{0xff}, amd64.Indirect{amd64.Rax, 0, 8})
			asm.JccRel(amd64.CC_Z, gojit.Addr(asm.Buf[asm.Off:]))
		case ']':
			if len(stack) == 0 {
				return nil, fmt.Errorf("mismatched []")
			}
			header := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			asm.JmpRel(gojit.Addr(asm.Buf[header:]))
			end := asm.Off

			asm.Off = header
			asm.Testb(amd64.Imm{0xff}, amd64.Indirect{amd64.Rax, 0, 8})
			asm.JccRel(amd64.CC_Z, gojit.Addr(asm.Buf[end:]))
			asm.Off = end
		}
	}

	asm.Ret()

	var f func([]byte)
	gojit.BuildTo(buf, &f)
	return f, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s file.bf\n", os.Args[0])
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Reading %s: %s\n", os.Args[1], err.Error())
	}

	f, e := compile(data, os.Stdin, os.Stdout)
	if e != nil {
		log.Fatalf("compiling: %s", e.Error())
	}
	var memory [4096]byte
	f(memory[:])
}
