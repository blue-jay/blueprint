// Package model provides access to the data in the MySQL database.
package model

import (
	"database/sql"
)

// Connection is an interface for making queries.
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

// Info is a collection of all the models.
type Info struct {
	Note
	User
}
