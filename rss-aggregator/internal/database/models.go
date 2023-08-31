// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type TbFeed struct {
	ID        uuid.UUID
	Name      sql.NullString
	Url       string
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TbFeedFollow struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FeedID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TbUser struct {
	ID        uuid.UUID
	Name      sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
	ApiKey    string
}
