package model_test

import (
	"fmt"
	"time"

	"eventstreamingtester/pkg/model"
)

func ExampleEvent_ToJSON() {
	t := time.Date(2023, time.September, 6, 0, 0, 0, 0, time.UTC)
	event := model.NewEvent(1, t, 25.5, 60.0, 1013.25)
	jsonStr := event.ToJSON()
	fmt.Println(jsonStr)
	// Output: {"id":1,"timestamp":"2023-09-06T00:00:00Z","temperature":25.5,"humidity":60,"pressure":1013.25}
}

func ExampleFromJSON() {
	jsonStr := `{"id":1,"timestamp":"2023-09-06T00:00:00Z","temperature":25.5,"humidity":60,"pressure":1013.25}`
	event, err := model.FromJSON([]byte(jsonStr))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(event)
	// Output: ID: 1, Timestamp: 2023-09-06 00:00:00 +0000 UTC, Temperature: 25.500000, Humidity: 60.000000, Pressure: 1013.250000
}
