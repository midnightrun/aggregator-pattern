package aggregator

import (
	"fmt"
	"time"
)

type mockProcessor struct {
	processedNotification SecurityNotification
	processedAggregation  Aggregation
	returnAggregation     Aggregation
	err                   error
}

func (p *mockProcessor) Process(n SecurityNotification, a Aggregation) (Aggregation, error) {
	p.processedNotification = n
	p.processedAggregation = a
	return p.returnAggregation, p.err
}

func (p *mockProcessor) ProcessWithoutEvent(state Aggregation) (*Aggregation, error) {
	fmt.Println("Processing not implemented for type mockAggregationProcessor")

	if time.Now().Add(-4 * time.Hour).Before(state.LastUpdate) {
		return nil, nil
	}

	return &state, nil
}
