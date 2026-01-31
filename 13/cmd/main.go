package main

import (
	"fmt"
	"os"

	"wb-cut/internal/config"
	"wb-cut/internal/cut"
)

func main() {
	cfg := config.InitConfig()

	c := cut.New(cfg)
	if err := c.Run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
