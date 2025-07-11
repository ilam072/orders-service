package main

import (
	"context"
	"log/slog"
	"os"
	"wb-l0/internal/cache"
	"wb-l0/internal/config"
	"wb-l0/internal/repo/postgres"
	"wb-l0/internal/rest"
	"wb-l0/internal/service"
	"wb-l0/pkg/db"
	"wb-l0/pkg/logger/handlers/slogpretty"
)

func main() {
	cfg := config.New()

	pool, err := db.OpenDB(context.Background(), cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	orderRepo := postgres.NewOrderRepo(pool)
	orderService := service.NewOrderService(orderRepo, cache.New())

	log := initLogger()

	rest.NewHandler(log, orderService, cache.New())

	// handler (http, kafka)

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
