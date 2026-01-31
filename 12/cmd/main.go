package main

import (
	"fmt"
	"io"
	"os"

	"wb-grep/internal/config"
	"wb-grep/internal/grep"
)

func main() {
	cfg := config.InitConfig()

	if cfg.Pattern == "" {
		fmt.Fprintln(os.Stderr, "usage: grep [options] pattern [file]")
		os.Exit(1)
	}

	var source io.Reader
	if cfg.FilePath != "" {
		file, err := os.Open(cfg.FilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		source = file
	} else {
		source = os.Stdin
	}

	g := grep.New(cfg, source)
	if err := g.Run(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
