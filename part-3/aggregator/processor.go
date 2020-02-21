package aggregator

import (
	"fmt"
)

// Processor handles incoming notifications and previous state.
//
// Executed within the context of a transaction, takes incoming event and the
// current aggregate state for that event into account, returns new state and an error if any.
type Processor interface {
	Process(*SecurityNotification, Aggregation) (Aggregation, error)
}

type PublishingProcessor struct{}

func (pp PublishingProcessor) Process(evt *SecurityNotification, existingState Aggregation) (Aggregation, error) {
	notification, newState := Strategy(evt, existingState)
	if notification == nil {
		return newState, nil
	}

	fmt.Printf("publishing new event for %s priority %s", evt.Email, evt.Priority)
	return nil, nil
}
