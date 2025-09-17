package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// App version.
const version = "1.0.0"

// App config: port and environment (from flags).
type config struct {
	port int
	env  string
}

// App dependencies: config and logger.
type application struct {
	config config
	logger *slog.Logger
}

func main() {
	// Config instance and flags with defaults.
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Structured logger to stdout.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Build app with config and logger.
	app := &application{
		config: cfg,
		logger: logger,
	}

	// HTTP server with timeouts; error logs via slog at Error level.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start server.
	logger.Info("server is starting", "addr", srv.Addr, "env", cfg.env)

	err := srv.ListenAndServe()
	logger.Error("server is stopping", "error", err.Error())
	os.Exit(1)
}
