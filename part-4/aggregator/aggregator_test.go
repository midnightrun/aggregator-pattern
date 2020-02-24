package aggregator

import (
	"errors"
	"testing"
	"time"
)

func TestProcessNotificationWithNoPreviousState(t *testing.T) {
	store, cleanup := createBadgerStore(t)
	defer cleanup()

	processor := &mockProcessor{}
	notification := makeNotification("testEmail")

	err := store.ProcessNotification(notification, processor)
	fatalIfError(t, err)

	if processor.processedNotification == nil {
		t.Fatalf("processor did not receive notification")
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
		t.Fatalf("expected publishing error but got %v", err)
	}

	loaded := store.Get(notification.id)

	if loaded != nil {
		t.Fatalf("expected nil but got %v", loaded)
	}
}

// Todo Deep Equal test missing

func TestProcessAggregationAfterTreshold(t *testing.T) {
	db, cleanup := createBadgerStore()
	defer cleanup()

	store := NewStore(db)

	aggregation := Aggregation{
		&SecurityNotification{
			Email:        "testEmail",
			Notification: "testing",
			Timestamp:    time.Now().H.UTC()},
	}

	err := store.Save(aggregation, "testEmail")
	if err != nil {
		t.Fatalf("saving aggregation to database failed due to %v\n", err)
	}

	processor := &mockAggregationProcessor{}
	err := store.ProcessAggregation(processor)

	loaded := store.Get("testEmail")

	if loaded != nil {
		t.Fatalf("expected nil but got %v", loaded)
	}
}
