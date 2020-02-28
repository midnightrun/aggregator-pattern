package aggregator

import "time"

// Aggregation represents an aggregation of SecurityNotification events.
type Aggregation struct {
	Email         string
	Notifications []SecurityNotification
	LastUpdate    time.Time
}

