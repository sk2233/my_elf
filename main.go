package main

import "my_elf/c"

func main() {
	//asm.Compile("source/asm/read.asm", "source/build/obj")
	c.Compilation("source/c/test.c", "source/build/asm.asm")
}
