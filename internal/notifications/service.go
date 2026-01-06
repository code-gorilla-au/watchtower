package notifications

import (
	"context"
	"database/sql"
	"watchtower/internal/database"
	"watchtower/internal/logging"
)

func NewService(db *database.Queries, txnDB *sql.DB) *Service {
	return &Service{
		store: db,
		txnDB: txnDB,
		txnFunc: func(tx *sql.Tx) Store {
			return db.WithTx(tx)
		},
	}
}

func (s *Service) CreateNotification(ctx context.Context, orgID int64, notificationType string, content string) (Notification, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Creating notification")

	model, err := s.store.CreateOrgNotification(ctx, database.CreateOrgNotificationParams{
		OrganisationID: sql.NullInt64{
			Int64: orgID,
			Valid: true,
		},
		Type:    notificationType,
		Content: content,
	})
	if err != nil {
		logger.Error("Error creating notification", "error", err)
		return Notification{}, err
	}

	return fromNotificationModel(model), nil
}
