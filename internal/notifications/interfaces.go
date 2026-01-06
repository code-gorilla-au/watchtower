package notifications

import (
	"context"

	"watchtower/internal/database"
)

type Store interface {
	CreateOrgNotification(ctx context.Context, arg database.CreateOrgNotificationParams) (database.OrganisationNotification, error)
	UpdateOrgNotificationByID(ctx context.Context, arg database.UpdateOrgNotificationByIDParams) (database.OrganisationNotification, error)
	UpdateOrgNotificationStatusByID(ctx context.Context, arg database.UpdateOrgNotificationStatusByIDParams) (database.OrganisationNotification, error)
}

var _ Store = (*database.Queries)(nil)
