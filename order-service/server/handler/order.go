package handler

import (
	"github.com/osmanunal/commerce-oms/order-service/internal/model"
	"github.com/osmanunal/commerce-oms/order-service/internal/service"
	"github.com/osmanunal/commerce-oms/order-service/server/viewmodel"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var orderRequest viewmodel.OrderRequest
	if err := c.BodyParser(&orderRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	order := orderRequest.ToModel(model.Order{})

	err := h.orderService.Create(c.Context(), &order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Order created successfully"})
}

func (h *OrderHandler) Get(c *fiber.Ctx) error {
	orderID := c.Params("id")
	order, err := h.orderService.Get(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var orderResponse viewmodel.OrderResponse
	orderResponse.FromModel(*order)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"order": orderResponse})
}
