# `gojit` -- pure-golang runtime code-generation

This is the result of my spending the hack day at
[Gophercon 2014](http://gophercon.com) playing with doing JIT from
golang code. This respository contains several packages:

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
