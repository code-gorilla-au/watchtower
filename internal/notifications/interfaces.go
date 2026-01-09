package notifications

import (
	"context"
	"watchtower/internal/database"
)

type Store interface {
	CreateOrgNotification(ctx context.Context, arg database.CreateOrgNotificationParams) (database.OrganisationNotification, error)
	UpdateOrgNotificationByID(ctx context.Context, arg database.UpdateOrgNotificationByIDParams) (database.OrganisationNotification, error)
	UpdateOrgNotificationStatusByID(ctx context.Context, arg database.UpdateOrgNotificationStatusByIDParams) (database.OrganisationNotification, error)
	GetNotificationByExternalID(ctx context.Context, externalID string) (database.OrganisationNotification, error)
	GetUnreadNotifications(ctx context.Context) ([]database.OrganisationNotification, error)
	DeleteOrgNotificationByDate(ctx context.Context, createdAt int64) error
}

var _ Store = (*database.Queries)(nil)
