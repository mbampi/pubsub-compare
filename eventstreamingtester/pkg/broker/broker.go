package broker

import (
	"context"
	"eventstreamingtester/pkg/metrics"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// ClientConfig holds configuration for the Broker client
type ClientConfig struct {
	MaxRetry   int
	RetryDelay time.Duration
	Timeout    time.Duration
	BrokerAddr string
}

// IBroker interface for the broker
type IBroker interface {
	RunProducer(totalMessages int)
	RunConsumer(totalMessages int)
}

// Run starts the broker client
func Run(broker IBroker, brokerAddr string, role string, numEvents int) {
	log := logrus.New()
	log.Out = os.Stdout

	log.Info("Starting Broker Client: ", time.Now().String())

	switch role {
	case "producer":
		broker.RunProducer(numEvents)

	case "consumer":
		// Start the metrics collector in a separate goroutine
		ctx, cancel := context.WithCancel(context.Background())
		go metrics.Metrics.CollectMetrics(ctx)

		broker.RunConsumer(numEvents)
		cancel()
		res := metrics.Metrics.MetricsResult()
		log.Infof("Metrics: %+v", res.String())

	default:
		log.Fatalf("Invalid role: %s", role)
	}

	log.Info("Broker Client stopped successfully")
}
