package notifications

import (
	"context"
	"database/sql"

	"watchtower/internal/database"
)

type Store interface {
	CreateOrgNotification(ctx context.Context, arg database.CreateOrgNotificationParams) (database.OrganisationNotification, error)
	UpdateOrgNotificationByID(ctx context.Context, arg database.UpdateOrgNotificationByIDParams) (database.OrganisationNotification, error)
	UpdateOrgNotificationStatusByID(ctx context.Context, arg database.UpdateOrgNotificationStatusByIDParams) (database.OrganisationNotification, error)
	GetUnreadNotificationsByOrgID(ctx context.Context, organisationID sql.NullInt64) ([]database.OrganisationNotification, error)
	DeleteOrgNotificationByDate(ctx context.Context, createdAt int64) error
}

var _ Store = (*database.Queries)(nil)
