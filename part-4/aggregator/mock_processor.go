package aggregator

type mockProcessor struct {
	processedNotification *SecurityNotification
	processedAggregation  Aggregation
	returnAggregation     Aggregation
	err                   error
}

func (p *mockProcessor) Process(n *SecurityNotification, a Aggregation) (Aggregation, error) {
	p.processedNotification = n
	p.processedAggregation = a
	return p.returnAggregation, p.err
}
