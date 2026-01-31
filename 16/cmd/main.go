package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"wb-wget/internal/config"
	"wb-wget/internal/downloader"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := config.InitConfig()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: wb-wget [options] <URL>")
		flag.PrintDefaults()
		os.Exit(1)
	}
	startURL := args[0]

	d := downloader.NewDownloader(
		startURL,
		cfg.OutputDir,
		cfg.Depth,
		cfg.TimeoutToSeconds(),
		cfg.NumWorkers,
	)

	slog.Info("Starting download",
		"url", startURL,
		"depth", cfg.Depth,
		"workers", cfg.NumWorkers,
		"output", cfg.OutputDir,
	)

	if err := d.Run(startURL); err != nil {
		slog.Error("Download failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Download completed", "output", cfg.OutputDir)
}
