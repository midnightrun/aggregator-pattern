package aggregator

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
)

func createBadgerStore(t *testing.T) (*AggregationStore, func() error) {
	t.Helper()

	options := badger.DefaultOptions("./tmp")
	options.Logger = nil

	db, err := badger.Open(options)
	if err != nil {
		t.Fatalf("could not open database: %v", err)
		return nil, nil
	}
	store := NewStore(db)
	return &store, db.Close
}

func dropAll() {
	options := badger.DefaultOptions("./tmp")
	options.Logger = nil

	db, err := badger.OpenManaged(options)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.DropAll()
}

func makeNotification(correlationId string) SecurityNotification {
	return SecurityNotification{
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

	if !reflect.DeepEqual(notification, processor.processedNotification) {
		t.Fatalf("processor did not receive correct notification")
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

	loaded, err := store.Get(notification.Email)
	fatalIfError(t, err)

	if len(loaded.Notifications) > 0 {
		t.Fatalf("expected nil but got %v", loaded)
	}
}

func TestProcessNotificationOnPublish(t *testing.T) {
	store, cleanup := createBadgerStore(t)
	defer cleanup()

	email := "testEmail"
	state := Aggregation{
		Email:      email,
		LastUpdate: time.Now().UTC(),
		Notifications: []SecurityNotification{
			{
				Email:        email,
				Notification: "test mail",
			},
			{
				Email:        email,
				Notification: "test mail",
			},
		},
	}

	err := store.Save(state, email)
	if err != nil {
		t.Fatal("error while saving state to store")
	}

	processor := &mockProcessor{}
	processor.returnAggregation = Aggregation{}

	notification := SecurityNotification{
		Email:        email,
		Notification: "notification",
	}

	err = store.ProcessNotification(notification, processor)
	if err != nil {
		t.Fatalf("error while processing notification")
	}

	loaded, err := store.Get(email)
	if err != nil {
		t.Fatalf("error while getting aggregation")
	}

	if len(loaded.Notifications) > 0 {
		t.Fatalf("expected nil but got %v", loaded)
	}
}

// Todo Deep Equal test missing

func TestProcessAggregationAfterTreshold(t *testing.T) {
	store, cleanup := createBadgerStore(t)
	defer cleanup()

	email := "testEmail"

	aggregation := Aggregation{
		Email:      email,
		LastUpdate: time.Now().Add(-3 * time.Hour),
		Notifications: []SecurityNotification{
			{
				Timestamp:    time.Now().UTC(),
				Priority:     0,
				Email:        email,
				Notification: "testMail",
			},
			{
				Timestamp:    time.Now().UTC(),
				Priority:     0,
				Email:        email,
				Notification: "testMail",
			},
			{
				Timestamp:    time.Now().UTC(),
				Priority:     0,
				Email:        email,
				Notification: "testMail",
			},
		},
	}

	err := store.Save(aggregation, email)
	if err != nil {
		t.Fatalf("saving aggregation to database failed due to %v\n", err)
	}

	// Todo: Add config for setting the time
	processor := &mockProcessor{}
	err = store.ProcessAggregations(processor)

	loaded, err := store.Get(email)

	fatalIfError(t, err)

	if len(loaded.Notifications) > 0 {
		t.Fatalf("expected no notifications but got %v", loaded)
	}
}
