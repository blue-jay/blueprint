// Package user provides access to the user table in the MySQL database.
package user

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

// Service defines the database connection.
type Service struct {
	DB *sqlx.DB
}

// ByEmail gets user information from email.
func (c Service) ByEmail(email string) (Item, bool, error) {
	result := Item{}
	err := c.DB.Get(&result, fmt.Sprintf(`
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
func (c Service) Create(firstName, lastName, email, password string) (sql.Result, error) {
	result, err := c.DB.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(first_name, last_name, email, password)
		VALUES
		(?,?,?,?)
		`, table),
		firstName, lastName, email, password)
	return result, err
}
