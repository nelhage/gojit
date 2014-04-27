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

func Build(b []byte) func() {
	addr := Addr(b)
	stub := &addr

	return *(*func())(unsafe.Pointer(&stub))
}

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
