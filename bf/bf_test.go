package main

import (
	"bytes"
	"io"
	"testing"
)

func TestSimple(t *testing.T) {
	cases := []struct {
		prog   string
		mem    []byte
		rd, wr []byte
	}{
		{"++", []byte{2, 0}, nil, nil},
		{"++---", []byte{0xff, 0}, nil, nil},
		{"+>+>+", []byte{1, 1, 1, 0}, nil, nil},
		{".", nil, nil, []byte{0}},
		{"++++.>++.>.", []byte{4, 2, 0}, nil, []byte{4, 2, 0}},
		{",", []byte{55}, []byte{55}, nil},
		{",>,>,>,,", []byte{55, 44, 33, 11}, []byte{55, 44, 33, 22, 11}, nil},
		{"+,", []byte{0}, []byte{}, nil},
		{"++++[-]", []byte{0}, nil, nil},
		{"++++[>+++<-]", []byte{0, 12}, nil, nil},
		{"+++[>+++[>+++<-]<-]", []byte{0, 0, 27}, nil, nil},
	}

	var rd io.Reader
	var wr io.Writer

	for _, tc := range cases {
		rd = bytes.NewBuffer(tc.rd)
		wr = &bytes.Buffer{}

		f, e := compile([]byte(tc.prog), rd, wr)
		if e != nil {
			t.Errorf("compile(%v): %s", tc.prog, e.Error())
			continue
		}

		mem := make([]byte, len(tc.mem))
		f(mem)
		if !bytes.Equal(mem, tc.mem) {
			t.Errorf("compile(%s): %v != %v (expected)",
				tc.prog, mem, tc.mem)
		}
		if tc.wr != nil && !bytes.Equal(tc.wr, wr.(*bytes.Buffer).Bytes()) {
			t.Errorf("compile(%s): output %v != %v (expected)",
				tc.prog, wr.(*bytes.Buffer).Bytes(), tc.wr)
		}
	}
}
