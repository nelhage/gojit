package main

import (
	"bufio"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/nelhage/gojit/bf"
)

func main() {
	var (
		buffer = flag.Bool("buffer", false, "buffer stdout")
	)
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalf("Usage: %s file.bf\n", os.Args[0])
	}

	data, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatalf("Reading %s: %s\n", flag.Arg(0), err.Error())
	}

	out := io.Writer(os.Stdout)
	if *buffer {
		out = bufio.NewWriter(out)
	}

	f, e := bf.Compile(data, os.Stdin, out)
	if e != nil {
		log.Fatalf("compiling: %s", e.Error())
	}
	var memory [4096]byte
	f(memory[:])
}
