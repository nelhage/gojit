package amd64

import (
	"reflect"
	"unsafe"
)

func (a *Assembler) CallFunc(f interface{}) {
	switch a.ABI {
	case CgoABI:
		a.CallFuncCgo(f)
	case GoABI:
		a.CallFuncGo(f)
	default:
		panic("bad ABI")
	}
}

// CallFuncGo assembles a call directly to the go function 'f'. No stack
// swtitching or other setup is performed.
func (a *Assembler) CallFuncGo(f interface{}) {
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

// CallFuncCgo assembles a sequence to call into the go function 'f',
// using the cgo callback runtime interface. Prior to CallFunc, a
// 6c/6g-layout stack frame should be set up containing arguments and
// return slots for f. The call will be effected by way of
// runtime.cgocallback_gofunc, involving a stack switch back to the
// goroutine stack.
//
// All registers are caller-save in the 6c ABI, and so all registers
// should be assumed clobbered across a CallFunc.
func (a *Assembler) CallFuncCgo(f interface{}) {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		panic("CallFunc: Can't call non-func")
	}
	ival := *(*struct {
		typ uintptr
		fun uintptr
	})(unsafe.Pointer(&f))

	// runtime·cgocallback_gofunc(f, frame, framesize)
	// plan9 ABI

	framesize := calculateFramesize(f)

	a.Sub(Imm{24}, Rsp)
	a.Mov(Imm{int32(framesize)}, Indirect{Rsp, 16, 64})
	a.Lea(Indirect{Rsp, 24, 64}, Rax)
	a.Mov(Rax, Indirect{Rsp, 8, 64})
	a.MovAbs(uint64(ival.fun), Rax)
	a.Mov(Rax, Indirect{Rsp, 0, 64})
	a.MovAbs(uint64(get_runtime_cgocallback_gofunc()), Rax)
	a.Call(Rax)
	a.Add(Imm{24}, Rsp)
}

func calculateFramesize(f interface{}) uintptr {
	t := reflect.TypeOf(f)
	s := uintptr(0)

	for i := 0; i < t.NumIn(); i++ {
		s += t.In(i).Size()
	}
	for i := 0; i < t.NumOut(); i++ {
		s += t.Out(i).Size()
	}
	return s
}

func get_runtime_cgocallback_gofunc() uintptr
