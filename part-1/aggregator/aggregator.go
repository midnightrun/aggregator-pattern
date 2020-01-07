package aggregator

func Strategy(evt *SecurityNotification, s Aggregation) (*AggregationNotification, Aggregation) {
	s = append(s, evt)

	if len(s) == 3 || evt.Priority == HIGH {
		return &AggregationNotification{
			Email:         evt.Email,
			Notifications: s,
		}, nil
	}

	return nil, s
}
