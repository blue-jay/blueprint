// Package userstatus provides access to the user_status table in the MySQL database.
package userstatus

import (
	"database/sql"
	"time"
)

var (
	// table is the table name.
	table = "user_status"
)

// Item defines the model
type Item struct {
	ID        uint8     `db:"id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

// Connection is an interface for making queries.
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}
