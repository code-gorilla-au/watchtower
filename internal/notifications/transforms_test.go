package notifications

import (
	"database/sql"
	"testing"
	"time"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/odize"
)

func TestTransforms(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("fromNotificationModel should correctly transform database model to Notification struct", func(t *testing.T) {
			now := time.Now().UTC().Truncate(time.Second)
			createdAtTs := now.Add(-1 * time.Hour).Unix()
			updatedAtTs := now.Unix()

			model := database.OrganisationNotification{
				ID: 123,
				OrganisationID: sql.NullInt64{
					Int64: 456,
					Valid: true,
				},
				Type:      "alert",
				Content:   "test content",
				Status:    string(Unread),
				CreatedAt: createdAtTs,
				UpdatedAt: updatedAtTs,
			}

			expected := Notification{
				ID:             123,
				OrganisationID: 456,
				Status:         Unread,
				Content:        "test content",
				Type:           "alert",
				CreatedAt:      time.Unix(createdAtTs, 0).UTC(),
				UpdatedAt:      time.Unix(updatedAtTs, 0).UTC(),
			}

			result := fromNotificationModel(model)

			odize.AssertEqual(t, expected.ID, result.ID)
			odize.AssertEqual(t, expected.OrganisationID, result.OrganisationID)
			odize.AssertEqual(t, expected.Status, result.Status)
			odize.AssertEqual(t, expected.Content, result.Content)
			odize.AssertEqual(t, expected.Type, result.Type)
			odize.AssertEqual(t, expected.CreatedAt, result.CreatedAt)
			odize.AssertEqual(t, expected.UpdatedAt, result.UpdatedAt)
		}).
		Test("fromNotificationModel should handle zero organisation ID", func(t *testing.T) {
			model := database.OrganisationNotification{
				OrganisationID: sql.NullInt64{
					Int64: 0,
					Valid: false,
				},
			}

			result := fromNotificationModel(model)
			odize.AssertEqual(t, int64(0), result.OrganisationID)
		}).
		Run()

	odize.AssertNoError(t, err)
}
