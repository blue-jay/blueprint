// Package model handles the loading of models.
package model

import (
	"github.com/blue-jay/blueprint/model/note"
	"github.com/blue-jay/blueprint/model/user"
	"github.com/jmoiron/sqlx"
)

var (
	Note note.Service // Note model
	User user.Service // User model
)

// Load injects the dependencies for the models
func Load(db *sqlx.DB) {
	Note = note.Service{db}
	User = user.Service{db}
}
