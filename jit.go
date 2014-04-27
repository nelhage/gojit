// Package gojit contains basic support for writing JITs in golang. It
// contains functions for allocating byte slices in executable memory,
// and converting between such slices and golang function types.
package gojit

import (
	"reflect"
	"syscall"
	"unsafe"
)

// Alloc returns a byte slice of the specified length that is marked
// RWX -- i.e. the memory in it can be both written and executed. This
// is just a simple wrapper around syscall.Mmap.
//
// len most likely needs to be a multiple of PageSize.
func Alloc(len int) ([]byte, error) {
	b, err := syscall.Mmap(-1, 0, len,
		syscall.PROT_EXEC|syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANON|syscall.MAP_PRIVATE)
	return b, err
}

// Release frees a buffer allocated by Alloc
func Release(b []byte) error {
	return syscall.Munmap(b)
}

// Addr returns the address in memory of a byte slice, as a uintptr
func Addr(b []byte) uintptr {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return hdr.Data
}

// Build returns a nullary golang function that will result in jumping
// into the specified byte slice. The slice should in most cases be a
// slice returned by Alloc, although you could also use syscall.Mmap
// or syscall.Mprotect directly.
func Build(b []byte) func() {
	addr := Addr(b)
	stub := &addr

	return *(*func())(unsafe.Pointer(&stub))
}

// BuildTo converts a byte-slice into an arbitrary-signatured
// function. The out argument should be a pointer to a variable of
// `func' type.
func BuildTo(b []byte, out interface{}) {
	v := reflect.ValueOf(out)
	if v.Type().Kind() != reflect.Ptr {
		panic("BuildTo: must pass a pointer")
	}
	if v.Elem().Type().Kind() != reflect.Func {
		panic("BuildTo: must pass a pointer to func")
	}

	f := Build(b)

	ival := *(*struct {
		typ uintptr
		val uintptr
	})(unsafe.Pointer(&out))

	// Since we know that out has concrete type of *func(..) ...,
	// we know it fits into a word, and thus `val' is just the
	// pointer itself (http://research.swtch.com/interfaces)

	*(*func())(unsafe.Pointer(ival.val)) = f
}
