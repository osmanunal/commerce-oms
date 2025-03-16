package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/osmanunal/commerce-oms/product-service/internal/port"
	"github.com/osmanunal/commerce-oms/product-service/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQOrderBroker struct {
	conn             *amqp.Connection
	pubCh            *amqp.Channel
	conCh            *amqp.Channel
	stockStatusPub   amqpPublisher
	stockDecreaseCon amqpConsumer
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

func NewRabbitMQOrderBroker(cfg *config.RabbitMQConfig) (*RabbitMQOrderBroker, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq bağlantı hatası: %w", err)
	}

	pubCh, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq publisher kanal hatası: %w", err)
	}

	conCh, err := conn.Channel()
	if err != nil {
		pubCh.Close()
		conn.Close()
		return nil, fmt.Errorf("rabbitmq consumer kanal hatası: %w", err)
	}

	broker := &RabbitMQOrderBroker{
		conn:   conn,
		pubCh:  pubCh,
		conCh:  conCh,
		config: cfg,
		stockStatusPub: amqpPublisher{
			channel:      pubCh,
			exchangeName: cfg.StockStatusExchange,
			exchangeType: "direct",
			routingKey:   "stock.status",
			queueName:    cfg.StockStatusQueue,
		},
		stockDecreaseCon: amqpConsumer{
			channel:      conCh,
			exchangeName: cfg.StockDecreaseExchange,
			exchangeType: "direct",
			routingKey:   "stock.decrease",
			queueName:    cfg.StockDecreaseQueue,
		},
	}

	err = broker.setupPublisher()
	if err != nil {
		broker.Close()
		return nil, err
	}

	err = broker.setupConsumer()
	if err != nil {
		broker.Close()
		return nil, err
	}

	return broker, nil
}

func (b *RabbitMQOrderBroker) setupPublisher() error {
	// StockStatus Exchange
	err := b.pubCh.ExchangeDeclare(
		b.stockStatusPub.exchangeName, // exchange name
		b.stockStatusPub.exchangeType, // exchange type
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
	_, err = b.pubCh.QueueDeclare(
		b.stockStatusPub.queueName, // name
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
	err = b.pubCh.QueueBind(
		b.stockStatusPub.queueName,    // queue name
		b.stockStatusPub.routingKey,   // routing key
		b.stockStatusPub.exchangeName, // exchange
		false,                         // no-wait
		nil,                           // arguments
	)
	if err != nil {
		return fmt.Errorf("stock status queue binding hatası: %w", err)
	}

	return nil
}

func (b *RabbitMQOrderBroker) setupConsumer() error {
	// StockDecrease Exchange
	err := b.conCh.ExchangeDeclare(
		b.stockDecreaseCon.exchangeName, // exchange name
		b.stockDecreaseCon.exchangeType, // exchange type
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
	_, err = b.conCh.QueueDeclare(
		b.stockDecreaseCon.queueName, // name
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
	err = b.conCh.QueueBind(
		b.stockDecreaseCon.queueName,    // queue name
		b.stockDecreaseCon.routingKey,   // routing key
		b.stockDecreaseCon.exchangeName, // exchange
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		return fmt.Errorf("stock decrease queue binding hatası: %w", err)
	}

	return nil
}

func (b *RabbitMQOrderBroker) PublishStockStatus(ctx context.Context, message port.StockStatusMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("mesaj marshal hatası: %w", err)
	}

	err = b.pubCh.PublishWithContext(
		ctx,
		b.stockStatusPub.exchangeName, // exchange
		b.stockStatusPub.routingKey,   // routing key
		false,                         // mandatory
		false,                         // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         data,
		},
	)
	if err != nil {
		return fmt.Errorf("mesaj publish hatası: %w", err)
	}

	return nil
}

func (b *RabbitMQOrderBroker) ConsumeStockDecrease(ctx context.Context, handler func(message port.StockDecreaseMessage) error) error {
	msgs, err := b.conCh.Consume(
		b.stockDecreaseCon.queueName, // queue
		"",                           // consumer
		false,                        // auto-ack
		false,                        // exclusive
		false,                        // no-local
		false,                        // no-wait
		nil,                          // args
	)
	if err != nil {
		return fmt.Errorf("rabbitmq consume hatası: %w", err)
	}

	go func() {
		for msg := range msgs {
			var stockDecrease port.StockDecreaseMessage
			err := json.Unmarshal(msg.Body, &stockDecrease)
			if err != nil {
				log.Printf("mesaj unmarshal hatası: %v", err)
				msg.Nack(false, true) // requeue
				continue
			}

			err = handler(stockDecrease)
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

func (b *RabbitMQOrderBroker) Close() {
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
