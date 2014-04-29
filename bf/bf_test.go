package bf

import (
	"bytes"
	"io"
	"reflect"
	"runtime"
	"testing"
)

var helloWorld = "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."

// http://esolangs.org/wiki/Dbfi
var dbfi = `
>>>+[[-]>>[-]++>+>+++++++[<++++>>++<-]++>>+>+>+++++[>++>++++++<<-]+>>>,<++[[>[
->>]<[>>]<<-]<[<]<+>>[>]>[<+>-[[<+>-]>]<[[[-]<]++<-[<+++++++++>[<->-]>>]>>]]<<
]<]<[[<]>[[>]>>[>>]+[<<]<[<]<+>>-]>[>]+[->>]<<<<[[<<]<[<]+<<[+>+<<-[>-->+<<-[>
+<[>>+<<-]]]>[<+>-]<]++>>-->[>]>>[>>]]<<[>>+<[[<]<]>[[<<]<[<]+[-<+>>-[<<+>++>-
[<->[<<+>>-]]]<[>+<-]>]>[>]>]>[>>]>>]<<[>>+>>+>>]<<[->>>>>>>>]<<[>.>>>>>>>]<<[
>->>>>>]<<[>,>>>]<<[>+>]<<[+<<]<]`

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
		{">+>+[<]", []byte{0, 1, 1}, nil, nil},
		{helloWorld, nil, nil, []byte("Hello World!\n")},
		{dbfi, nil, []byte(helloWorld + "!"), []byte("Hello World!\n")},
	}

	var rd io.Reader
	var wr io.Writer

	for _, tc := range cases {
		rd = bytes.NewBuffer(tc.rd)
		wr = &bytes.Buffer{}

		f, e := Compile([]byte(tc.prog), rd, wr)
		if e != nil {
			t.Errorf("Compile(%v): %s", tc.prog, e.Error())
			continue
		}

		runtime.GC()

		mem := make([]byte, 4096)
		f(mem)
		if tc.mem != nil && !bytes.Equal(mem[:len(tc.mem)], tc.mem) {
			t.Errorf("Compile(%s): %v != %v (expected)",
				tc.prog, mem, tc.mem)
		}
		if tc.wr != nil && !bytes.Equal(tc.wr, wr.(*bytes.Buffer).Bytes()) {
			t.Errorf("Compile(%s): output %v != %v (expected)",
				tc.prog, wr.(*bytes.Buffer).Bytes(), tc.wr)
		}
	}
}

func TestOptimize(t *testing.T) {
	cases := []struct {
		prog string
		ops  []opcode
	}{
		{"+", []opcode{{'+', 1}}},
		{"+++++", []opcode{{'+', 5}}},
		{"++XX+++--<>+", []opcode{{'+', 5}, {'-', 2}, {'<', 1}, {'>', 1}, {'+', 1}}},
	}

	for _, tc := range cases {
		got := optimize([]byte(tc.prog))
		if !reflect.DeepEqual(got, tc.ops) {
			t.Errorf("Optimize(%s): got %v, expect %v",
				tc.prog, got, tc.ops)
		}
	}
}

func BenchmarkCompileHello(b *testing.B) {
	var rw bytes.Buffer
	for i := 0; i < b.N; i++ {
		Compile([]byte(helloWorld), &rw, &rw)
	}
}

func BenchmarkRunHello(b *testing.B) {
	var rw bytes.Buffer
	prog, e := Compile([]byte(helloWorld), &rw, &rw)
	if e != nil {
		b.Fatalf("Compile: %s", e.Error())
	}
	mem := make([]byte, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, _ := range mem {
			mem[j] = 0
		}
		prog(mem)
	}
}

func BenchmarkRunDbfiHello(b *testing.B) {
	var r bytes.Buffer
	var w bytes.Buffer
	input := []byte(helloWorld + "!")
	prog, e := Compile([]byte(dbfi), &r, &r)
	if e != nil {
		b.Fatalf("Compile: %s", e.Error())
	}
	mem := make([]byte, 4096)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, _ := range mem {
			mem[j] = 0
		}
		r.Reset()
		r.Write(input)
		w.Reset()
		prog(mem)
	}
}
