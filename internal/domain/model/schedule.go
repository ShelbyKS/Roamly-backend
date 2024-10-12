package model

import "time"

type Event struct {
	Place     string
	StartTime time.Time
	EndTime   time.Time
	Payload   map[string]any
}

type Schedule struct {
	Events  []Event
	Payload map[string]any
}
