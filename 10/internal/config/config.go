package config

import "flag"

// Config holds the configuration for the sort utility.
type Config struct {
	Column         int
	IsNumeric      bool
	IsReverse      bool
	IsUnique       bool
	IsMonthSort    bool
	IsIgnoreBlanks bool
	IsCheckSorted  bool
	IsHumanNumeric bool

	InputFile string
}

// InitConfig initializes and returns a Config with command-line flags parsed.
func InitConfig() *Config {
	f := Config{}
	flag.IntVar(&f.Column, "k", 1, "sort by column N (1-based)")
	flag.BoolVar(&f.IsNumeric, "n", false, "sort numerically")
	flag.BoolVar(&f.IsReverse, "r", false, "reverse order")
	flag.BoolVar(&f.IsUnique, "u", false, "output unique lines only")
	flag.BoolVar(&f.IsMonthSort, "M", false, "sort by month name")
	flag.BoolVar(&f.IsIgnoreBlanks, "b", false, "ignore trailing blanks")
	flag.BoolVar(&f.IsCheckSorted, "c", false, "check if already sorted")
	flag.BoolVar(&f.IsHumanNumeric, "h", false, "sort by human-readable numeric values")
	flag.Parse()

	f.InputFile = flag.Arg(0)

	return &f
}
