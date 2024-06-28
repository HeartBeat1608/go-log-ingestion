package internals_test

import (
	"log"
	"os"
	"testing"

	"github.com/HeartBeat1608/logmaster/internals"
)

var cm *internals.DBConnectionManager

func TestMain(m *testing.M) {
	// setup
	internals.LoadConfig("../config.json")
	cm = internals.NewConnectionManager()

	// main code
	code := m.Run()

	// teardown

	// exit
	os.Exit(code)
}

func TestAddService(t *testing.T) {
	serviceName := "test"

	_, err := cm.AddService(serviceName)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetOrAddConnection(t *testing.T) {
	serviceName := "test"

	t.Skip("Skipped")

	db, err := cm.GetOrAddConnection(serviceName)
	if err != nil {
		t.Fatalf("%v", err)
	}

	rows, err := db.Queryx(".tables")
	if err != nil {
		t.Fatalf("%v", err)
	}

	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			log.Printf("%v\n", rows)
		}
	}
}
