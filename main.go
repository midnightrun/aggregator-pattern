package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/midnightrun/aggregator-pattern/aggregator"
)

var aggregationStore map[string]aggregator.Aggregation

func main() {
	aggregationStore = make(map[string]aggregator.Aggregation, 0)

	http.HandleFunc("/notifications", aggregatorHandler)

	fmt.Printf("receiving on http://localhost:8080/notifications\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
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

	existingState, ok := aggregationStore[correlationID]
	if !ok {
		existingState = make(aggregator.Aggregation, 0)
	}

	var n *aggregator.AggregationNotification

	n, aggregationStore[correlationID] = aggregator.Strategy(&sn, existingState)
	if n != nil {
		log.Printf("new event emitted for user %s\n", n.Email)
		return
	}

	log.Println("event processed - no event emitted")
}
