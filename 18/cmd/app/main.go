package main

import (
	"calendar/internal/config"
	"calendar/internal/repository"
	"calendar/internal/router"
	"calendar/internal/service"
	"calendar/internal/transport/handler"
	"calendar/internal/transport/middleware"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	logger := newLogger(cfg.Logger)
	slog.SetDefault(logger)

	logger.Info("starting application")

	repo := repository.NewRepo()
	serv := service.NewService(repo, logger)
	hand := handler.NewHandler(serv, logger)

	r := router.InitRouter(hand)
	r.Use(
		gin.Recovery(),
		middleware.LoggingMiddleware(logger),
	)

	httpServer := &http.Server{
		Addr:         cfg.HTTP.Addr(),
		Handler:      r,
		ReadTimeout:  cfg.HTTP.Timeout.Read,
		WriteTimeout: cfg.HTTP.Timeout.Write,
		IdleTimeout:  cfg.HTTP.Timeout.Idle,
	}

	go func() {
		logger.Info("http server started", slog.String("addr", cfg.HTTP.Addr()))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	shutdown(httpServer, logger)
}

func newLogger(cfg config.LoggerConfig) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: parseLevel(cfg.Level),
	}

	switch cfg.Format {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func shutdown(server *http.Server, logger *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", slog.Any("error", err))
	}

	logger.Info("server stopped gracefully")
}
