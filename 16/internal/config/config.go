package config

import (
	"flag"
	"time"
)

type Config struct {
	Depth      uint
	NumWorkers int
	Timeout    int64
	OutputDir  string
}

func InitConfig() *Config {
	timeout := flag.Int64("t", 10, "timeout seconds per page (default 10)")
	depth := flag.Uint("d", 0, "depth (default 0 for only root page)")
	numWorkers := flag.Int("w", 10, "num of workers(default 10)")
	outputDir := flag.String("o", "./output", "output directory (default ./output)")
	flag.Parse()

	if *numWorkers == 0 {
		*numWorkers = 1
	}

	if *timeout <= 0 {
		*timeout = 10
	}
	return &Config{
		Depth:      *depth,
		NumWorkers: *numWorkers,
		Timeout:    *timeout,
		OutputDir:  *outputDir,
	}
}

func (c *Config) TimeoutToSeconds() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}
