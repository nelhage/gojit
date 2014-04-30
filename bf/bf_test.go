package bf

import (
	"bytes"
	"github.com/nelhage/gojit/amd64"
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

func TestCompile(t *testing.T) {
	testImplementation(t, Compile)
}

func TestInterpret(t *testing.T) {
	testImplementation(t, Interpret)
}

func testImplementation(t *testing.T,
	prepare func([]byte, io.Reader, io.Writer) (func([]byte), error)) {
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

		f, e := prepare([]byte(tc.prog), rd, wr)
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
		got, _ := optimize([]byte(tc.prog))
		if !reflect.DeepEqual(got, tc.ops) {
			t.Errorf("Optimize(%s): got %v, expect %v",
				tc.prog, got, tc.ops)
		}
	}
}

func TestGC(t *testing.T) {
	var rw bytes.Buffer
	prog, e := Compile([]byte(helloWorld), &rw, &rw)
	if e != nil {
		t.Fatalf("Compile: %s", e.Error())
	}
	var m runtime.MemStats

	for i := 0; i < 1000; i++ {
		runtime.GC()
		runtime.ReadMemStats(&m)

		mem := make([]byte, 2048)
		prog(mem)
	}
}

func BenchmarkCompileHello(b *testing.B) {
	var rw bytes.Buffer
	for i := 0; i < b.N; i++ {
		Compile([]byte(helloWorld), &rw, &rw)
	}
}

func benchmark(b *testing.B,
	prepare func([]byte, io.Reader, io.Writer) (func([]byte), error),
	code, in []byte) {

	var r bytes.Buffer
	var w bytes.Buffer

	prog, e := prepare(code, &r, &r)
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
		if in != nil {
			r.Write(in)
		}
		w.Reset()
		prog(mem)
	}
}

func use_goabi() {
	abi = amd64.GoABI
}

func reset_abi() {
	abi = amd64.CgoABI
}

func BenchmarkCompiledHello(b *testing.B) {
	use_goabi()
	defer reset_abi()
	benchmark(b, Compile, []byte(helloWorld), nil)
}

func BenchmarkCompiledHelloCgo(b *testing.B) {
	benchmark(b, Compile, []byte(helloWorld), nil)
}

func BenchmarkInterpretHello(b *testing.B) {
	benchmark(b, Interpret, []byte(helloWorld), nil)
}

func BenchmarkCompiledDbfiHello(b *testing.B) {
	use_goabi()
	defer reset_abi()
	benchmark(b, Compile, []byte(dbfi), []byte(helloWorld+"!"))
}

func BenchmarkCompiledDbfiHelloCgo(b *testing.B) {
	benchmark(b, Compile, []byte(dbfi), []byte(helloWorld+"!"))
}

func BenchmarkInterpretDbfiHello(b *testing.B) {
	benchmark(b, Interpret, []byte(dbfi), []byte(helloWorld+"!"))
}
