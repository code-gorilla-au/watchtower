package watchtower

import (
	"time"
	"watchtower/internal/notifications"
)

// GetUnreadNotifications retrieves a list of unread notifications for the specified organization ID.
func (s *Service) GetUnreadNotifications() ([]notifications.Notification, error) {
	return s.notificationSvc.GetUnreadNotifications(s.ctx)
}

// MarkNotificationAsRead marks a notification as read based on the provided notification ID.
func (s *Service) MarkNotificationAsRead(notificationID int64) error {
	return s.notificationSvc.MarkNotificationAsRead(s.ctx, notificationID)
}

func (s *Service) DeleteOldNotifications() error {
	return s.notificationSvc.DeleteNotificationsByDate(s.ctx, time.Now())
}
