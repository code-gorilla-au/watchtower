package watchtower

import (
	"context"
	"strings"
	"testing"
	"time"
	"watchtower/internal/notifications"
	"watchtower/internal/products"

	"github.com/code-gorilla-au/odize"
	"github.com/google/uuid"
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
			orgID := int64(0)

			err := s.notificationSvc.CreateNotification(ctx, notifications.CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "test-type",
				Content:          "test-content",
				ExternalID:       uuid.New().String(),
			})
			odize.AssertNoError(t, err)

			unread, err := s.GetUnreadNotifications()
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, 1, len(unread))
			odize.AssertEqual(t, "test-content", unread[0].Content)
			odize.AssertEqual(t, notifications.StatusUnread, unread[0].Status)
		}).
		Test("MarkNotificationAsRead should mark a notification as read", func(t *testing.T) {
			orgID := int64(1002)
			err := s.notificationSvc.CreateNotification(ctx, notifications.CreateNotificationParams{
				OrgID:            orgID,
				NotificationType: "type",
				Content:          "content",
				ExternalID:       "test-external-id-2",
			})
			odize.AssertNoError(t, err)

			unread, err := s.notificationSvc.GetNotificationByExternalID(ctx, "test-external-id-2")
			odize.AssertNoError(t, err)

			err = s.MarkNotificationAsRead(unread.ID)
			odize.AssertNoError(t, err)

			verifyUnread, err := s.notificationSvc.GetNotificationByExternalID(ctx, "test-external-id-2")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, notifications.StatusRead, verifyUnread.Status)
		}).
		Test("DeleteOldNotifications should delete notifications", func(t *testing.T) {
			orgID := int64(1003)
			err := s.notificationSvc.CreateNotification(ctx, notifications.CreateNotificationParams{
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

			unread, err := s.GetUnreadNotifications()
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, 0, len(unread))
		}).
		Test("CreateUnreadPRNotification should create notifications for recent PRs", func(t *testing.T) {
			// Setup: Org, Product, Repo, PR
			org, err := s.CreateOrganisation("Test Org PR", "test-org-pr", "token", "desc")
			odize.AssertNoError(t, err)

			_, err = s.CreateProduct("Test Product PR", "desc", []string{"tag-pr"}, org.ID)
			odize.AssertNoError(t, err)

			err = s.productSvc.CreateRepo(ctx, products.CreateRepoParams{
				Name:  "repo-pr",
				Url:   "url",
				Topic: "tag-pr",
				Owner: "owner",
			})
			odize.AssertNoError(t, err)

			externalID := uuid.New().String()
			err = s.productSvc.CreatePullRequest(ctx, products.CreatePRParams{
				ExternalID:     externalID,
				Title:          "PR Title",
				RepositoryName: "repo-pr",
				Url:            "url",
				State:          "OPEN",
				Author:         "author",
				CreatedAt:      time.Now(),
			})
			odize.AssertNoError(t, err)

			// Action
			err = s.CreateUnreadPRNotification()
			odize.AssertNoError(t, err)

			// Verify
			unreadNotification, err := s.notificationSvc.GetNotificationByExternalID(ctx, externalID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, "OPEN_PULL_REQUEST", unreadNotification.Type)
			odize.AssertTrue(t, strings.Contains(unreadNotification.Content, "repo-pr"))
			odize.AssertTrue(t, strings.Contains(unreadNotification.Content, "New pull request"))
		}).
		Test("CreateUnreadSecurityNotification should create notifications for recent security alerts", func(t *testing.T) {
			// Setup: Org, Product, Repo, Security Alert
			org, err := s.CreateOrganisation("Test Org Sec", "test-org-sec", "token", "desc")
			odize.AssertNoError(t, err)

			_, err = s.CreateProduct("Test Product Sec", "desc", []string{"tag-sec"}, org.ID)
			odize.AssertNoError(t, err)

			err = s.productSvc.CreateRepo(ctx, products.CreateRepoParams{
				Name:  "repo-sec",
				Url:   "url",
				Topic: "tag-sec",
				Owner: "owner",
			})
			odize.AssertNoError(t, err)

			externalID := uuid.New().String()
			err = s.productSvc.UpsertSecurity(ctx, products.CreateSecurityParams{
				ExternalID:     externalID,
				RepositoryName: "repo-sec",
				PackageName:    "pkg",
				State:          "OPEN",
				Severity:       "HIGH",
				CreatedAt:      time.Now(),
			})
			odize.AssertNoError(t, err)

			// Action
			err = s.CreateUnreadSecurityNotification()
			odize.AssertNoError(t, err)

			// Verify
			unreadNotification, err := s.notificationSvc.GetNotificationByExternalID(ctx, externalID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, "OPEN_SECURITY_ALERT", unreadNotification.Type)
			odize.AssertTrue(t, strings.Contains(unreadNotification.Content, "repo-sec"))
			odize.AssertTrue(t, strings.Contains(unreadNotification.Content, "New security alert"))
		}).
		Run()

	odize.AssertNoError(t, err)
}
