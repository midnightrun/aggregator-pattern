package aggregator

import "time"

type Aggregation struct {
	Notifications []*SecurityNotification
	LastUpdate    time.Time
}
