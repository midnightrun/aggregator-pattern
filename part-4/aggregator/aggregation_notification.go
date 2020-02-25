package aggregator

type AggregationNotification struct {
	Email         string
	Notifications []*SecurityNotification
}
