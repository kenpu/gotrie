package main

import (
	"fmt"
	. "gotrie"
)

func main() {
	a := uint64(0xD8)
	b := uint64(0x08)
	fmt.Printf("%s\n%s\n", Uint64_string(a), Uint64_string(b))
	fmt.Printf("%d vs %d\n", PopCount(a^b), PopCountPartial(a^b, 7))

	fmt.Printf("%s\n%s\n", Uint64_string(a), Uint64_string(a<<3))
}
