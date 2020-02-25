package aggregator

type AggregationProcessor interface {
	Process(state Aggregation) (*Aggregation, error)
}
