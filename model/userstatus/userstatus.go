package userstatus

import "time"

var (
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
