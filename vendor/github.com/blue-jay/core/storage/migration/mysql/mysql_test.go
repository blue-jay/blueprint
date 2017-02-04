// Package mysql_test tests the MySQL migration process.
package mysql_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/blue-jay/core/storage"
	"github.com/blue-jay/core/storage/migration"
	"github.com/blue-jay/core/storage/migration/mysql"

	"github.com/jmoiron/sqlx"
)

var (
	migrationFolder = "testdata/migration_files"
	conf            mysql.Configuration
	con             Connection
)

// TestMain runs setup, tests, and then teardown.
func TestMain(m *testing.M) {
	// For use with coveralls
	file := "envtest.json"

	_, conf = mysql.SetUp("testdata/"+file, "database_test")

	// Connect to the database
	db, _ := conf.Connect(true)
	con = Connection{
		db: db,
	}

	returnCode := m.Run()
	teardown()
	os.Exit(returnCode)
}

// loadConfig will read the config from the env.json file
func loadConfig() mysql.Configuration {
	info, err := storage.LoadConfig(os.Getenv("JAYCONFIG"))
	if err != nil {
		log.Fatalf("%v", err)
	}
	info.MySQL.Database = "database_test"
	info.MySQL.Migration.Folder = migrationFolder

	// Connect to the database
	return mysql.Configuration{
		info.MySQL,
	}
}

// setup handles any start up tasks.
func setup() *migration.Info {
	mig, err := conf.New()
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Remove table
	con.deleteTable("test_brother")

	// Remove the folder
	err = os.RemoveAll(migrationFolder)
	if err != nil {
		log.Fatal(err)
	}

	// Make the folder
	err = os.MkdirAll(migrationFolder, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return mig
}

// teardown handles any clean up tasks.
func teardown() {
	// Remove the folder
	err := os.RemoveAll(migrationFolder)
	if err != nil {
		log.Fatal(err)
	}

	// Remove the database
	mysql.TearDown(con.db, "database_test")
}

// TestCreateTable.
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

// TestDropTable.
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

// TestInsertRows.
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
	result, _ := con.byID("1")
	if result.Name != "Joey" {
		t.Errorf("record retrieved is incorrect: '%v'", result.Name)
	}
}

// TestDeleteRows.
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
	result, err := con.byID("1")
	if result.Name != "Joey" {
		t.Errorf("record retrieved is incorrect: '%v'", result.Name)
	}

	// Run the migration
	err = mig.DownAll()
	if err != nil {
		t.Errorf("could not migrate down: %v", err)
	}

	// Test querying the data
	result, _ = con.byID("1")
	if result.Name == "Joey" {
		t.Errorf("record retrieved is incorrect: '%v'", result.Name)
	}
}

// TestAlterRows.
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
	mig.UpOne()

	// Insert rows migration
	setupMigrateInsert(mig)
	err = mig.Create("Insert brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Run the migration
	mig.UpOne()

	// Test querying the data
	result, err := con.byID("1")
	if result.Name != "Joey" {
		t.Errorf("record retrieved is incorrect: '%v'", result.Name)
	}

	// Run the migration
	mig.DownOne()

	// Test querying the data
	result, err = con.byID("1")
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
	mig.UpAll()

	// Test querying the data
	result, err = con.byID("1")
	if result.Age != 0 {
		t.Errorf("record retrieved is incorrect: '%v'", result.Age)
	}

	// Run the migration
	mig.DownAll()

	// Update column migration
	setupMigrateUpdate(mig)
	err = mig.Create("Update brother table")
	if err != nil {
		t.Errorf("could not create migration: %v", err)
	}

	// Run the migration
	mig.UpAll()

	// Test querying the data
	result, _ = con.byID("1")
	if result.Age != 28 {
		t.Errorf("record retrieved is incorrect: '%v'", result.Age)
	}

	// Run the migration
	mig.DownOne()

	// Test querying the data
	result, _ = con.byID("1")
	if result.Age == 28 {
		t.Errorf("record retrieved is incorrect: '%v'", result.Age)
	}

	// Run the migration
	mig.UpOne()

	// Test querying the data
	result, _ = con.byID("1")
	if result.Age != 28 {
		t.Errorf("record retrieved is incorrect: '%v'", result.Age)
	}
}

// *****************************************************************************
// Models
// *****************************************************************************

// Entity defines the brother table.
type Entity struct {
	ID   uint32 `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

// Connection defines the shared database interface.
type Connection struct {
	db *sqlx.DB
}

// byID gets note by ID.
func (c Connection) byID(ID string) (Entity, error) {
	result := Entity{}
	err := c.db.Get(&result, "SELECT * FROM test_brother WHERE id = ? LIMIT 1", ID)
	return result, err
}

// deleteTable drops a table.
func (c Connection) deleteTable(table string) (sql.Result, error) {
	result, err := c.db.Exec(fmt.Sprintf("DROP TABLE %v", table))
	return result, err
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
