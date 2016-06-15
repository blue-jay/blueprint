package migration_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/blue-jay/blueprint/config"
	"github.com/blue-jay/blueprint/lib/database"
	"github.com/blue-jay/blueprint/lib/migration"
	"github.com/blue-jay/blueprint/lib/migration/mysql"
	"github.com/blue-jay/blueprint/model"
)

var (
	migrationFolder = filepath.Join("database", "migration_test")
)

// TestMain runs setup, tests, and then teardown
func TestMain(m *testing.M) {
	globalSetup()
	returnCode := m.Run()
	teardown()
	os.Exit(returnCode)
}

// globalSetup handles any start up tasks
func globalSetup() {
	// Change the working directory to the root
	os.Chdir("../../")

	// Load the configuration
	info := config.Load()

	// Set the table name
	mysql.Table = "test_migration"

	// Create MySQL entity
	db := &mysql.Entity{}

	// Update the config
	db.UpdateConfig(&info.Database)

	// Connect to database
	database.Connect(info.Database)
}

func setup() *migration.Info {
	// Delete anything from previous test funs
	teardown()

	// Make the folder
	err := os.MkdirAll(migrationFolder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Create MySQL entity
	db := &mysql.Entity{}

	// Create a new migration
	mig, err := migration.New(db, migrationFolder)
	if err != nil {
		log.Fatal(err)
	}

	return mig
}

// teardown handles any cleanup teasks
func teardown() {
	// Remove the folder
	err := os.RemoveAll(migrationFolder)
	if err != nil {
		log.Fatal(err)
	}

	_, err = deleteTable("test_brother")
	_, err = deleteTable("test_migration")

}

// TestCreateTable
func TestCreateTable(t *testing.T) {
	var err error
	mig := setup()

	// Create table migration
	setupMigrateCreate(mig)
	err = mig.Create("Create brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Run the migration
	err = mig.UpOne()
	if err != nil {
		t.Errorf("could not migrate up: %v", err)
	}
}

// TestDrppTable
func TestDropTable(t *testing.T) {
	var err error
	mig := setup()

	// Create table migration
	setupMigrateCreate(mig)
	err = mig.Create("Create brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Run the migration
	err = mig.UpOne()
	if err != nil {
		t.Errorf("could not migrate up: %v", err)
	}

	// Run the migration
	err = mig.DownOne()
	if err != nil {
		t.Errorf("could not migrate down: %v", err)
	}
}

// TestInsertRows
func TestInsertRows(t *testing.T) {
	var err error
	mig := setup()

	// Create table migration
	setupMigrateCreate(mig)
	err = mig.Create("Create brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Insert rows migration
	setupMigrateInsert(mig)
	err = mig.Create("Insert brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Run the migration
	err = mig.UpAll()
	if err != nil {
		t.Errorf("could not migrate up: %v", err)
	}

	// Test querying the data
	result, err := byID("1")
	if result.Name != "Joey" {
		t.Errorf("record retrieved is incorrect: '%v'", result.Name)
	}
}

// TestDeleteRows
func TestDeleteRows(t *testing.T) {
	var err error
	mig := setup()

	// Create table migration
	setupMigrateCreate(mig)
	err = mig.Create("Create brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Insert rows migration
	setupMigrateInsert(mig)
	err = mig.Create("Insert brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Run the migration
	err = mig.UpAll()
	if err != nil {
		t.Errorf("could not migrate up: %v", err)
	}

	// Test querying the data
	result, err := byID("1")
	if result.Name != "Joey" {
		t.Errorf("record retrieved is incorrect: '%v'", result.Name)
	}

	// Run the migration
	err = mig.DownAll()
	if err != nil {
		t.Errorf("could not migrate down: %v", err)
	}

	// Test querying the data
	result, err = byID("1")
	if result.Name == "Joey" {
		t.Errorf("record retrieved is incorrect: '%v'", result.Name)
	}
}

// TestAlterRows
func TestAlterRows(t *testing.T) {
	var err error
	mig := setup()

	// Create table migration
	setupMigrateCreate(mig)
	err = mig.Create("Create brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Run the migration
	err = mig.UpOne()
	if err != nil {
		t.Errorf("could not migrate up: %v", err)
	}

	// Insert rows migration
	setupMigrateInsert(mig)
	err = mig.Create("Insert brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Run the migration
	err = mig.UpOne()
	if err != nil {
		t.Errorf("could not migrate up: %v", err)
	}

	// Test querying the data
	result, err := byID("1")
	if result.Name != "Joey" {
		t.Errorf("record retrieved is incorrect: '%v'", result.Name)
	}

	// Run the migration
	err = mig.DownOne()
	if err != nil {
		t.Errorf("could not migrate down: %v", err)
	}

	// Test querying the data
	result, err = byID("1")
	if result.Name == "Joey" {
		t.Errorf("record retrieved is incorrect: '%v'", result.Name)
	}

	// Alter column migration
	setupMigrateAlter(mig)
	err = mig.Create("Alter brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Run the migration
	err = mig.UpAll()
	if err != nil {
		t.Errorf("could not migrate up: %v", err)
	}

	// Test querying the data
	result, err = byID("1")
	if result.Age != 0 {
		t.Errorf("record retrieved is incorrect: '%v'", result.Age)
	}

	// Run the migration
	err = mig.DownAll()
	if err != nil {
		t.Errorf("could not migrate down: %v", err)
	}

	// Update column migration
	setupMigrateUpdate(mig)
	err = mig.Create("Update brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Run the migration
	err = mig.UpAll()
	if err != nil {
		t.Errorf("could not migrate up: %v", err)
	}

	// Test querying the data
	result, err = byID("1")
	if result.Age != 28 {
		t.Errorf("record retrieved is incorrect: '%v'", result.Age)
	}

	// Run the migration
	err = mig.DownOne()
	if err != nil {
		t.Errorf("could not migrate down: %v", err)
	}

	// Test querying the data
	result, err = byID("1")
	if result.Age == 28 {
		t.Errorf("record retrieved is incorrect: '%v'", result.Age)
	}

	// Run the migration
	err = mig.UpOne()
	if err != nil {
		t.Errorf("could not migrate up: %v", err)
	}

	// Test querying the data
	result, err = byID("1")
	if result.Age != 28 {
		t.Errorf("record retrieved is incorrect: '%v'", result.Age)
	}
}

// *****************************************************************************
// Models
// *****************************************************************************

// Entity defines the brother table
type Entity struct {
	ID   uint32 `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

// byID gets note by ID
func byID(ID string) (Entity, error) {
	result := Entity{}
	err := database.SQL.Get(&result, "SELECT * FROM test_brother WHERE id = ? LIMIT 1", ID)
	return result, model.StandardError(err)
}

// deleteTable drops a table
func deleteTable(table string) (sql.Result, error) {
	result, err := database.SQL.Exec(fmt.Sprintf("DROP TABLE %v", table))
	return result, model.StandardError(err)
}

// *****************************************************************************
// Test Migrations
// *****************************************************************************

func setupMigrateCreate(mig *migration.Info) {
	mig.TemplateUp = `
CREATE TABLE test_brother (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    PRIMARY KEY (id)
);
`

	mig.TemplateDown = `
DROP TABLE test_brother;
`
}

func setupMigrateInsert(mig *migration.Info) {
	mig.TemplateUp = `
INSERT INTO test_brother (name) VALUES ("Joey");
INSERT INTO test_brother (name) VALUES ("Jarrod");
INSERT INTO test_brother (name) VALUES ("Trent");
INSERT INTO test_brother (name) VALUES ("Troy");
`

	mig.TemplateDown = `
DELETE FROM test_brother;
ALTER TABLE test_brother AUTO_INCREMENT = 1;
`
}

func setupMigrateAlter(mig *migration.Info) {
	mig.TemplateUp = `
ALTER TABLE test_brother ADD COLUMN age int;
`

	mig.TemplateDown = `
ALTER TABLE test_brother DROP COLUMN age;
`
}

func setupMigrateUpdate(mig *migration.Info) {
	mig.TemplateUp = `
UPDATE test_brother SET age = 28 WHERE id = 1;
UPDATE test_brother SET age = 26 WHERE id = 2;
UPDATE test_brother SET age = 24 WHERE id = 3;
UPDATE test_brother SET age = 23 WHERE id = 4;
`

	mig.TemplateDown = `
UPDATE test_brother SET age = NULL;
`
}

// *****************************************************************************
// Helpers
// *****************************************************************************

// folderExists will exit if the folder doesn't exist
func folderExists(dir string) bool {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
