package aggregator

import "time"

type Priority int

const (
	LOW = iota
	MEDIUM
	HIGH
)

func (p Priority) String() string {
	return [...]string{"low", "medium", "high"}[p]
}

type SecurityNotification struct {
	Email        string
	Notification string
	Timestamp    time.Time
	Priority     Priority
}

type Aggregation []*SecurityNotification

type AggregationNotification struct {
	Email         string
	Notifications Aggregation
}
