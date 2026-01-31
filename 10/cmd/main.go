package main

import (
	"bufio"
	"io"
	"os"
	"wb-sort/internal/config"
	"wb-sort/internal/parser"
	"wb-sort/internal/sorter"
)

func main() {
	cfg := config.InitConfig()
	source := getSource(cfg.InputFile)

	p := parser.NewParser(cfg, source)
	s := sorter.NewSorter(cfg, p)

	lines, err := s.Sort()
	if err != nil {
		panic(err)
	}
	if err := write(lines, os.Stdout); err != nil {
		panic(err)
	}
}

func getSource(filepath string) io.Reader {
	if filepath != "" {
		f, err := os.Open(filepath)
		if err != nil {
			panic(err)
		}
		return f
	}

	return os.Stdin
}

func write(lines []string, w io.Writer) error {
	bw := bufio.NewWriter(w)
	for _, line := range lines {
		if _, err := bw.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return bw.Flush()
}
