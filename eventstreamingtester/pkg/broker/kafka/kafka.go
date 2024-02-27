package kafka

import (
	"eventstreamingtester/pkg/broker"
	"eventstreamingtester/pkg/metrics"
	"eventstreamingtester/pkg/model"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

const topic = "iot-events"

// Client is the Kafka client
type Client struct {
	config *broker.ClientConfig
	log    *logrus.Logger
}

// NewClient creates a new Client with default or provided configurations
func NewClient(cfg *broker.ClientConfig, logger *logrus.Logger) *Client {
	return &Client{
		config: cfg,
		log:    logger,
	}
}

// retryOperation attempts to execute an operation with retry logic
func (c *Client) retryOperation(operation func() error) error {
	var err error
	for retry := 0; retry < c.config.MaxRetry; retry++ {
		if err = operation(); err == nil {
			return nil
		}
		time.Sleep(c.config.RetryDelay)
	}
	return err
}

// RunProducer sends messages to the Kafka topic at a specified rate
func (c *Client) RunProducer(totalMessages int) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true

	err := c.retryOperation(func() error {
		producer, err := sarama.NewSyncProducer([]string{c.config.BrokerAddr}, kafkaConfig)
		if err != nil {
			return err
		}
		defer producer.Close()

		c.log.Info("Producer started successfully")

		for i := 0; i < totalMessages; i++ {
			id := int64(i)
			nowTime := time.Now()
			event := model.NewEvent(id, nowTime, 25.1, 60.5, 1013.0)

			message := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(event.ToJSON()),
			}
			if _, _, err := producer.SendMessage(message); err != nil {
				c.log.Errorf("Failed to send message: %v", err)
			}
		}

		c.log.Info("All messages sent")
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}
}

// RunConsumer consumes messages from the Kafka topic
func (c *Client) RunConsumer(numMessages int) {
	kafkaConfig := sarama.NewConfig()

	err := c.retryOperation(func() error {
		consumer, err := sarama.NewConsumer([]string{c.config.BrokerAddr}, kafkaConfig)
		if err != nil {
			return err
		}
		defer consumer.Close()

		partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
		if err != nil {
			return err
		}
		defer partitionConsumer.Close()

		c.log.Info("Consumer started successfully")

		for i := 0; i < numMessages; i++ {
			select {
			case <-time.After(c.config.Timeout):
				c.log.Info("Consumer stopped by timeout")
				return nil
			case msg := <-partitionConsumer.Messages():
				metrics.Metrics.ReceivedMessageChan <- metrics.ReceivedEvent{
					Timestamp: time.Now(),
					Value:     msg.Value,
				}
			}
		}
		return nil
	})

	if err != nil {
		c.log.Fatalf("Failed to start Sarama consumer: %v", err)
	}
}
