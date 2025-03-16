package router

import (
	"context"
	"log"

	"github.com/osmanunal/commerce-oms/product-service/internal/broker"
	"github.com/osmanunal/commerce-oms/product-service/internal/repository"
	"github.com/osmanunal/commerce-oms/product-service/internal/service"
	"github.com/osmanunal/commerce-oms/product-service/pkg/config"
	"github.com/osmanunal/commerce-oms/product-service/server/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func Setup(app *fiber.App, db *bun.DB, cfg *config.Config) {
	productRepo := repository.NewProductRepository(db)

	orderBroker, err := broker.NewRabbitMQOrderBroker(&cfg.RabbitMQConfig)
	if err != nil {
		log.Fatalf("RabbitMQ broker oluşturma hatası: %v", err)
	}

	productService := service.NewProductService(productRepo, orderBroker)

	err = productService.StartConsuming(context.Background())
	if err != nil {
		log.Fatalf("RabbitMQ consumer başlatma hatası: %v", err)
	}

	productHandler := handler.NewProductHandler(*productService)

	api := app.Group("/api")

	products := api.Group("/product")
	products.Post("/", productHandler.Create)
	products.Get("/", productHandler.GetAll)
	products.Get("/:id", productHandler.GetByID)
}
