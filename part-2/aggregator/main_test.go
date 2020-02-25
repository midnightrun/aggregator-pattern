package aggregator

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Starting Aggregator Tests ...")
	exitVal := m.Run()
	log.Println("Shutting down Aggregator Tests")
	dropAll()
	os.Exit(exitVal)
}
