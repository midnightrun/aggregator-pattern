package aggregator

import (
	"fmt"
)

type PublishingProcessor struct{}

func (pp PublishingProcessor) Process(evt *SecurityNotification, existingState Aggregation) (Aggregation, error) {
	notification, newState := Strategy(evt, existingState)
	if notification == nil {
		return newState, nil
	}

	fmt.Printf("publishing new event for %s priority %s", evt.Email, evt.Priority)
	return newState, nil
}

type mockProcessor struct {
	processedNotification *SecurityNotification
	processedAggregation  Aggregation
	returnAggregation     Aggregation
	err                   error
}

func (p *mockProcessor) Process(n *SecurityNotification, a Aggregation) (Aggregator, error) {
	p.processedNotification = n
	p.processedAggregation = a
	return p.returnAggregation, p.err
}
