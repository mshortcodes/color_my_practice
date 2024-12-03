// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID         uuid.UUID
	Date       time.Time
	ColorDepth int32
	Confirmed  bool
	UserID     uuid.UUID
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	HashedPassword string
}
