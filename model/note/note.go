// Package note provides access to the note table in the MySQL database.
package note

import (
	"database/sql"
	"fmt"

	"github.com/blue-jay/blueprint/lib/database"
	"github.com/blue-jay/blueprint/model"
	"github.com/go-sql-driver/mysql"
)

var (
	table = "note"
)

// Item defines the model.
type Item struct {
	ID        uint32         `db:"id"`
	Name      string         `db:"name"`
	Filename  string         `db:"filename"`
	FileID    string         `db:"file_id"`
	UserID    uint32         `db:"user_id"`
	CreatedAt mysql.NullTime `db:"created_at"`
	UpdatedAt mysql.NullTime `db:"updated_at"`
	DeletedAt mysql.NullTime `db:"deleted_at"`
}

// ByID gets item by ID.
func ByID(ID string, userID string) (Item, error) {
	result := Item{}
	err := database.SQL.Get(&result, fmt.Sprintf(`
		SELECT id, name, filename, file_id, user_id, created_at, updated_at, deleted_at
		FROM %v
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		ID, userID)
	return result, model.StandardError(err)
}

// ByUserID gets all entities for a user.
func ByUserID(userID string) ([]Item, error) {
	var result []Item
	err := database.SQL.Select(&result, fmt.Sprintf(`
		SELECT id, name, filename, file_id, user_id, created_at, updated_at, deleted_at
		FROM %v
		WHERE user_id = ?
			AND deleted_at IS NULL
		`, table),
		userID)
	return result, model.StandardError(err)
}

// Create adds an item.
func Create(name string, filename string, fileID string, userID string) (sql.Result, error) {
	result, err := database.SQL.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(name, filename, file_id, user_id)
		VALUES
		(?,?,?,?)
		`, table),
		name, filename, fileID, userID)
	return result, model.StandardError(err)
}

// Update makes changes to an existing item.
func Update(name string, ID string, userID string) (sql.Result, error) {
	result, err := database.SQL.Exec(fmt.Sprintf(`
		UPDATE %v
		SET name = ?
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		name, ID, userID)
	return result, model.StandardError(err)
}

// DeleteHard removes an item.
func DeleteHard(ID string, userID string) (sql.Result, error) {
	result, err := database.SQL.Exec(fmt.Sprintf(`
		DELETE FROM %v
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		`, table),
		ID, userID)
	return result, model.StandardError(err)
}

// Delete marks an item as removed.
func Delete(ID string, userID string) (sql.Result, error) {
	result, err := database.SQL.Exec(fmt.Sprintf(`
		UPDATE %v
		SET deleted_at = NOW()
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		ID, userID)
	return result, model.StandardError(err)
}
