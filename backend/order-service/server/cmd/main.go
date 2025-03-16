package main

import (
	"context"
	"fmt"
	"github.com/osmanunal/commerce-oms/order-service/internal/broker"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/osmanunal/commerce-oms/order-service/pkg/config"
	"github.com/osmanunal/commerce-oms/order-service/server/router"

	"github.com/osmanunal/commerce-oms/order-service/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Read()
	db := database.ConnectDB(cfg.DBConfig)

	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(cfg.ServerConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.ServerConfig.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.ServerConfig.IdleTimeout) * time.Second,
	})

	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New())

	productBroker, err := broker.NewRabbitMQProductBroker(&cfg.RabbitMQConfig)
	if err != nil {
		log.Fatalf("RabbitMQ broker oluşturma hatası: %v", err)
	}

	router.Setup(app, db, cfg)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.ServerConfig.Port)); err != nil {
			log.Fatalf("Server çalıştırma hatası: %v", err)
		}
	}()

	log.Printf("Server started on port %d", cfg.ServerConfig.Port)

	// Graceful shutdown için sinyal bekle
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server kapatılıyor...")

	// 5 saniye içinde bağlantıları kapat
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server kapatma hatası: %v", err)
	}

	log.Println("Server başarıyla kapatıldı")
}
