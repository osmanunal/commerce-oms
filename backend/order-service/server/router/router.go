package router

import (
	"context"
	"log"

	"github.com/osmanunal/commerce-oms/order-service/internal/broker"
	"github.com/osmanunal/commerce-oms/order-service/internal/repository"
	"github.com/osmanunal/commerce-oms/order-service/internal/service"
	"github.com/osmanunal/commerce-oms/order-service/pkg/config"
	"github.com/osmanunal/commerce-oms/order-service/server/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func Setup(app *fiber.App, db *bun.DB, cfg *config.Config) {
	orderRepo := repository.NewOrderRepository(db)

	productBroker, err := broker.NewRabbitMQProductBroker(&cfg.RabbitMQConfig)
	if err != nil {
		log.Fatalf("RabbitMQ broker oluşturma hatası: %v", err)
	}

	orderService := service.NewOrderService(orderRepo, productBroker)

	err = orderService.StartConsuming(context.Background())
	if err != nil {
		log.Fatalf("RabbitMQ consumer başlatma hatası: %v", err)
	}

	orderHandler := handler.NewOrderHandler(*orderService)

	api := app.Group("/api")

	orders := api.Group("/order")
	orders.Post("/", orderHandler.Create)
	orders.Get("/:id", orderHandler.Get)
}
