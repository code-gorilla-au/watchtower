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

			notif, err := s.CreateNotification(ctx, CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: notifType,
				Content:          content,
				ExternalID:       "test-external-id",
			})
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, orgID, notif.OrganisationID)
			odize.AssertEqual(t, notifType, notif.Type)
			odize.AssertEqual(t, content, notif.Content)
			odize.AssertEqual(t, StatusUnread, notif.Status)
			odize.AssertTrue(t, notif.ID > 0)
			odize.AssertFalse(t, notif.CreatedAt.IsZero())
			odize.AssertFalse(t, notif.UpdatedAt.IsZero())
		}).
		Test("GetUnreadNotifications should return only unread notifications for an org", func(t *testing.T) {
			orgID := int64(2)

			_, err := s.CreateNotification(ctx, CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "type1",
				Content:          "content1",
				ExternalID:       "ext1",
			})
			odize.AssertNoError(t, err)

			notif2, err := s.CreateNotification(ctx, CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "type2",
				Content:          "content2",
				ExternalID:       "ext2",
			})
			odize.AssertNoError(t, err)

			err = s.MarkNotificationAsRead(ctx, notif2.ID)
			odize.AssertNoError(t, err)

			unread, err := s.GetUnreadNotifications(ctx, orgID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, 1, len(unread))
			odize.AssertEqual(t, "content1", unread[0].Content)
			odize.AssertEqual(t, StatusUnread, unread[0].Status)
		}).
		Test("MarkNotificationAsRead should update notification status", func(t *testing.T) {
			orgID := int64(3)
			notif, err := s.CreateNotification(ctx, CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "type",
				Content:          "content",
				ExternalID:       "ext3",
			})
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, StatusUnread, notif.Status)

			err = s.MarkNotificationAsRead(ctx, notif.ID)
			odize.AssertNoError(t, err)

			unread, err := s.GetUnreadNotifications(ctx, orgID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 0, len(unread))
		}).
		Test("DeleteNotificationsByDate should delete old notifications", func(t *testing.T) {
			orgID := int64(4)

			_, err := s.CreateNotification(ctx, CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "type",
				Content:          "content",
				ExternalID:       "ext4",
			})
			odize.AssertNoError(t, err)

			cutoff := time.Now().Add(1 * time.Minute)

			err = s.DeleteNotificationsByDate(ctx, cutoff)
			odize.AssertNoError(t, err)

			unread, err := s.GetUnreadNotifications(ctx, orgID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 0, len(unread))
		}).
		Run()

	odize.AssertNoError(t, err)
}
