// Package user provides access to the user table in the MySQL database.
package user

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var (
	// table is the table name.
	table = "user"
)

// Item defines the model.
type Item struct {
	ID        uint32         `db:"id"`
	FirstName string         `db:"first_name"`
	LastName  string         `db:"last_name"`
	Email     string         `db:"email"`
	Password  string         `db:"password"`
	StatusID  uint8          `db:"status_id"`
	CreatedAt mysql.NullTime `db:"created_at"`
	UpdatedAt mysql.NullTime `db:"updated_at"`
	DeletedAt mysql.NullTime `db:"deleted_at"`
}

// Connection is an interface for making queries.
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

// ByEmail gets user information from email.
func ByEmail(db Connection, email string) (Item, bool, error) {
	result := Item{}
	err := db.Get(&result, fmt.Sprintf(`
		SELECT id, password, status_id, first_name
		FROM %v
		WHERE email = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		email)
	return result, err == sql.ErrNoRows, err
}

// Create creates user.
func Create(db Connection, firstName, lastName, email, password string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(first_name, last_name, email, password)
		VALUES
		(?,?,?,?)
		`, table),
		firstName, lastName, email, password)
	return result, err
}
