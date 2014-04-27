package main

import (
	"github.com/nelhage/gojit/bf"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s file.bf\n", os.Args[0])
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Reading %s: %s\n", os.Args[1], err.Error())
	}

	f, e := bf.Compile(data, os.Stdin, os.Stdout)
	if e != nil {
		log.Fatalf("compiling: %s", e.Error())
	}
	var memory [4096]byte
	f(memory[:])
}
