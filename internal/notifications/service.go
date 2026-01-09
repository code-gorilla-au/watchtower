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

type CreateNotificationParams struct {
	OrgID            int64
	ExternalID       string
	NotificationType string
	Content          string
}

// CreateNotification creates a new notification for a specific organisation with the given type and content.
func (s *Service) CreateNotification(ctx context.Context, params CreateNotificationParams) error {
	logger := logging.FromContext(ctx).With("orgID", params.OrgID, "service", "notifications")
	logger.Debug("Creating notification for org")

	_, err := s.store.CreateOrgNotification(ctx, database.CreateOrgNotificationParams{
		ExternalID: params.ExternalID,
		OrganisationID: sql.NullInt64{
			Int64: params.OrgID,
			Valid: true,
		},
		Type:    params.NotificationType,
		Content: params.Content,
	})
	if err == nil || database.IsErrUniqueConstraint(err) {
		return nil
	}

	logger.Error("Error creating notification", "error", err)
	return err
}

// BulkCreateNotifications creates multiple notifications in a single operation and returns the count or an error if any fail.
func (s *Service) BulkCreateNotifications(ctx context.Context, notifications []CreateNotificationParams) (int, error) {
	logger := logging.FromContext(ctx).With("service", "notifications")
	logger.Debug("Creating notifications in bulk")

	for _, item := range notifications {
		err := s.CreateNotification(ctx, CreateNotificationParams{
			ExternalID:       item.ExternalID,
			OrgID:            item.OrgID,
			NotificationType: item.NotificationType,
			Content:          item.Content,
		})
		if err != nil {
			logger.Error("Error creating notification", "error", err)
			return 0, err
		}
	}

	return len(notifications), nil
}

// GetUnreadNotifications fetches all unread notifications for the specified organisation ID. Returns a list of notifications or an error.
func (s *Service) GetUnreadNotifications(ctx context.Context) ([]Notification, error) {
	logger := logging.FromContext(ctx).With("service", "notifications")
	logger.Debug("Fetching unread notifications")

	models, err := s.store.GetUnreadNotifications(ctx)
	if err != nil {
		logger.Error("Error fetching unread notifications", "error", err)
		return []Notification{}, err
	}

	return fromNotificationModels(models), nil
}

func (s *Service) GetNotificationByExternalID(ctx context.Context, externalID string) (Notification, error) {
	logger := logging.FromContext(ctx).With("externalID", externalID, "service", "notifications")
	logger.Debug("Fetching notification by external ID")

	model, err := s.store.GetNotificationByExternalID(ctx, externalID)
	if err != nil {
		logger.Error("Error fetching notification by external ID", "error", err)
		return Notification{}, err
	}

	return fromNotificationModel(model), nil
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
