package parser

import (
	"bufio"
	"errors"
	"io"
	"wb-sort/internal/config"
)

// ErrSource is returned when the source reader is nil.
var ErrSource = errors.New("source is nil")

// Parser reads and parses lines from an input source.
type Parser struct {
	cfg    *config.Config
	source io.Reader
}

// NewParser creates a new Parser with the given configuration and source.
func NewParser(cfg *config.Config, source io.Reader) *Parser {
	return &Parser{
		cfg:    cfg,
		source: source,
	}
}

// Parse reads all lines from the source and returns them as a slice of strings.
func (p *Parser) Parse() ([]string, error) {
	if p.source == nil {
		return nil, ErrSource
	}

	scanner := bufio.NewScanner(p.source)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
