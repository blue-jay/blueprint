// Package mysql implements MySQL migrations.
package mysql

import (
	"fmt"
	"time"

	"github.com/blue-jay/blueprint/lib/database"
	"github.com/blue-jay/blueprint/model"
)

var (
	// Table name
	Table     = "migration"
	extension = "sql"
)

// Entity defines the migration table
type Entity struct {
	ID        uint32    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

// Extension returns the file extension with a period
func (t *Entity) Extension() string {
	return "." + extension
}

// UpdateConfig will update any parameters necessary
func (t *Entity) UpdateConfig(config *database.Info) {
	config.MySQL.Parameter = "?parseTime=true&multiStatements=true"
}

// TableExist returns true if the migration table exists
func (t *Entity) TableExist() error {
	_, err := database.SQL.Exec(fmt.Sprintf("SELECT 1 FROM %v LIMIT 1;", Table))
	if err != nil {
		return err
	}

	return err
}

// CreateTable returns true if the migration was created
func (t *Entity) CreateTable() error {
	_, err := database.SQL.Exec(fmt.Sprintf(`CREATE TABLE %v (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY (name),
  		PRIMARY KEY (id)
		);`, Table))
	if err != nil {
		return err
	}

	return err
}

// Status returns last migration name
func (t *Entity) Status() (string, error) {
	result := &Entity{}
	err := database.SQL.Get(result, fmt.Sprintf("SELECT * FROM %v ORDER BY id DESC LIMIT 1;", Table))
	return result.Name, model.StandardError(err)
}

// statusID returns last migration ID
func statusID() (uint32, error) {
	result := &Entity{}
	err := database.SQL.Get(result, fmt.Sprintf("SELECT * FROM %v ORDER BY id DESC LIMIT 1;", Table))
	return result.ID, model.StandardError(err)
}

// Migrate runs a query and returns error
func (t *Entity) Migrate(qry string) error {
	_, err := database.SQL.Exec(qry)
	return err
}

// RecordUp adds a record to the database
func (t *Entity) RecordUp(name string) error {
	_, err := database.SQL.Exec(fmt.Sprintf("INSERT INTO %v (name) VALUES (?);", Table), name)
	return err
}

// RecordDown removes a record from the database and updates the AUTO_INCREMENT value
func (t *Entity) RecordDown(name string) error {
	_, err := database.SQL.Exec(fmt.Sprintf("DELETE FROM %v WHERE name = ? LIMIT 1;", Table), name)

	// If the record was removed successfully
	if err == nil {
		var ID uint32
		var nextID uint32 = 1

		// Get the last migration record now
		ID, err = statusID()

		// If there are no more migrations in the table
		if err == model.ErrNoResult {
			// Leave ID at 1
		} else if err != nil {
			return err
		} else {
			nextID = ID
		}

		_, err = database.SQL.Exec(fmt.Sprintf("ALTER TABLE %v AUTO_INCREMENT = %v;", Table, nextID))
	}
	return err
}
