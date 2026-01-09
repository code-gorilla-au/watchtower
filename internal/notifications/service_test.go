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
		Run()

	odize.AssertNoError(t, err)
}
