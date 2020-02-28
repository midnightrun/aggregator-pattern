package aggregator

import (
	"time"
)

// Strategy receive an Event process the behaviour of the system relating to the current Aggregation state.
func Strategy(evt SecurityNotification, state Aggregation) (*AggregationNotification, Aggregation) {
	state.Notifications = append(state.Notifications, evt)
	state.LastUpdate = time.Now().UTC()

	if len(state.Notifications) >= 3 || evt.Priority == HIGH {
		return aggregationToNotification(state), Aggregation{}
	}

	return nil, state
}

func StrategyWithoutEvent(state Aggregation) (*AggregationNotification, Aggregation) {
	if time.Now().Add(-30 * time.Second).UTC().Before(state.LastUpdate) {
		return nil, state
	}

	if len(state.Notifications) > 0 {
		return aggregationToNotification(state), Aggregation{}
	}

	return nil, state
}

func aggregationToNotification(state Aggregation) *AggregationNotification {
	return &AggregationNotification{
		Email:         state.Email,
		Notifications: state.Notifications,
	}
}
