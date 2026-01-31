package main

import (
	"fmt"
	"unpack/unpacker"
)

func main() {
	s := `\12`
	fmt.Printf("Исходная строка: %s\n", s)
	res, err := unpacker.Unpack(s)
	if err != nil {
		panic(err)
	}

	fmt.Printf("res: %s\n", res)
}
