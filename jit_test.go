package gojit

import (
	"testing"
)

func TestRET(t *testing.T) {
	b, e := NewBuffer()
	if e != nil {
		t.Fatalf("NewBuffer: %s", e.Error())
	}
	defer b.Release()
	b.Buf = append(b.Buf, 0xc3)
	b.Call()
}
