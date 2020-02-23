package aggregator

func Strategy(evt *SecurityNotification, state Aggregation) (*AggregationNotification, Aggregation) {
	state = append(state, evt)

	if len(state) >= 3 || evt.Priority == HIGH {
		return &AggregationNotification{
			Email:         evt.Email,
			Notifications: state,
		}, state
	}

	return nil, state
}
