package notifications

import (
	"database/sql"
	"time"
)

type Service struct {
	store   Store
	txnDB   *sql.DB
	txnFunc func(tx *sql.Tx) Store
}

type NotificationStatus string

const (
	StatusUnread NotificationStatus = "unread"
	StatusRead   NotificationStatus = "read"
)

type Notification struct {
	ID             int64              `json:"id"`
	OrganisationID int64              `json:"organisation_id"`
	Status         NotificationStatus `json:"status"`
	Content        string             `json:"content"`
	Type           string             `json:"type"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}
