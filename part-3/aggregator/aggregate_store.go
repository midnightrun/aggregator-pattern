package aggregator

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger"
)

type AggregateStore struct {
	db *badger.DB
}

func (a *AggregateStore) ProcessNotification(n *SecurityNotification, p Processor) error {
	return a.db.Update(func(txn *badger.Txn) error {
		correlationId := n.Email

		previousState, err := getOrNil(correlationId)
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

func getOrNil(txn *badger.Txn, correlationId string) (Aggregation, error) {
	item, err := txn.Get(keyForId(defaultPrefix, correlationId))
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}

	var sns Aggregation
	err = item.Value(func(val []byte) error {
		sns, err = unmarshal(val)
		return err
	})

	return sns, err
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
