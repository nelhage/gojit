package gojit

import (
	"syscall"
)

type Buffer struct {
	Buf   []byte
	alloc []byte
}

func NewBuffer() (*Buffer, error) {
	b := &Buffer{}
	var err error
	b.alloc, err = syscall.Mmap(-1, 0, 8192,
		syscall.PROT_EXEC|syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANONYMOUS|syscall.MAP_PRIVATE)
	if err != nil {
		return nil, err
	}
	b.Buf = b.alloc[:0]
	return b, nil
}

func (b *Buffer) Call() {
	call(b.alloc)
}

func (b *Buffer) Release() {
	syscall.Munmap(b.alloc)
}

// go:noescape
func call(b []byte)
