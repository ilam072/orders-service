package rest

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetOrderHandler(ctx *fiber.Ctx) error {
	orderId := ctx.Params("id")

	order, ok := h.cache.Get(orderId)
	if ok {
		return ctx.Status(fiber.StatusOK).JSON(order)
	}

	order, err := h.s.GetOrder(ctx.Context(), orderId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			errorResponse("something went wrong, try again later"))
	}

	return ctx.Status(fiber.StatusOK).JSON(order)
}
