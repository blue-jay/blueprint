package model

import (
	"database/sql"
	"errors"
)

var (
	// ErrNoResult is when no results are found
	ErrNoResult = errors.New("Result not found.")
)

// StrandardError returns a model defined error
func StandardError(err error) error {
	if err == sql.ErrNoRows {
		return ErrNoResult
	}

	return err
}
