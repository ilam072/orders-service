package handler

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"wb-l0/internal/types/dto"
	"wb-l0/pkg/e"
	"wb-l0/pkg/logger/sl"
)

type Consumer interface {
	Consume(context.Context) (kafka.Message, error)
	Close() error
}

type Service interface {
	CreateOrder(context.Context, dto.Order) error
}

type OrderConsumerHandler struct {
	c       Consumer
	log     *slog.Logger
	service Service
}

func (h *OrderConsumerHandler) Start(ctx context.Context) error {
	const op = "kafka.handler.Start()"

	log := h.log.With(
		slog.String("op", op),
	)

	for {
		select {
		case <-ctx.Done():
			log.Info("kafka consumer shutting down...")
			return nil
		default:
			message, err := h.c.Consume(ctx)
			if err != nil {
				log.Error("failed to read message", sl.Err(err))
				return e.Wrap(op, err)
			}

			order := dto.Order{}

			if err := json.Unmarshal(message.Value, &order); err != nil {
				log.Error("failed to decode json message to order", sl.Err(err))
				continue
			}

			if err := order.Validate(); err != nil {
				log.Error("failed to validate order", sl.Err(err))
				continue
			}

			if err = h.service.CreateOrder(ctx, order); err != nil {
				log.Error("failed to create order", sl.Err(err))
			}
		}
	}
}
