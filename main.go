package main

import (
	"canary/internal/config"
	"canary/internal/database"
	"canary/internal/modules/logs"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	db, err := database.Connect(cfg.Database.Filepath)
	if err != nil {
		return err
	}
	logger.Info("database connected successfully")
	//nolint:errcheck
	defer db.Close()

	views, err := template.ParseGlob("views/*/*.html")
	if err != nil {
		return fmt.Errorf("failed to load views: %w", err)
	}

	mux := http.NewServeMux()
	logs.InitModule(db, views, mux)

	address := cfg.Server.Address()
	//nolint:exhaustruct
	server := &http.Server{
		Handler:           mux,
		Addr:              address,
		ReadTimeout:       cfg.Server.Timeout,
		WriteTimeout:      cfg.Server.Timeout,
		IdleTimeout:       cfg.Server.Timeout,
		ReadHeaderTimeout: cfg.Server.Timeout,
	}

	// serve public files.
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public", fs))

	logger.Info("starting server", "address", address)
	return server.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
