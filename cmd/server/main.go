package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"wb-l0/internal/broker/kafka/consumer"
	"wb-l0/internal/broker/kafka/handler"
	"wb-l0/internal/cache"
	"wb-l0/internal/config"
	"wb-l0/internal/repo/postgres"
	"wb-l0/internal/rest"
	"wb-l0/internal/service"
	"wb-l0/pkg/db"
	"wb-l0/pkg/logger/handlers/slogpretty"
	"wb-l0/pkg/logger/sl"
)

func main() {
	cfg := config.New()

	pool, err := db.OpenDB(context.Background(), cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	log := initLogger()

	orderRepo := postgres.NewOrderRepo(pool)
	cache := cache.New()
	orderService := service.NewOrderService(orderRepo, cache)

	kafkaConsumer := consumer.New(
		cfg.KafkaConfig.Topic,
		cfg.KafkaConfig.GroupID,
		cfg.KafkaConfig.Brokers...,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	orderConsumerHandler := handler.NewOrderConsumerHandler(
		log,
		kafkaConsumer,
		orderService,
	)

	go func() {
		if err := orderConsumerHandler.Start(ctx); err != nil {
			log.Error("failed to start consumer", sl.Err(err))
		}
	}()

	h := rest.NewHandler(log, orderService)
	go func() {
		if err := h.Listen(cfg.ServerConfig.Address()); err != nil {
			log.Error("failed to start server", sl.Err(err))
			cancel()
		}
	}()

	<-sigs
	log.Info("shutting down...")
	cancel()

	if err := h.Shutdown(); err != nil {
		log.Error("failed to shutdown server", sl.Err(err))
	}

	if err := kafkaConsumer.Close(); err != nil {
		log.Error("failed to close kafka consumer", sl.Err(err))
	}

	log.Info("application stopped")
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
