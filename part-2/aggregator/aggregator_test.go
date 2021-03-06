package aggregator

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
)

func createStore(t *testing.T) (*AggregationStore, func() error) {
	t.Helper()
	db, err := badger.Open(badger.DefaultOptions("./tmp"))
	if err != nil {
		log.Fatalf("could not open database: %v", err)
		return nil, nil
	}
	store := NewStore(db)
	return &store, db.Close
}

func dropAll() {
	db, err := badger.OpenManaged(badger.DefaultOptions("./tmp"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.DropAll()
}

func makeNotification(t *testing.T, email string) *SecurityNotification {
	t.Helper()
	return &SecurityNotification{
		Email:        email,
		Notification: "testing",
		Timestamp:    time.Now().UTC(),
		Priority:     LOW,
	}
}

func fatalIfError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error %v\n", err)
	}
}

func TestAggregation(t *testing.T) {
	n, s := Strategy(makeNotification(t, "testEmail"), make(Aggregation, 0))
	if n != nil {
		t.Fatalf("unexpected received a notification: got %#v", n)
	}

	n, s = Strategy(makeNotification(t, "testEmail"), s)
	if n != nil {
		t.Fatalf("unexpected received a notification: got %#v", n)
	}

	n, s = Strategy(makeNotification(t, "testEmail"), s)
	if n == nil {
		t.Fatal("expected a notification, got nil")
	}

	if n.Email != "testEmail" {
		t.Fatalf("incorrect aggregate email: got %s, wanted %s", n.Email, "testEmail")
	}

	if len(n.Notifications) != 3 {
		t.Fatalf("expected 3 messages in the aggregation, got %d", 1)
	}
}

func TestAggregationPublishesOnHighPriorityEvent(t *testing.T) {
	n, s := Strategy(makeNotification(t, "testEmail"), make(Aggregation, 0))
	if n != nil {
		t.Fatalf("unexpected received a notification: got %#v", n)
	}

	evt := makeNotification(t, "testEmail")
	evt.Priority = HIGH

	n, s = Strategy(evt, s)
	if n == nil {
		t.Fatalf("expected a notification, got nil")
	}

	if l := len(n.Notifications); l != 2 {
		t.Fatalf("expected 2 messages in the aggregation, got %d", l)
	}
}

func TestGetWithUnkownID(t *testing.T) {
	store, cleanup := createStore(t)
	defer cleanup()

	a, err := store.Get("testEmail")
	fatalIfError(t, err)

	if a != nil {
		t.Fatalf("unknown ID wanted nil, got %#v", a)
	}
}

func TestSave(t *testing.T) {
	store, cleanup := createStore(t)
	defer cleanup()

	notification := Aggregation{makeNotification(t, "testEmail")}

	err := store.Save("testEmail", notification)
	fatalIfError(t, err)

	loaded, err := store.Get("testEmail")
	fatalIfError(t, err)

	if !reflect.DeepEqual(notification, loaded) {
		t.Fatalf("save failed to save: wanted %#v, got %#v", notification, loaded)
	}
}
