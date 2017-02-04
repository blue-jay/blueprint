// Package postgresql provides a wrapper around the pq package.
package postgresql

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

// Info holds the details for the connection.
type Info struct {
	Username        string
	Password        string
	Database        string
	Hostname        string
	Port            int
	Parameter       string
	MigrationFolder string
	Extension       string
}

// *****************************************************************************
// Database Handling
// *****************************************************************************

// Connect to the database.
func (c Info) Connect(specificDatabase bool) (*sqlx.DB, error) {
	// Connect to database and ping
	return sqlx.Connect("postgres", c.dsn(specificDatabase))
}

// Create a new database.
func (c Info) Create(sql *sqlx.DB) error {
	// Create the database
	_, err := sql.Exec(fmt.Sprintf(`CREATE DATABASE %v;`, c.Database))

	return err
}

// Drop a database.
func (c Info) Drop(sql *sqlx.DB) error {
	// Drop the database
	_, err := sql.Exec(fmt.Sprintf(`DROP DATABASE %v;`, c.Database))

	return err
}

// *****************************************************************************
// Database Specific
// *****************************************************************************

// DSN returns the Data Source Name.
func (c Info) dsn(includeDatabase bool) string {
	// Set defaults
	ci := c.setDefaults()

	// Build parameters
	param := ci.Parameter

	// If parameter is specified, add a question mark
	// Don't add one if a question mark is already there
	if len(ci.Parameter) > 0 && !strings.HasPrefix(ci.Parameter, "?") {
		param = "?" + ci.Parameter

	}

	// Example: postgres://pqgotest:password@localhost/pqgotest
	s := fmt.Sprintf("postgres://%v:%v@%v:%d/%v", ci.Username, ci.Password, ci.Hostname, ci.Port, param)

	if includeDatabase {
		s = fmt.Sprintf("postgres://%v:%v@%v:%d/%v%v", ci.Username, ci.Password, ci.Hostname, ci.Port, ci.Database, param)
	}

	return s
}

// setDefaults sets the charset and collation if they are not set.
func (c Info) setDefaults() Info {
	ci := c

	return ci
}
