package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/osmanunal/commerce-oms/order-service/internal/port"
	"github.com/osmanunal/commerce-oms/order-service/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProductBroker struct {
	conn             *amqp.Connection
	pubCh            *amqp.Channel
	conCh            *amqp.Channel
	stockDecreasePub amqpPublisher
	stockStatusCon   amqpConsumer
	config           *config.RabbitMQConfig
}

type amqpPublisher struct {
	channel      *amqp.Channel
	exchangeName string
	exchangeType string
	routingKey   string
	queueName    string
}

type amqpConsumer struct {
	channel      *amqp.Channel
	exchangeName string
	exchangeType string
	routingKey   string
	queueName    string
}

func NewRabbitMQProductBroker(cfg *config.RabbitMQConfig) (*RabbitMQProductBroker, error) {
	// RabbitMQ bağlantısı kur
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq bağlantı hatası: %w", err)
	}

	// Publisher kanalı aç
	pubCh, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq publisher kanal hatası: %w", err)
	}

	// Consumer kanalı aç
	conCh, err := conn.Channel()
	if err != nil {
		pubCh.Close()
		conn.Close()
		return nil, fmt.Errorf("rabbitmq consumer kanal hatası: %w", err)
	}

	broker := &RabbitMQProductBroker{
		conn:   conn,
		pubCh:  pubCh,
		conCh:  conCh,
		config: cfg,
		stockDecreasePub: amqpPublisher{
			channel:      pubCh,
			exchangeName: cfg.StockDecreaseExchange,
			exchangeType: "direct",
			routingKey:   "stock.decrease",
			queueName:    cfg.StockDecreaseQueue,
		},
		stockStatusCon: amqpConsumer{
			channel:      conCh,
			exchangeName: cfg.StockStatusExchange,
			exchangeType: "direct",
			routingKey:   "stock.status",
			queueName:    cfg.StockStatusQueue,
		},
	}

	// Publisher exchange ve queue tanımla
	err = broker.setupPublisher()
	if err != nil {
		broker.Close()
		return nil, err
	}

	// Consumer exchange ve queue tanımla
	err = broker.setupConsumer()
	if err != nil {
		broker.Close()
		return nil, err
	}

	return broker, nil
}

func (b *RabbitMQProductBroker) setupPublisher() error {
	// StockDecrease Exchange
	err := b.pubCh.ExchangeDeclare(
		b.stockDecreasePub.exchangeName, // exchange name
		b.stockDecreasePub.exchangeType, // exchange type
		true,                            // durable
		false,                           // auto-deleted
		false,                           // internal
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		return fmt.Errorf("stock decrease exchange tanımlama hatası: %w", err)
	}

	// StockDecrease Queue
	_, err = b.pubCh.QueueDeclare(
		b.stockDecreasePub.queueName, // name
		true,                         // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		return fmt.Errorf("stock decrease queue tanımlama hatası: %w", err)
	}

	// StockDecrease Binding
	err = b.pubCh.QueueBind(
		b.stockDecreasePub.queueName,    // queue name
		b.stockDecreasePub.routingKey,   // routing key
		b.stockDecreasePub.exchangeName, // exchange
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		return fmt.Errorf("stock decrease queue binding hatası: %w", err)
	}

	return nil
}

func (b *RabbitMQProductBroker) setupConsumer() error {
	// StockStatus Exchange
	err := b.conCh.ExchangeDeclare(
		b.stockStatusCon.exchangeName, // exchange name
		b.stockStatusCon.exchangeType, // exchange type
		true,                          // durable
		false,                         // auto-deleted
		false,                         // internal
		false,                         // no-wait
		nil,                           // arguments
	)
	if err != nil {
		return fmt.Errorf("stock status exchange tanımlama hatası: %w", err)
	}

	// StockStatus Queue
	_, err = b.conCh.QueueDeclare(
		b.stockStatusCon.queueName, // name
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		return fmt.Errorf("stock status queue tanımlama hatası: %w", err)
	}

	// StockStatus Binding
	err = b.conCh.QueueBind(
		b.stockStatusCon.queueName,    // queue name
		b.stockStatusCon.routingKey,   // routing key
		b.stockStatusCon.exchangeName, // exchange
		false,                         // no-wait
		nil,                           // arguments
	)
	if err != nil {
		return fmt.Errorf("stock status queue binding hatası: %w", err)
	}

	return nil
}

func (b *RabbitMQProductBroker) PublishStockDecrease(ctx context.Context, messages []port.StockDecreaseMessage) error {
	for _, msg := range messages {
		data, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("mesaj marshal hatası: %w", err)
		}

		err = b.pubCh.PublishWithContext(
			ctx,
			b.stockDecreasePub.exchangeName, // exchange
			b.stockDecreasePub.routingKey,   // routing key
			false,                           // mandatory
			false,                           // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         data,
			},
		)
		if err != nil {
			return fmt.Errorf("mesaj publish hatası: %w", err)
		}
	}
	return nil
}

func (b *RabbitMQProductBroker) ConsumeStockStatus(ctx context.Context, handler func(message port.StockStatusMessage) error) error {
	msgs, err := b.conCh.Consume(
		b.stockStatusCon.queueName, // queue
		"",                         // consumer
		false,                      // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)
	if err != nil {
		return fmt.Errorf("rabbitmq consume hatası: %w", err)
	}

	go func() {
		for msg := range msgs {
			var stockStatus port.StockStatusMessage
			err := json.Unmarshal(msg.Body, &stockStatus)
			if err != nil {
				log.Printf("mesaj unmarshal hatası: %v", err)
				msg.Nack(false, true) // requeue
				continue
			}

			err = handler(stockStatus)
			if err != nil {
				log.Printf("handler hata: %v", err)
				msg.Nack(false, true) // requeue
				continue
			}

			msg.Ack(false)
		}
	}()

	return nil
}

func (b *RabbitMQProductBroker) Close() {
	if b.pubCh != nil {
		b.pubCh.Close()
	}
	if b.conCh != nil {
		b.conCh.Close()
	}
	if b.conn != nil {
		b.conn.Close()
	}
}
