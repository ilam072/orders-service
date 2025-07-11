package rest

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"wb-l0/internal/types/dto"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderId string) (dto.Order, error)
}

type Handler struct {
	log *slog.Logger
	api *fiber.App
	s   OrderService
}

func NewHandler(log *slog.Logger, s OrderService) *Handler {
	api := fiber.New()

	h := &Handler{
		log: log,
		api: api,
		s:   s,
	}
	h.api.Get("/api/order/:id", h.GetOrderHandler)

	return h
}

func (h *Handler) Listen(addr string) error {
	return h.api.Listen(addr)
}
