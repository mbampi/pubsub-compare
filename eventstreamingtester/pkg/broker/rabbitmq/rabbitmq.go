package rabbitmq

import (
	"context"
	"eventstreamingtester/pkg/broker"
	"eventstreamingtester/pkg/metrics"
	"eventstreamingtester/pkg/model"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

const (
	queueName   = "iot-events"
	consumerTag = "iotConsumer"
)

type Client struct {
	config *broker.ClientConfig
	log    *logrus.Logger
}

func NewClient(cfg *broker.ClientConfig, logger *logrus.Logger) *Client {
	if logger == nil {
		logger = logrus.New()
	}
	return &Client{
		config: cfg,
		log:    logger,
	}
}

// dialRabbitMQ attempts to connect to RabbitMQ with a retry mechanism
func (c *Client) dialRabbitMQ() (*amqp.Connection, error) {
	username := "guest"
	password := "guest"
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s/", username, password, c.config.BrokerAddr)
	var conn *amqp.Connection
	var err error
	for retry := 0; retry < c.config.MaxRetry; retry++ {
		conn, err = amqp.Dial(amqpURL)
		if err == nil {
			return conn, nil
		}
		time.Sleep(c.config.RetryDelay)
	}
	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts", c.config.MaxRetry)
}

// RunProducer publishes the specified number of messages to the queue
func (c *Client) RunProducer(totalMessages int) {
	conn, err := c.dialRabbitMQ()
	if err != nil {
		c.log.Fatalf("Producer: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	ch.Confirm(true)
	if err != nil {
		c.log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, false, true, false, false, nil)
	if err != nil {
		c.log.Fatalf("Failed to declare a queue: %s", err)
	}

	time.Sleep(3 * time.Second)
	c.log.Info("Producer started successfully")

	for i := 0; i < totalMessages; i++ {
		id := int64(i)
		nowTime := time.Now()
		event := model.NewEvent(id, nowTime, 25.1, 60.5, 1013.0)
		body := event.ToJSON()
		err = ch.PublishWithContext(context.Background(), "", queueName, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
		if err != nil {
			c.log.Errorf("Failed to publish a message: %s", err)
			continue
		}
	}
	c.log.Infof("All %d messages sent", totalMessages)
}

// RunConsumer starts consuming messages from the queue
func (c *Client) RunConsumer(totalMessages int) {
	conn, err := c.dialRabbitMQ()
	if err != nil {
		c.log.Fatalf("Consumer: %s", err)
	}
	defer conn.Close()

	time.Sleep(2 * time.Second)

	ch, err := conn.Channel()
	if err != nil {
		c.log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	var retry int
	var msgs <-chan amqp.Delivery
	for retry = 0; retry < c.config.MaxRetry; retry++ {
		msgs, err = ch.Consume(queueName, consumerTag, true, false, false, false, nil)
		if err == nil {
			break
		}
		time.Sleep(c.config.RetryDelay)
	}
	if err != nil {
		c.log.Fatalf("Failed to consume messages: %s", err)
	}

	c.log.Info("Consumer started successfully with retry attempts: ", retry)

	received := 0
	for i := 0; i < totalMessages; i++ {
		select {
		case d, ok := <-msgs:
			if !ok {
				c.log.Errorf("Failed to receive message: %s", err)
				continue
			}
			metrics.Metrics.ReceivedMessageChan <- metrics.ReceivedEvent{
				Timestamp: time.Now(),
				Value:     d.Body,
			}
			received++
		case <-time.After(c.config.Timeout):
			c.log.Info("Consumer stopped by timeout")
			break
		}
	}

	c.log.Infof("%d messages received", received)
}
