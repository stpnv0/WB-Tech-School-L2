package main

import (
	"fmt"
	"shell/internal/shell"
)

func main() {
	fmt.Println("Simple Shell. Type 'exit' or Ctrl+D to quit.")
	sh := shell.New()
	sh.Run()
}
