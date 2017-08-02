package model

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// Note defines the model.
type Note struct {
	Connection
	ID        uint32         `db:"id"`
	Name      string         `db:"name"`
	UserID    uint32         `db:"user_id"`
	CreatedAt mysql.NullTime `db:"created_at"`
	UpdatedAt mysql.NullTime `db:"updated_at"`
	DeletedAt mysql.NullTime `db:"deleted_at"`
}

// NewNote returns a note.
func NewNote(db Connection) Note {
	return Note{
		Connection: db,
	}
}

// ByID gets an item by ID.
func (db Note) ByID(ID string, userID string) (Note, bool, error) {
	result := Note{}
	err := db.Get(&result, `
		SELECT id, name, user_id, created_at, updated_at, deleted_at
		FROM note
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`,
		ID, userID)
	return result, err == sql.ErrNoRows, err
}

// ByUserID gets all items for a user.
func (db Note) ByUserID(userID string) ([]Note, bool, error) {
	var result []Note
	err := db.Select(&result, `
		SELECT id, name, user_id, created_at, updated_at, deleted_at
		FROM note
		WHERE user_id = ?
			AND deleted_at IS NULL
		`,
		userID)
	return result, err == sql.ErrNoRows, err
}

// ByUserIDPaginate gets items for a user based on page and max variables.
func (db Note) ByUserIDPaginate(userID string, max int, page int) ([]Note, bool, error) {
	var result []Note
	err := db.Select(&result, fmt.Sprintf(`
		SELECT id, name, user_id, created_at, updated_at, deleted_at
		FROM note
		WHERE user_id = ?
			AND deleted_at IS NULL
		LIMIT %v OFFSET %v
		`, max, page),
		userID)
	return result, err == sql.ErrNoRows, err
}

// ByUserIDCount counts the number of items for a user.
func (db Note) ByUserIDCount(userID string) (int, error) {
	var result int
	err := db.Get(&result, `
		SELECT count(*)
		FROM note
		WHERE user_id = ?
			AND deleted_at IS NULL
		`,
		userID)
	return result, err
}

// Create adds an item.
func (db Note) Create(name string, userID string) (sql.Result, error) {
	result, err := db.Exec(`
		INSERT INTO note
		(name, user_id)
		VALUES
		(?,?)
		`,
		name, userID)
	return result, err
}

// Update makes changes to an existing item.
func (db Note) Update(name string, ID string, userID string) (sql.Result, error) {
	result, err := db.Exec(`
		UPDATE note
		SET name = ?
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`,
		name, ID, userID)
	return result, err
}

// DeleteHard removes an item.
func (db Note) DeleteHard(ID string, userID string) (sql.Result, error) {
	result, err := db.Exec(`
		DELETE FROM note
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		`,
		ID, userID)
	return result, err
}

// DeleteSoft marks an item as removed.
func (db Note) DeleteSoft(ID string, userID string) (sql.Result, error) {
	result, err := db.Exec(`
		UPDATE note
		SET deleted_at = NOW()
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`,
		ID, userID)
	return result, err
}
