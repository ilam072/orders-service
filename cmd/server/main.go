package main

import (
	"context"
	"log/slog"
	"os"
	"wb-l0/internal/cache"
	"wb-l0/internal/config"
	"wb-l0/pkg/db"
	"wb-l0/pkg/logger/handlers/slogpretty"
)

func main() {
	cfg := config.New()

	log := initLogger()

	db, err := db.OpenDB(context.Background(), cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	cache := cache.New()

	// repo

	// service

	// order consumer
}

func initLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	h := opts.NewPrettyHandler(os.Stdout)

	return slog.New(h)
}
