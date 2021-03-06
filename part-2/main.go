package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dgraph-io/badger"
	"github.com/midnightrun/aggregator-pattern/part-2/aggregator"
)

var store aggregator.AggregationStore

func main() {
	db, err := badger.Open(badger.DefaultOptions("./tmp"))
	if err != nil {
		fmt.Printf("terminated service on http://localhost:8080/notifications due to %s\n", err)
		return
	}

	defer db.Close()

	store = aggregator.NewStore(db)

	http.HandleFunc("/notifications", aggregatorHandler)

	errs := make(chan error, 2)

	go func() {
		fmt.Printf("receiving on http://localhost:8080/notifications\n")
		errs <- http.ListenAndServe(":8080", nil)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("terminated service on http://localhost:8080/notifications due to %s\n", <-errs)
}

func aggregatorHandler(w http.ResponseWriter, r *http.Request) {
	var sn aggregator.SecurityNotification

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&sn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	correlationID := sn.Email
	log.Printf("processing %s event for %s\n", sn.Priority, correlationID)

	existingState, err := store.Get(correlationID)
	if err != nil {
		existingState = make(aggregator.Aggregation, 0)
	}

	var n *aggregator.AggregationNotification

	n, s := aggregator.Strategy(&sn, existingState)

	err = store.Save(sn.Email, s)
	if err != nil {
		log.Printf("failed save operation: %v\n", err)
	}

	if n != nil {
		log.Printf("new event emitted for user %s\n", n.Email)
		return
	}

	log.Println("event processed - no event emitted")
}
