package gojit

import (
	"reflect"
	"syscall"
	"unsafe"
)

func Alloc(len int) ([]byte, error) {
	b, err := syscall.Mmap(-1, 0, len,
		syscall.PROT_EXEC|syscall.PROT_READ|syscall.PROT_WRITE,
		mapAnon|syscall.MAP_PRIVATE)
	return b, err
}

func Release(b []byte) error {
	return syscall.Munmap(b)
}

func Addr(b []byte) uintptr {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return hdr.Data
}

func Call(b []byte) {
	call(b)
}

type funcStub struct {
	stub uintptr
	code uintptr
}

func Build(b []byte) func() {
	dummy := funcImpl
	stubAddr := **(**uintptr)(unsafe.Pointer(&dummy))

	stub := funcStub{stub: stubAddr, code: Addr(b)}
	dummy2 := &stub

	return *(*func())(unsafe.Pointer(&dummy2))
}

func call(b []byte)
func funcImpl()
