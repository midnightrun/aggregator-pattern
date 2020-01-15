package aggregator

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Starting Aggegator Tests ...")
	exitVal := m.Run()
	log.Println("Shutting down Aggegator Tests")
	dropAll()
	os.Exit(exitVal)
}
