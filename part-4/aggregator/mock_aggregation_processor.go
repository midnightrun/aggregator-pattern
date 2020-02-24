package aggregator

import "fmt"

type mockAggregationProcessor struct{}

func (ap mockAggregationProcessor) Process(state Aggregation) (*Aggregation, error) {
	fmt.Println("Processing not implemented for type mockAggregationProcessor")

	return nil, nil
}
