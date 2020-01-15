package aggregator

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger"
)

var defaultPrefix string = "aggregator_"

func Strategy(evt *SecurityNotification, s Aggregation) (*AggregationNotification, Aggregation) {
	s = append(s, evt)

	if len(s) == 3 || evt.Priority == HIGH {
		return &AggregationNotification{
			Email:         evt.Email,
			Notifications: s,
		}, nil
	}

	return nil, s
}

func getOrNil(txn *badger.Txn, id string) (Aggregation, error) {
	item, err := txn.Get(keyForId(defaultPrefix, id))
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

func keyForId(prefix string, id string) []byte {
	return []byte(fmt.Sprintf("%s%s", prefix, id))
}

func marshal(state Aggregation) ([]byte, error) {
	return json.Marshal(state)
}

func unmarshal(input []byte) (Aggregation, error) {
	var agg Aggregation
	err := json.Unmarshal(input, &agg)
	return agg, err
}
