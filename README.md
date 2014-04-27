# `gojit` -- pure-golang runtime code-generation

This is the result of my spending the hack day at
[Gophercon 2014](http://gophercon.com) playing with doing JIT from
golang code. This repository contains several packages:

- `gojit`

   Contains the basic JIT support -- allocate executable chunks of
   memory, and convert them into callable golang functions.

- `amd64`

   Contains a simplistic amd64 assembler designed for use with `gojit`

- `bf`

   Contains a just-in-time compiler for
   [Brainfuck](http://esolangs.org/wiki/Brainfuck) that demos the
   above packages

- `gobf`

   Contains a binary that provides a command-line interface to `bf`


## Using

`gobf` can be fetched using

    go get github.com/nelhage/gojit/gobf

And then run as `gobf file.bf`. For some built-in examples:

    $ gobf $GOPATH/src/github.com/nelhage/gojit/bf/test/hello.bf
    Hello World!
    $ gobf $GOPATH/src/github.com/nelhage/gojit/bf/test/hello.bf | gobf $GOPATH/src/github.com/nelhage/gojit/bf/test/rot13.bf
    Uryyb Jbeyq!

## Portability

This code has been tested on `darwin/amd64` and `linux/amd64`. It is
extremely unlikely to work anywhere else.
