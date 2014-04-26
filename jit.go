package gojit

import (
	"reflect"
	"syscall"
	"unsafe"
)

func Alloc(len int) ([]byte, error) {
	b, err := syscall.Mmap(-1, 0, len,
		syscall.PROT_EXEC|syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANONYMOUS|syscall.MAP_PRIVATE)
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

// go:noescape
func call(b []byte)
