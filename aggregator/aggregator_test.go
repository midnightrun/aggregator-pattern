package aggregator

import (
	"testing"
	"time"
)

func makeNotification(email string) *SecurityNotification {
	return &SecurityNotification{
		Email:        email,
		Notification: "testing",
		Timestamp:    time.Now().UTC(),
		Priority:     LOW,
	}
}

func TestAggregation(t *testing.T) {
	n, s := Strategy(makeNotification("testEmail"), make(Aggregation, 0))
	if n != nil {
		t.Fatalf("unexpected received a notification: got %#v", n)
	}

	n, s = Strategy(makeNotification("testEmail"), s)
	if n != nil {
		t.Fatalf("unexpected received a notification: got %#v", n)
	}

	n, s = Strategy(makeNotification("testEmail"), s)
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
	n, s := Strategy(makeNotification("testEmail"), make(Aggregation, 0))
	if n != nil {
		t.Fatalf("unexpected received a notification: got %#v", n)
	}

	evt := makeNotification("testEmail")
	evt.Priority = HIGH

	n, s = Strategy(evt, s)
	if n == nil {
		t.Fatalf("expected a notification, got nil")
	}

	if l := len(n.Notifications); l != 2 {
		t.Fatalf("expected 2 messages in the aggregation, got %d", l)
	}
}
