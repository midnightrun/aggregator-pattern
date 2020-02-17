package aggregator

// Processor handles incoming notifications and previous state.
//
// Executed within the context of a transaction, takes incoming event and the
// current aggregate state for that event into account, returns new state and an error if any.
type Processor interface {
	Process(*SecurityNotification, Aggregation) (Aggregation, error)
}
