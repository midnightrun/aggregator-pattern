package aggregator

import (
	"time"
)

func Strategy(evt *SecurityNotification, state Aggregation) (*AggregationNotification, Aggregation) {
	state.Notifications = append(state.Notifications, evt)

	if len(state.Notifications) >= 3 || evt.Priority == HIGH {
		state.Notifications = []*SecurityNotification{}
		return &AggregationNotification{
			Email:         evt.Email,
			Notifications: state.Notifications,
		}, state
	}

	return nil, state
}

func StrategyWithoutEvent(correlationId string, state Aggregation) (*AggregationNotification, Aggregation) {
	if state.LastUpdate < time.Now().UTC() {
		return &AggregationNotification{

			Email:         correlationId,
			Notifications: state.Notifications,
		}, state
	}

	return nil, state
}
