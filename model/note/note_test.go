package note_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/blue-jay/blueprint/config"
	"github.com/blue-jay/blueprint/lib/database"
	"github.com/blue-jay/blueprint/lib/migration"
	"github.com/blue-jay/blueprint/lib/migration/mysql"
	"github.com/blue-jay/blueprint/model/note"
)

var (
	mig *migration.Info
)

// TestMain runs setup, tests, and then teardown
func TestMain(m *testing.M) {
	setup()
	returnCode := m.Run()
	teardown()
	os.Exit(returnCode)
}

// setup handles any start up tasks
func setup() {
	var err error

	// Change the working directory to the root
	os.Chdir("../../")

	// Load the configuration
	info := config.Load()

	// Create MySQL entity
	db := &mysql.Entity{}

	// Update the config
	db.UpdateConfig(&info.Database)

	// Connect to database
	database.Connect(info.Database)

	// Create a new migration
	mig, err = migration.New(db, filepath.Join("database", "migration"))
	if err != nil {
		log.Fatal(err)
	}

	// Run the up migrations
	err = mig.UpAll()
	if err != nil {
		log.Fatal(err)
	}
}

// teardown handles any cleanup teasks
func teardown() {
	err := mig.DownAll()
	if err != nil {
		log.Fatal(err)
	}
}

// TestCrud
func TestCrud(t *testing.T) {
	userID := "1"
	data := "This is test data."
	dataNew := "This is new test data."

	// Create a record
	result, err := note.Create(data, userID)
	if err != nil {
		t.Errorf("could not create record:", err)
	}

	// Get the last ID
	ID, err := result.LastInsertId()
	if err != nil {
		t.Errorf("could not convert ID:", err)
	}

	// Convert ID to string
	lastID := fmt.Sprintf("%v", ID)

	// Select a record
	record, err := note.ByID(lastID, userID)
	if err != nil {
		t.Errorf("could not retrieve record:", err)
	} else if record.Content != data {
		t.Errorf("retieved wrong record: got '%v' want '%v'", record.Content, data)
	}

	// Update a record
	result, err = note.Update(dataNew, lastID, userID)
	if err != nil {
		t.Errorf("could not update record:", err)
	}

	// Select a record
	record, err = note.ByID(lastID, userID)
	if err != nil {
		t.Errorf("could not retrieve record:", err)
	} else if record.Content != dataNew {
		t.Errorf("retieved wrong record: got '%v' want '%v'", record.Content, dataNew)
	}

	// Delete a record by ID
	result, err = note.Delete(lastID, userID)
	if err != nil {
		t.Errorf("could not delete record:", err)
	}

}
