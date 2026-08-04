package rabbitmq

import (
	"fmt"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn *amqp.Connection
	url  string
	mu   sync.Mutex
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	log.Println("✅ Connected to RabbitMQ")
	return &RabbitMQClient{conn: conn, url: url}, nil
}

// NewChannel هر بار یک channel جدید می‌سازد (مطمئن‌ترین روش)
func (c *RabbitMQClient) NewChannel() (*amqp.Channel, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil || c.conn.IsClosed() {
		// reconnect
		conn, err := amqp.Dial(c.url)
		if err != nil {
			return nil, fmt.Errorf("failed to reconnect: %w", err)
		}
		c.conn = conn
		log.Println("🔄 Reconnected to RabbitMQ")
	}

	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return ch, nil
}

func (c *RabbitMQClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil && !c.conn.IsClosed() {
		return c.conn.Close()
	}
	return nil
}
