package aggregator

import (
	"time"
)

// Strategy receive an Event process the behaviour of the system relating to the current Aggregation state.
func Strategy(evt *SecurityNotification, state Aggregation) (*AggregationNotification, Aggregation) {
	state.Notifications = append(state.Notifications, evt)

	if len(state.Notifications) >= 3 || evt.Priority == HIGH {
		state.Notifications = []*SecurityNotification{}
		return &AggregationNotification{
			Email:         evt.Email,
			Notifications: state.Notifications,
		}, Aggregation{}
	}

	return nil, state
}

func StrategyWithoutEvent(correlationId string, state Aggregation) (*AggregationNotification, Aggregation) {
	if time.Now().Add(-3 * time.Hour).UTC().Before(state.LastUpdate) {
		return &AggregationNotification{

			Email:         correlationId,
			Notifications: state.Notifications,
		}, state
	}

	return nil, state
}
