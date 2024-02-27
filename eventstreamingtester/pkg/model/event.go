package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// Event represents the event data
type Event struct {
	ID          int64     `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Pressure    float64   `json:"pressure"`
}

// NewEvent creates a new instance of Event
func NewEvent(id int64, timestamp time.Time, temperature float64, humidity float64, pressure float64) *Event {
	return &Event{
		ID:          id,
		Timestamp:   timestamp,
		Temperature: temperature,
		Humidity:    humidity,
		Pressure:    pressure,
	}
}

// String returns the string representation of Event
func (e *Event) String() string {
	return fmt.Sprintf("ID: %d, Timestamp: %s, Temperature: %f, Humidity: %f, Pressure: %f", e.ID, e.Timestamp, e.Temperature, e.Humidity, e.Pressure)
}

// GetID returns the ID of Event from the given JSON string
func GetID(jsonStr string) (int64, error) {
	var event Event
	err := json.Unmarshal([]byte(jsonStr), &event)
	if err != nil {
		return 0, err
	}
	return event.ID, nil
}

// ToJSON returns the JSON string representation of Event
func (e *Event) ToJSON() string {
	jsonStr, _ := json.Marshal(e)
	return string(jsonStr)
}

// FromJSON creates a new instance of Event from the given JSON string
func FromJSON(jsonBytes []byte) (*Event, error) {
	var event Event
	err := json.Unmarshal(jsonBytes, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
