package note_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/blue-jay/blueprint/model/note"
	"github.com/blue-jay/blueprint/model/user"

	"github.com/blue-jay/core/storage/migration"
	"github.com/blue-jay/core/storage/migration/mysql"
)

var (
	mig *migration.Info
)

// TestMain runs setup, tests, and then teardown.
func TestMain(m *testing.M) {
	setup()
	returnCode := m.Run()
	teardown()
	os.Exit(returnCode)
}

// setup handles any start up tasks.
func setup() {
	mysql.SetUp("../../env.json", "database_test")
}

// teardown handles any clean up tasks.
func teardown() {
	mysql.TearDown()
}

// TestComplete
func TestComplete(t *testing.T) {
	data := "Test data."
	dataNew := "New test data."

	result, err := user.Create("John", "Doe", "jdoe@domain.com", "p@$$W0rD")
	if err != nil {
		t.Error("could not create user:", err)
	}

	uID, err := result.LastInsertId()
	if err != nil {
		t.Error("could not convert user ID:", err)
	}

	// Convert ID to string
	userID := fmt.Sprintf("%v", uID)

	// Create a record
	result, err = note.Create(data, userID)
	if err != nil {
		t.Error("could not create record:", err)
	}

	// Get the last ID
	ID, err := result.LastInsertId()
	if err != nil {
		t.Error("could not convert ID:", err)
	}

	// Convert ID to string
	lastID := fmt.Sprintf("%v", ID)

	// Select a record
	record, err := note.ByID(lastID, userID)
	if err != nil {
		t.Error("could not retrieve record:", err)
	} else if record.Name != data {
		t.Errorf("retrieved wrong record: got '%v' want '%v'", record.Name, data)
	}

	// Update a record
	result, err = note.Update(dataNew, lastID, userID)
	if err != nil {
		t.Error("could not update record:", err)
	}

	// Select a record
	record, err = note.ByID(lastID, userID)
	if err != nil {
		t.Error("could not retrieve record:", err)
	} else if record.Name != dataNew {
		t.Errorf("retieved wrong record: got '%v' want '%v'", record.Name, dataNew)
	}

	// Delete a record by ID
	result, err = note.Delete(lastID, userID)
	if err != nil {
		t.Error("could not delete record:", err)
	}

	// Count the number of deleted rows
	rows, err := result.RowsAffected()
	if err != nil {
		t.Error("could not count affected rows:", err)
	} else if rows != 1 {
		t.Error("incorrect number of affected rows:", rows)
	}
}
