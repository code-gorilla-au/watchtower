package notifications

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/code-gorilla-au/odize"
)

func TestService(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()

	group.BeforeEach(func() {
		s = New(_testDB, _testTxnDB, func(tx *sql.Tx) Store {
			return _testDB.WithTx(tx)
		})
	})

	err := group.
		Test("CreateNotification should create a notification", func(t *testing.T) {
			orgID := int64(1)
			notifType := "test-type"
			content := "test-content"

			err := s.CreateNotification(ctx, CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: notifType,
				Content:          content,
				ExternalID:       "test-external-id",
			})
			odize.AssertNoError(t, err)

			unread, err := s.GetNotificationByExternalID(ctx, "test-external-id")
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, orgID, unread.OrganisationID)
			odize.AssertEqual(t, notifType, unread.Type)
			odize.AssertEqual(t, content, unread.Content)
			odize.AssertEqual(t, StatusUnread, unread.Status)
			odize.AssertTrue(t, unread.ID > 0)
			odize.AssertFalse(t, unread.CreatedAt.IsZero())
			odize.AssertFalse(t, unread.UpdatedAt.IsZero())
		}).
		Test("MarkNotificationAsRead should update notification status", func(t *testing.T) {
			orgID := int64(3)
			err := s.CreateNotification(ctx, CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "type",
				Content:          "content",
				ExternalID:       "ext3",
			})
			odize.AssertNoError(t, err)

			notif, err := s.GetNotificationByExternalID(ctx, "ext3")
			odize.AssertNoError(t, err)

			err = s.MarkNotificationAsRead(ctx, notif.ID)
			odize.AssertNoError(t, err)

			updatedNotif, updateErr := s.GetNotificationByExternalID(ctx, "ext3")
			odize.AssertNoError(t, updateErr)
			odize.AssertEqual(t, StatusRead, updatedNotif.Status)
		}).
		Test("DeleteNotificationsByDate should delete old notifications", func(t *testing.T) {
			orgID := int64(4)

			err := s.CreateNotification(ctx, CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "type",
				Content:          "content",
				ExternalID:       "ext4",
			})
			odize.AssertNoError(t, err)

			unread, err := s.GetNotificationByExternalID(ctx, "ext4")
			odize.AssertNoError(t, err)

			err = s.MarkNotificationAsRead(ctx, unread.ID)
			odize.AssertNoError(t, err)

			cutoff := time.Now().Add(1 * time.Minute)

			err = s.DeleteNotificationsByDate(ctx, cutoff)
			odize.AssertNoError(t, err)

			_, err = s.GetNotificationByExternalID(ctx, "ext4")
			odize.AssertError(t, err)
		}).
		Test("BulkCreateNotifications should create multiple notifications", func(t *testing.T) {
			notifications := []CreateNotificationParams{
				{
					OrgID:            1,
					NotificationType: "bulk-type-1",
					Content:          "bulk-content-1",
					ExternalID:       "bulk-ext-1",
				},
				{
					OrgID:            1,
					NotificationType: "bulk-type-2",
					Content:          "bulk-content-2",
					ExternalID:       "bulk-ext-2",
				},
			}

			count, err := s.BulkCreateNotifications(ctx, notifications)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 2, count)

			n1, err := s.GetNotificationByExternalID(ctx, "bulk-ext-1")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "bulk-content-1", n1.Content)

			n2, err := s.GetNotificationByExternalID(ctx, "bulk-ext-2")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "bulk-content-2", n2.Content)
		}).
		Test("BulkCreateNotifications should handle empty list", func(t *testing.T) {
			count, err := s.BulkCreateNotifications(ctx, []CreateNotificationParams{})
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 0, count)
		}).
		Test("BulkCreateNotifications should handle duplicate external IDs by ignoring them", func(t *testing.T) {
			notifications := []CreateNotificationParams{
				{
					OrgID:            1,
					NotificationType: "type",
					Content:          "content-original",
					ExternalID:       "dup-ext",
				},
				{
					OrgID:            1,
					NotificationType: "type",
					Content:          "content-duplicate",
					ExternalID:       "dup-ext",
				},
			}

			count, err := s.BulkCreateNotifications(ctx, notifications)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 2, count)

			n, err := s.GetNotificationByExternalID(ctx, "dup-ext")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "content-original", n.Content)
		}).
		Test("GetUnreadNotifications should return only unread notifications", func(t *testing.T) {
			orgID := int64(5)

			// Create an unread notification
			err := s.CreateNotification(ctx, CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "unread-type",
				Content:          "unread-content",
				ExternalID:       "unread-ext",
			})
			odize.AssertNoError(t, err)

			// Create another notification and mark it as read
			err = s.CreateNotification(ctx, CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "read-type",
				Content:          "read-content",
				ExternalID:       "read-ext",
			})
			odize.AssertNoError(t, err)

			notif, err := s.GetNotificationByExternalID(ctx, "read-ext")
			odize.AssertNoError(t, err)

			err = s.MarkNotificationAsRead(ctx, notif.ID)
			odize.AssertNoError(t, err)

			// Fetch unread notifications
			unread, err := s.GetUnreadNotifications(ctx)
			odize.AssertNoError(t, err)

			// Verify only the unread one is returned
			// Note: The database might contain notifications from previous tests if it's not cleared.
			// However, based on the current implementation of GetUnreadNotifications, it doesn't take orgID.
			// Let's check the service.go again.

			foundUnread := false
			foundRead := false
			for _, n := range unread {
				if n.Content == "unread-content" {
					foundUnread = true
				}
				if n.Content == "read-content" {
					foundRead = true
				}
			}

			odize.AssertTrue(t, foundUnread)
			odize.AssertFalse(t, foundRead)
		}).
		Run()

	odize.AssertNoError(t, err)
}
