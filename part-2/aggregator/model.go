package aggregator

import (
	"time"

	"github.com/dgraph-io/badger"
)

type Priority int

const (
	LOW = iota
	MEDIUM
	HIGH
)

func (p Priority) String() string {
	return [...]string{"low", "medium", "high"}[p]
}

type SecurityNotification struct {
	Email        string
	Notification string
	Timestamp    time.Time
	Priority     Priority
}

type Aggregation []*SecurityNotification

type AggregationNotification struct {
	Email         string
	Notifications Aggregation
}

type AggregationStore struct {
	db *badger.DB
}

func NewStore(db *badger.DB) AggregationStore {
	return AggregationStore{db: db}
}

func (a *AggregationStore) Get(id string) (Aggregation, error) {
	var sns Aggregation

	err := a.db.View(func(txn *badger.Txn) error {
		var err error
		sns, err = getOrNil(txn, id)
		return err
	})

	return sns, err
}

func (a *AggregationStore) Save(id string, state Aggregation) error {
	b, err := marshal(state)
	if err != nil {
		return err
	}

	return a.db.Update(func(txn *badger.Txn) error {
		return txn.Set(keyForId(defaultPrefix, id), b)
	})
}
