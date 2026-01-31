// Package config parses command line arguments and validates configuration.
package config

import (
	"flag"
	"fmt"
	"os"
)

// Config holds the configuration for the grep utility.
type Config struct {
	FilePath string
	Pattern  string

	After      int  // -A N
	Before     int  // -B N
	Context    int  // -C N
	CountOnly  bool // -c
	IgnoreCase bool // -i
	Invert     bool // -v
	Fixed      bool // -F
	LineNum    bool // -n
}

// InitConfig initializes and returns a Config with command-line flags parsed.
func InitConfig() *Config {
	cfg := Config{}
	flag.IntVar(&cfg.After, "A", 0, "print N lines after match")
	flag.IntVar(&cfg.Before, "B", 0, "print N lines before match")
	flag.IntVar(&cfg.Context, "C", 0, "print N lines around match")
	flag.BoolVar(&cfg.CountOnly, "c", false, "print only count of matching lines")
	flag.BoolVar(&cfg.IgnoreCase, "i", false, "ignore case distinction")
	flag.BoolVar(&cfg.Invert, "v", false, "invert match")
	flag.BoolVar(&cfg.Fixed, "F", false, "pattern is fixed string")
	flag.BoolVar(&cfg.LineNum, "n", false, "print line number with output lines")
	flag.Parse()

	if cfg.After < 0 || cfg.Before < 0 || cfg.Context < 0 {
		fmt.Fprintln(os.Stderr, "values for -A, -B, -C cannot be negative")
		os.Exit(2)
	}

	cfg.Pattern = flag.Arg(0)
	if flag.NArg() > 1 {
		cfg.FilePath = flag.Arg(1)
	}

	if cfg.Context > 0 {
		cfg.After = cfg.Context
		cfg.Before = cfg.Context
	}

	return &cfg
}
