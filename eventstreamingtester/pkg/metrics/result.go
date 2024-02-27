package metrics

import (
	"fmt"
	"time"
)

// MetricsResult holds the metrics result
type MetricsResult struct {
	Duration   time.Duration
	Throughput float64 // messages per second
	AvgLatency time.Duration
	ErrorRate  float64
	Errors     int
	Start      time.Time
	End        time.Time
}

// String returns the string representation of MetricsResult
func (m *MetricsResult) String() string {
	return fmt.Sprintf("Duration: %s, Throughput: %f messages/second, AvgLatency: %s, ErrorRate: %f, Errors: %d", m.Duration, m.Throughput, m.AvgLatency, m.ErrorRate, m.Errors)
}
