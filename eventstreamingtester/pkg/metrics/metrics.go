package metrics

import (
	"context"
	"eventstreamingtester/pkg/model"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// Metrics holds the metrics for throughput, latency, error rate, and resource usage
	Metrics *MetricsCollector
	log     = logrus.New()
)

type ReceivedEvent struct {
	Timestamp time.Time
	Value     []byte
}

// MetricsCollector holds the metrics for throughput, latency, error rate, and resource usage
type MetricsCollector struct {
	startTime        time.Time
	endTime          time.Time
	messagesSent     map[int64]time.Time
	messagesReceived []ReceivedEvent
	numEvents        int

	SentMessageChan     chan int64
	ReceivedMessageChan chan ReceivedEvent
}

// NewMetricsCollector initializes a new instance of MetricsCollector
func NewMetricsCollector(numEvents int) *MetricsCollector {
	dataStructureSize := min(numEvents, 1000000)
	return &MetricsCollector{
		// SentMessageChan:     make(chan int64, dataStructureSize),
		ReceivedMessageChan: make(chan ReceivedEvent, dataStructureSize),
		// messagesSent:        make(map[int64]time.Time, dataStructureSize),
		messagesReceived: make([]ReceivedEvent, 0, dataStructureSize),
		numEvents:        numEvents,
	}
}

// collectMetrics collects metrics for throughput, latency, error rate, and resource usage
func (m *MetricsCollector) CollectMetrics(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case sent := <-m.SentMessageChan:
			m.messagesSent[sent] = time.Now()
		case received := <-m.ReceivedMessageChan:
			m.messagesReceived = append(m.messagesReceived, received)
		}
	}
}

// lastMessageSentTime returns the time of the last message sent
func (m *MetricsCollector) lastMessageSentTime() time.Time {
	lastSent := time.Time{}
	for _, sentTime := range m.messagesSent {
		if sentTime.After(lastSent) {
			lastSent = sentTime
		}
	}
	return lastSent
}

// // messagesNotReceived returns the IDs of the messages that were not received
// func (m *MetricsCollector) messagesNotReceived() []int64 {
// 	var notReceived []int64
// 	for sent := range m.messagesSent {
// 		if _, ok := m.messagesReceived[sent]; !ok {
// 			notReceived = append(notReceived, sent)
// 		}
// 	}
// 	return notReceived
// }

func (m *MetricsCollector) MetricsResult() MetricsResult {
	if len(m.messagesReceived) == 0 {
		log.Error("Metrics: No messages received")
		return MetricsResult{}
	}

	firstReceivedEvent := m.messagesReceived[0]
	firstReceived := firstReceivedEvent.Timestamp
	m.startTime = firstReceived

	lastReceived := m.messagesReceived[len(m.messagesReceived)-1].Timestamp
	m.endTime = lastReceived

	return MetricsResult{
		Start:      m.startTime,
		End:        m.endTime,
		Duration:   m.endTime.Sub(m.startTime),
		Throughput: m.Throughput(),
		AvgLatency: m.AvgLatency(),
		Errors:     m.numEvents - len(m.messagesReceived),
		ErrorRate:  m.ErrorRate(),
	}
}

func (m *MetricsCollector) Throughput() float64 {
	totalMessages := len(m.messagesReceived)
	return float64(totalMessages) / m.endTime.Sub(m.startTime).Seconds()
}

func (m *MetricsCollector) AvgLatency() time.Duration {
	var totalLatency time.Duration
	for _, rcvd := range m.messagesReceived {
		event, err := model.FromJSON(rcvd.Value)
		if err != nil {
			log.WithError(err).WithField("value", string(rcvd.Value)).Error("Failed to unmarshal received message")
			continue
		}
		sentTime := event.Timestamp
		receivedTime := rcvd.Timestamp
		totalLatency += receivedTime.Sub(sentTime)
	}
	return totalLatency / time.Duration(len(m.messagesReceived))
}

func (m *MetricsCollector) ErrorRate() float64 {
	errors := m.numEvents - len(m.messagesReceived)
	return float64(errors) / float64(m.numEvents)
}
