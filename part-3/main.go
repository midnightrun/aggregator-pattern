package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dgraph-io/badger"
	"github.com/midnightrun/aggregator-pattern/part-3/aggregator"
)

var store aggregator.AggregationStore
var processor aggregator.PublishingProcessor

func main() {
	options := badger.DefaultOptions("./tmp")
	options.Logger = nil
	db, err := badger.Open(options)
	if err != nil {
		fmt.Printf("terminated service due to %v", err)
		return
	}

	defer db.Close()

	store = aggregator.NewStore(db)
	processor = aggregator.PublishingProcessor{}

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
	fmt.Printf("start handling request\n")

	var sn aggregator.SecurityNotification

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&sn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = store.ProcessNotification(&sn, processor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
