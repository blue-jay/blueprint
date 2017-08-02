package model

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

// User defines the model.
type User struct {
	Connection
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

// NewUser returns a user.
func NewUser(db Connection) User {
	return User{
		Connection: db,
	}
}

// ByEmail gets user information from email.
func (db User) ByEmail(email string) (User, bool, error) {
	result := User{}
	err := db.Get(&result, `
		SELECT id, password, status_id, first_name
		FROM user
		WHERE email = ?
			AND deleted_at IS NULL
		LIMIT 1
		`,
		email)
	return result, err == sql.ErrNoRows, err
}

// Create creates user.
func (db User) Create(firstName, lastName, email, password string) (sql.Result, error) {
	result, err := db.Exec(`
		INSERT INTO user
		(first_name, last_name, email, password)
		VALUES
		(?,?,?,?)
		`,
		firstName, lastName, email, password)
	return result, err
}
