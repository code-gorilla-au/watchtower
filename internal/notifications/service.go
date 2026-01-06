package notifications

import (
	"context"
	"database/sql"
	"time"
	"watchtower/internal/database"
	"watchtower/internal/logging"
)

func New(db Store, txnDB *sql.DB, txnFunc func(tx *sql.Tx) Store) *Service {
	return &Service{
		store:   db,
		txnDB:   txnDB,
		txnFunc: txnFunc,
	}
}

// CreateNotification creates a new notification for a specific organisation with the given type and content.
func (s *Service) CreateNotification(ctx context.Context, orgID int64, notificationType string, content string) (Notification, error) {
	logger := logging.FromContext(ctx).With("orgID", orgID, "service", "notifications")
	logger.Debug("Creating notification for org")

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

// GetUnreadNotifications fetches all unread notifications for the specified organisation ID. Returns a list of notifications or an error.
func (s *Service) GetUnreadNotifications(ctx context.Context, orgID int64) ([]Notification, error) {
	logger := logging.FromContext(ctx).With("orgID", orgID, "service", "notifications")
	logger.Debug("Fetching unread notifications")

	models, err := s.store.GetUnreadNotificationsByOrgID(ctx, sql.NullInt64{
		Int64: orgID,
		Valid: true,
	})
	if err != nil {
		logger.Error("Error fetching unread notifications", "error", err)
		return []Notification{}, err
	}

	return fromNotificationModels(models), nil
}

// MarkNotificationAsRead updates the status of a notification to "read" based on the provided notification ID.
func (s *Service) MarkNotificationAsRead(ctx context.Context, notificationID int64) error {
	logger := logging.FromContext(ctx).With("notificationID", notificationID, "service", "notifications")
	logger.Debug("Marking notification as read")

	if _, err := s.store.UpdateOrgNotificationStatusByID(ctx, database.UpdateOrgNotificationStatusByIDParams{
		Status: string(StatusRead),
		ID:     notificationID,
	}); err != nil {
		logger.Error("could not update notification status", "error", err)
		return err
	}

	return nil
}

// DeleteNotificationsByDate removes notifications older than the specified date from the store. Returns an error if deletion fails.
func (s *Service) DeleteNotificationsByDate(ctx context.Context, date time.Time) error {
	logger := logging.FromContext(ctx).With("date", date, "service", "notifications")
	logger.Debug("Deleting notifications older than date")

	if err := s.store.DeleteOrgNotificationByDate(ctx, date.Unix()); err != nil {
		logger.Error("Error deleting notifications", "error", err)
		return err
	}

	return nil
}
