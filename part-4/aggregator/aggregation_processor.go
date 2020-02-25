package aggregator

type AggregationProcessor interface {
	Process(Aggregation) error
}
