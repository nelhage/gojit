// Package gojit contains basic support for writing JITs in golang. It
// contains functions for allocating byte slices in executable memory,
// and converting between such slices and golang function types.
package gojit

import (
	_ "github.com/JamesDunne/gojit/cgo"
	"github.com/edsrzf/mmap-go"
	"reflect"
	"unsafe"
)

type ABI int

// Alloc returns a byte slice of the specified length that is marked
// RWX -- i.e. the memory in it can be both written and executed. This
// is just a simple wrapper around syscall.Mmap.
//
// len most likely needs to be a multiple of PageSize.
func Alloc(len int) ([]byte, error) {
	b, err := mmap.MapRegion(nil, len, mmap.EXEC|mmap.RDWR, mmap.ANON, int64(0))
	return b, err
}

// Release frees a buffer allocated by Alloc
func Release(b []byte) error {
	m := mmap.MMap(b)
	return m.Unmap()
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
	dummy := jitcall
	fn := &struct {
		trampoline uintptr
		jitcode    uintptr
	}{**(**uintptr)(unsafe.Pointer(&dummy)), Addr(b)}

	return *(*func())(unsafe.Pointer(&fn))
}

// BuildCgo is like Build, but the resulting provided code will be
// called by way of the cgo runtime. This has the advantage of being
// much easier and safer to program against (your JIT'd code need only
// conform to your platform's C ABI), at the cost of significant
// overhead for each call into your code.
func BuildCgo(b []byte) func() {
	dummy := cgocall
	fn := &struct {
		trampoline uintptr
		jitcode    uintptr
	}{**(**uintptr)(unsafe.Pointer(&dummy)), Addr(b)}

	return *(*func())(unsafe.Pointer(&fn))
}

// BuildTo converts a byte-slice into an arbitrary-signatured
// function. The out argument should be a pointer to a variable of
// `func' type.
//
// Arguments to the resulting function will be passed to code using a
// hybrid of the GCC and 6c ABIs: The compiled code will receive, via
// the GCC ABI, a single argument, a void* pointing at the beginning
// of the 6c argument frame. For concreteness, on amd64, a
// func([]byte) int would result in %rdi pointing at the 6c stack
// frame, like so:
//
//     24(%rdi) [ return value ]
//     16(%rdi) [  cap(slice)  ]
//     8(%rdi)  [  len(slice)  ]
//     0(%rdi)  [ uint8* data  ]
func BuildTo(b []byte, out interface{}) {
	buildToInternal(b, out, Build)
}

// BuildToCgo is as Build, but uses cgo like BuildCGo
func BuildToCgo(b []byte, out interface{}) {
	buildToInternal(b, out, BuildCgo)
}

func buildToInternal(b []byte, out interface{}, build func([]byte) func()) {
	v := reflect.ValueOf(out)
	if v.Type().Kind() != reflect.Ptr {
		panic("BuildTo: must pass a pointer")
	}
	if v.Elem().Type().Kind() != reflect.Func {
		panic("BuildTo: must pass a pointer to func")
	}

	f := build(b)

	ival := *(*struct {
		typ uintptr
		val uintptr
	})(unsafe.Pointer(&out))

	// Since we know that out has concrete type of *func(..) ...,
	// we know it fits into a word, and thus `val' is just the
	// pointer itself (http://research.swtch.com/interfaces)

	*(*func())(unsafe.Pointer(ival.val)) = f
}

func jitcall()
func cgocall()
