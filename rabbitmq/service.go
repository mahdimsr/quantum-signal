package rabbitmq

import (
	"context"
	"fmt"
	"time"
)

type RabbitMqService struct {
	Publisher *Publisher
	Queue     string
}

type RabbiMqMessage struct {
	Version       int                `json:"version"`
	Price         float64            `json:"price"`
	Symbol        string             `json:"symbol"`
	Side          string             `json:"side"`
	StrategyName  string             `json:"strategy_name"`
	UTCSignalTime int64              `json:"utc_signal_time"`
	StopLoss      float64            `json:"stop_loss"`
	Meta          RabbiMqMessageMeta `json:"meta"`
}

type RabbiMqMessageMeta struct {
	TakeProfit float64                `json:"take_profit"`
	Config     map[string]interface{} `json:"config"`
}

func NewRabbitMqService() (rabbitMqService *RabbitMqService, err error) {

	rabbitURL := "amqp://guest:guest@localhost:5672/"
	queueName := "signals"

	client, err := NewRabbitMQClient(rabbitURL)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	publisher := NewPublisher(client)

	return &RabbitMqService{
		Publisher: publisher,
		Queue:     queueName,
	}, nil
}

func (rabbitService *RabbitMqService) GenerateMessageV1(price, sl float64, symbol, side, strategy string, utcTime int64, meta RabbiMqMessageMeta) RabbiMqMessage {

	return RabbiMqMessage{
		Version:       1,
		Price:         price,
		StopLoss:      sl,
		Symbol:        symbol,
		Side:          side,
		StrategyName:  strategy,
		UTCSignalTime: utcTime,
		Meta:          meta,
	}
}

func (rabbitService *RabbitMqService) SendSignal(message RabbiMqMessage) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := rabbitService.Publisher.Publish(ctx, rabbitService.Queue, message)
	if err != nil {
		fmt.Printf("Failed to publish in rabbit service: %s", err)
	}

}
