package watchtower

import (
	"context"
	"testing"
	"time"
	"watchtower/internal/notifications"

	"github.com/code-gorilla-au/odize"
)

func TestService_Notifications(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
	})

	err := group.
		Test("GetUnreadNotifications should return unread notifications", func(t *testing.T) {
			orgID := int64(1001)

			// Seed a notification using the internal notification service
			notif, err := s.notificationSvc.CreateNotification(ctx, notifications.CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "test-type",
				Content:          "test-content",
				ExternalID:       "test-external-id-1",
			})
			odize.AssertNoError(t, err)

			// Fetch unread notifications
			unread, err := s.GetUnreadNotifications(orgID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, 1, len(unread))
			odize.AssertEqual(t, notif.ID, unread[0].ID)
			odize.AssertEqual(t, "test-content", unread[0].Content)
			odize.AssertEqual(t, notifications.StatusUnread, unread[0].Status)
		}).
		Test("MarkNotificationAsRead should mark a notification as read", func(t *testing.T) {
			orgID := int64(1002)
			notif, err := s.notificationSvc.CreateNotification(ctx, notifications.CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "type",
				Content:          "content",
				ExternalID:       "test-external-id-2",
			})
			odize.AssertNoError(t, err)

			err = s.MarkNotificationAsRead(notif.ID)
			odize.AssertNoError(t, err)

			unread, err := s.GetUnreadNotifications(orgID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 0, len(unread))
		}).
		Test("DeleteOldNotifications should delete notifications", func(t *testing.T) {
			orgID := int64(1003)
			_, err := s.notificationSvc.CreateNotification(ctx, notifications.CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "type",
				Content:          "content",
				ExternalID:       "test-external-id-3",
			})
			odize.AssertNoError(t, err)

			// DeleteOldNotifications uses time.Now() and the query uses created_at < ?.
			// Since we just created it in the same second, it might not be deleted.
			// Let's wait 1 second to ensure created_at < time.Now().
			time.Sleep(1100 * time.Millisecond)

			err = s.DeleteOldNotifications()
			odize.AssertNoError(t, err)

			unread, err := s.GetUnreadNotifications(orgID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, 0, len(unread))
		}).
		Run()

	odize.AssertNoError(t, err)
}
