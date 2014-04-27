package main

import (
	"bytes"
	"testing"
)

func TestSimple(t *testing.T) {
	cases := []struct {
		prog string
		out  []byte
	}{
		{"++", []byte{2, 0}},
		{"++---", []byte{0xff, 0}},
		{"+>+>+", []byte{1, 1, 1, 0}},
	}

	for _, tc := range cases {
		f, e := compile([]byte(tc.prog))
		if e != nil {
			t.Errorf("compile(%v): %s", tc.prog, e.Error())
		} else {
			mem := make([]byte, len(tc.out))
			f(mem)
			if !bytes.Equal(mem, tc.out) {
				t.Errorf("compile(%s): %v != %v (expected)",
					tc.prog, mem, tc.out)
			}
		}
	}
}
