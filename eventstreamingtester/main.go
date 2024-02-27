package main

import (
	"eventstreamingtester/config"
	"eventstreamingtester/pkg/broker"
	"eventstreamingtester/pkg/broker/kafka"
	"eventstreamingtester/pkg/broker/rabbitmq"
	"eventstreamingtester/pkg/metrics"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	brokerAddr := os.Getenv("BROKER_ADDRESS")
	role := os.Getenv("ROLE")
	brokerName := os.Getenv("BROKER")

	appConfig, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("Failed to load configurations: %v", err)
	}
	log.Info("Configurations loaded successfully")
	numEvents := appConfig.TestConfig.NumberOfEvents

	log.Infof("NumberOfEvents=%d", numEvents)

	metrics.Metrics = metrics.NewMetricsCollector(numEvents)

	var brokerClient broker.IBroker

	cfg := &broker.ClientConfig{
		MaxRetry:   10000,
		RetryDelay: time.Duration(2) * time.Millisecond,
		Timeout:    time.Duration(20) * time.Second,
		BrokerAddr: brokerAddr,
	}

	log = logrus.New()

	switch brokerName {
	case "kafka":
		brokerClient = kafka.NewClient(cfg, log)
	case "rabbitmq":
		brokerClient = rabbitmq.NewClient(cfg, log)
	default:
		log.Fatalf("Invalid broker: %s", brokerName)
	}

	broker.Run(brokerClient, brokerAddr, role, numEvents)
}
