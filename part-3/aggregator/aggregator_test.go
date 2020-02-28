package aggregator

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
)

func createBadgerStore(t *testing.T) (*AggregationStore, func() error) {
	t.Helper()
	db, err := badger.Open(badger.DefaultOptions("./tmp"))
	if err != nil {
		t.Fatalf("could not open database: %v", err)
		return nil, nil
	}
	store := NewStore(db)
	return &store, db.Close
}

func dropAll() {
	db, err := badger.OpenManaged(badger.DefaultOptions("./tmp"))
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.DropAll()
}

func makeNotification(correlationId string) *SecurityNotification {
	return &SecurityNotification{
		Email:        correlationId,
		Priority:     0,
		Timestamp:    time.Now().UTC(),
		Notification: "test notification",
	}
}

func fatalIfError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestProcessNotificationWithNoPreviousState(t *testing.T) {
	store, cleanup := createBadgerStore(t)
	defer cleanup()

	processor := &mockProcessor{}
	notification := makeNotification("testEmail")

	err := store.ProcessNotification(notification, processor)
	fatalIfError(t, err)

	if processor.processedNotification == nil {
		t.Fatalf("processor did not receive notification\n")
	}
}

func TestProcessNotificationErrorOnPublish(t *testing.T) {
	store, cleanup := createBadgerStore(t)
	defer cleanup()

	publishingError := errors.New("error on publishing")

	processor := &mockProcessor{}
	processor.err = publishingError

	notification := makeNotification("testEmail")

	err := store.ProcessNotification(notification, processor)
	if err != publishingError {
		t.Fatalf("expected publishing error but got %v\n", err)
	}

	loaded, err := store.Get("testEmail")
	if err != nil {
		t.Fatalf("error while getting aggregation")
	}

	if loaded != nil {
		t.Fatalf("expected nil but got %v\n", loaded)
	}
}
