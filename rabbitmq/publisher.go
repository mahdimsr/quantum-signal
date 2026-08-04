package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	client *RabbitMQClient
}

func NewPublisher(client *RabbitMQClient) *Publisher {
	return &Publisher{client: client}
}

// Publish هر بار یک channel جدید می‌سازد، publish می‌کند و channel را می‌بندد
func (p *Publisher) Publish(ctx context.Context, queueName string, message interface{}) error {
	// ۱. ساخت channel جدید
	ch, err := p.client.NewChannel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}
	defer ch.Close() // مهم: بعد از publish حتماً بسته شود

	// ۲. اطمینان از وجود queue
	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// ۳. تبدیل به JSON
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	log.Printf("📤 Publishing to: %s", queueName)
	log.Printf("📦 Body: %s", string(body))

	// ۴. ارسال
	err = ch.PublishWithContext(
		ctx,
		"",        // candles
		queueName, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	log.Printf("✅ Published successfully")
	return nil
}
