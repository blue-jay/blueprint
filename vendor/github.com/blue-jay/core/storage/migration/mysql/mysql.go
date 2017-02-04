// Package mysql implements MySQL migrations.
package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/blue-jay/core/storage"
	driver "github.com/blue-jay/core/storage/driver/mysql"
	"github.com/blue-jay/core/storage/migration"
	"github.com/jmoiron/sqlx"
)

// *****************************************************************************
// Thread-Safe Configuration
// *****************************************************************************
/*
var (
	info      database.Info
	infoMutex sync.RWMutex
)

// SetConfig stores the config.
func SetConfig(i database.Info) {
	infoMutex.Lock()
	info = i
	infoMutex.Unlock()
}

// ResetConfig removes the config.
func ResetConfig() {
	infoMutex.Lock()
	info = database.Info{}
	infoMutex.Unlock()
}

// Config returns the config.
func Config() database.Info {
	infoMutex.RLock()
	defer infoMutex.RUnlock()
	return info
}*/

// Configuration defines the shared configuration interface.
type Configuration struct {
	driver.Info
}

/*
// Shared returns the global configuration information.
func Shared() Configuration {
	return Configuration{
		Config(),
	}
}*/

// *****************************************************************************
// Migration Creation
// *****************************************************************************

// New creates a migration connection to the database.
func (c Configuration) New() (*migration.Info, error) {
	var mig *migration.Info

	// Load the config
	i := c.Info

	// Update the config
	i.Parameter = "parseTime=true&multiStatements=true"

	// Create MySQL entity
	mi := &Entity{}

	// Connect to the database
	con, err := i.Connect(true)

	// If the database doesn't exist or can't connect
	if err != nil {
		// Close the open connection (since 'unknown database' is still an
		// active connection)
		con.Close()

		// Connect to database without a database
		con, err = i.Connect(false)
		if err != nil {
			return mig, err
		}

		// Create the database
		err = i.Create(con)
		if err != nil {
			return mig, err
		}

		// Close connection
		con.Close()

		// Reconnect to the database
		con, err = i.Connect(true)
		if err != nil {
			return mig, err
		}
	}

	// Store the connection in the entity
	mi.sql = con

	// Store the migration table name
	mi.table = c.Migration.Table

	if len(mi.table) == 0 {
		return mig, errors.New("MySQL.Migration.Table key is missing in config file.")
	}

	// Setup logic was here
	return migration.New(mi, mi.table, c.Migration.Folder)
}

// *****************************************************************************
// Interface
// *****************************************************************************

// Item defines the migration table.
type Item struct {
	ID        uint32    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

// Entity defines fulfills the migration interface.
type Entity struct {
	table string
	sql   *sqlx.DB
}

// Extension returns the file extension with a period
func (t *Entity) Extension() string {
	return ".sql"
}

// TableExist returns true if the migration table exists
func (t *Entity) TableExist() error {
	_, err := t.sql.Exec(fmt.Sprintf("SELECT 1 FROM %v LIMIT 1;", t.table))
	if err != nil {
		return err
	}

	return err
}

// CreateTable returns true if the migration was created
func (t *Entity) CreateTable() error {
	_, err := t.sql.Exec(fmt.Sprintf(`CREATE TABLE %v (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  		name VARCHAR(191) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY (name),
  		PRIMARY KEY (id)
		);`, t.table))

	if err != nil {
		return err
	}

	return err
}

// Status returns last migration name
func (t *Entity) Status() (string, error) {
	result := &Item{}
	err := t.sql.Get(result, fmt.Sprintf("SELECT * FROM %v ORDER BY id DESC LIMIT 1;", t.table))

	// If no rows, then set to nil
	if err == sql.ErrNoRows {
		err = nil
	}

	return result.Name, err
}

// statusID returns last migration ID
func (t *Entity) statusID() (uint32, error) {
	result := &Item{}
	err := t.sql.Get(result, fmt.Sprintf("SELECT * FROM %v ORDER BY id DESC LIMIT 1;", t.table))
	return result.ID, err
}

// Migrate runs a query and returns error
func (t *Entity) Migrate(qry string) error {
	_, err := t.sql.Exec(qry)
	return err
}

// RecordUp adds a record to the database
func (t *Entity) RecordUp(name string) error {
	_, err := t.sql.Exec(fmt.Sprintf("INSERT INTO %v (name) VALUES (?);", t.table), name)
	return err
}

// RecordDown removes a record from the database and updates the AUTO_INCREMENT value
func (t *Entity) RecordDown(name string) error {
	_, err := t.sql.Exec(fmt.Sprintf("DELETE FROM %v WHERE name = ? LIMIT 1;", t.table), name)

	// If the record was removed successfully
	if err == nil {
		var ID uint32
		var nextID uint32 = 1

		// Get the last migration record now
		ID, err = t.statusID()

		// If there are no more migrations in the table
		if err == sql.ErrNoRows {
			// Leave ID at 1
		} else if err != nil {
			return err
		} else {
			nextID = ID
		}

		_, err = t.sql.Exec(fmt.Sprintf("ALTER TABLE %v AUTO_INCREMENT = %v;", t.table, nextID))
	}
	return err
}

// *****************************************************************************
// Test Helpers
// *****************************************************************************

// SetUp is a function for unit tests on a separate database.
func SetUp(envPath string, dbName string) (*migration.Info, Configuration) {
	// Get the environment variable
	if len(os.Getenv("JAYCONFIG")) == 0 {
		// Attempt to find env.json
		p, err := filepath.Abs(envPath)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Set the environment variable
		os.Setenv("JAYCONFIG", p)
	}

	// Get the config file path
	configFile := os.Getenv("JAYCONFIG")

	// Load the config
	config, err := storage.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Set the database name
	config.MySQL.Database = dbName

	// Set the migration folder to the absolute path
	config.MySQL.Migration.Folder = filepath.Join(filepath.Dir(configFile), config.MySQL.Migration.Folder)

	log.Println(config.MySQL.Migration.Folder)

	// Create the migration configuration
	conf := Configuration{
		config.MySQL,
	}

	// Create the migration
	mig, err := conf.New()
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Refresh the data
	mig.DownAll()
	mig.UpAll()

	return mig, conf
}

// TearDown removes the unit test database.
func TearDown(db *sqlx.DB, dbName string) error {
	// Drop the database
	_, err := db.Exec(fmt.Sprintf(`DROP DATABASE %v;`, dbName))
	return err
}
