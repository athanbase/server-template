// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package queries

import (
	"time"
)

// user
type User struct {
	ID int64 `json:"id"`
	// user name
	Name string `json:"name"`
	// user country
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// user detail
type UserDetail struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	// user email
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
