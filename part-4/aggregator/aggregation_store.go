package aggregator

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger"
)

var defaultPrefix string = "aggregator_"

type AggregationStore struct {
	db *badger.DB
}

func NewStore(db *badger.DB) AggregationStore {
	return AggregationStore{db: db}
}

func (a *AggregationStore) ProcessNotification(n *SecurityNotification, p Processor) error {
	fmt.Printf("get database entry for %s\n", n.Email)

	return a.db.Update(func(txn *badger.Txn) error {
		correlationId := n.Email

		previousState, err := getOrNil(txn, correlationId)
		if err != nil {
			return err
		}

		newState, err := p.Process(n, previousState)
		if err != nil {
			return err
		}

		b, err := marshal(newState)
		if err != nil {
			return err
		}

		return txn.Set(keyForId(defaultPrefix, correlationId), b)
	})
}

func (a *AggregationStore) Save(aggregation Aggregation, correlationId string) error {
	return a.db.Update(func(txn *badger.Txn) error {
		b, err := json.Marshal(aggregation)
		if err != nil {
			return err
		}

		return txn.Set(keyForId(defaultPrefix, correlationId), b)
	})
}

func (a *AggregationStore) Get(correlationId string) (*Aggregation, error) {
	var aggregation Aggregation

	err := a.db.View(func(txn *badger.Txn) error {
		var err error
		aggregation, err = getOrNil(txn, correlationId)
		return err
	})

	return aggregation, err
}

func getOrNil(txn *badger.Txn, correlationId string) (*Aggregation, error) {
	item, err := txn.Get(keyForId(defaultPrefix, correlationId))
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}

	var sns Aggregation
	err = item.Value(func(val []byte) error {
		sns, err = unmarshal(val)
		return err
	})

	return &sns, err
}

func keyForId(prefix string, correlationId string) []byte {
	return []byte(fmt.Sprintf("%s%s", prefix, correlationId))
}

func marshal(state Aggregation) ([]byte, error) {
	return json.Marshal(state)
}

func unmarshal(input []byte) (Aggregation, error) {
	var agg Aggregation
	err := json.Unmarshal(input, &agg)

	return agg, err
}
