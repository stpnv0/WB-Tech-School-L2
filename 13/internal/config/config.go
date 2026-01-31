package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Fields    string
	Delimiter string
	Separated bool
}

func InitConfig() *Config {
	cfg := Config{}
	flag.StringVar(&cfg.Fields, "f", "", "fields to select (e.g. 1,3-5)")
	flag.StringVar(&cfg.Delimiter, "d", "\t", "delimiter character")
	flag.BoolVar(&cfg.Separated, "s", false, "only output lines containing delimiter")
	flag.Parse()

	if cfg.Fields == "" {
		fmt.Fprintln(os.Stderr, "usage: cut -f fields [-d delimiter] [-s]")
		os.Exit(1)
	}

	return &cfg
}
