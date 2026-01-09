package watchtower

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"watchtower/internal/database"
	"watchtower/internal/notifications"
	"watchtower/internal/organisations"
	"watchtower/internal/products"

	"watchtower/internal/github"
	"watchtower/internal/logging"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// NewService creates and returns a new Service instance with the provided database queries.
func NewService(ctx context.Context, db *database.Queries, txnDB *sql.DB) *Service {
	return &Service{
		ghClient: github.New(logging.FromContext(ctx)),
		ctx:      ctx,
		orgSvc: organisations.New(db, txnDB, func(tx *sql.Tx) organisations.OrgStore {
			return db.WithTx(tx)
		}),
		productSvc: products.New(db),
		notificationSvc: notifications.New(db, txnDB, func(tx *sql.Tx) notifications.Store {
			return db.WithTx(tx)
		}),
	}
}

func (s *Service) Startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *Service) CreateUnreadNotification() error {
	logger := logging.FromContext(s.ctx)

	prCount, err := s.CreateUnreadPRNotification()
	if err != nil {
		logger.Error("Error creating unread pull request notification", "error", err)
		return err
	}

	secCount, err := s.CreateUnreadSecurityNotification()
	if err != nil {
		logger.Error("Error creating unread security notification", "error", err)
		return err
	}

	if (prCount + secCount) > 0 {
		runtime.EventsEmit(s.ctx, "UNREAD_NOTIFICATIONS")
	}

	return nil
}

// CreateUnreadPRNotification generates unread notifications for recent pull requests by fetching their IDs and creating notifications.
func (s *Service) CreateUnreadPRNotification() (int, error) {
	logger := logging.FromContext(s.ctx)

	prs, err := s.productSvc.GetRecentPullRequests(s.ctx)
	if err != nil {
		logger.Error("Error fetching recent pull requests", "error", err)
		return 0, err
	}

	return s.createNotification("OPEN_PULL_REQUEST", "New pull request", prs)
}

// CreateUnreadSecurityNotification generates unread security notifications for recent security alerts.
// It retrieves recent security-related IDs and creates notifications for each using the notification service.
// Returns an error if fetching security IDs or creating notifications fails.
func (s *Service) CreateUnreadSecurityNotification() (int, error) {
	logger := logging.FromContext(s.ctx)
	secList, err := s.productSvc.GetRecentSecurity(s.ctx)
	if err != nil {
		logger.Error("Error fetching recent security", "error", err)
		return 0, err
	}

	return s.createNotification("OPEN_SECURITY_ALERT", "New security alert", secList)
}

func (s *Service) createNotification(notificationType string, content string, recentlyChanged []products.RecentlyChangedEntity) (int, error) {

	notificationsList := make([]notifications.CreateNotificationParams, len(recentlyChanged))
	for i, entity := range recentlyChanged {
		notificationsList[i] = notifications.CreateNotificationParams{
			OrgID:            entity.OrganisationID,
			ExternalID:       entity.ExternalID,
			NotificationType: notificationType,
			Content:          fmt.Sprintf("%s: %s", entity.RepositoryName, content),
		}
	}

	return s.notificationSvc.BulkCreateNotifications(s.ctx, notificationsList)
}

// SyncOrgs synchronizes stale organisations by retrieving them and invoking the sync process for each.
func (s *Service) SyncOrgs() error {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Syncing orgs")

	orgs, err := s.orgSvc.GetStaleOrgs(s.ctx)
	if err != nil {
		logger.Error("Error fetching orgs", "error", err)

		return err
	}

	logger.Debug("syncing number of orgs", "count", len(orgs))

	for _, org := range orgs {
		if err = s.SyncOrg(org.ID); err != nil {
			logger.Error("Error syncing org", "error", err)

			continue
		}
	}

	return nil
}

// SyncOrg synchronizes the products and associated organization for the given organization ID.
func (s *Service) SyncOrg(orgId int64) error {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Syncing org", "org", orgId)

	prodList, err := s.GetAllProductsForOrganisation(orgId)
	if err != nil {
		logger.Error("Error fetching products for org", "error", err)

		return err
	}

	if len(prodList) == 0 {
		logger.Debug("No products found for org", "org", orgId)
		return nil
	}

	org, err := s.orgSvc.GetOrgAssociatedToProduct(s.ctx, prodList[0].ID)
	if err != nil {
		logger.Error("Error fetching organisation for product", "error", err)

		return err
	}

	for _, p := range prodList {
		if err = s.syncProductFromGithub(p, org); err != nil {
			logger.Error("Error syncing product", "error", err)

			return err
		}
	}

	if err = s.productSvc.UpdateSyncDateNow(s.ctx, org.ID); err != nil {
		logger.Error("Error updating organisation sync", "error", err)

		return err
	}

	return nil
}

// SyncProduct synchronizes a product with the given ID by retrieving its details and associated organization data.
func (s *Service) SyncProduct(id int64) error {
	logger := logging.FromContext(s.ctx)

	product, err := s.GetProductByID(id)
	if err != nil {
		logger.Error("Error fetching product", "error", err)

		return err
	}

	org, err := s.orgSvc.GetOrgAssociatedToProduct(s.ctx, product.ID)
	if err != nil {
		logger.Error("Error fetching organisation for product", "error", err)

		return err
	}

	return s.syncProductFromGithub(product, org)
}

// syncProductFromGithub synchronizes product repository data from GitHub for the specified product and organization.
// The method iterates through the product's tags to fetch repository data and updates the product's sync date upon success.
// Returns an error if syncing repositories or updating the sync date fails.
func (s *Service) syncProductFromGithub(product products.ProductDTO, org organisations.InternalOrganisation) error {
	logger := logging.FromContext(s.ctx)

	for _, tag := range product.Tags {
		if err := s.syncRepoDataByTag(tag, org.Namespace, org.Token); err != nil {
			logger.Error("Error syncing repos", "error", err)

			return err
		}
	}

	if err := s.productSvc.UpdateSyncDateNow(s.ctx, product.ID); err != nil {
		logger.Error("Error updating product sync", "error", err)

		return err
	}

	return nil
}

// syncRepoDataByTag synchronizes repository data based on a given tag by searching, retrieving, and inserting repository details.
func (s *Service) syncRepoDataByTag(tag string, owner string, ghToken string) error {
	logger := logging.FromContext(s.ctx)

	logger.Debug("Searching for repo with tag", "tag", tag)

	repos, apiErr := s.ghClient.SearchRepos(owner, strings.TrimSpace(tag), ghToken)
	if apiErr != nil {
		logger.Error("Error searching for repos", "error", apiErr)

		return apiErr
	}

	if err := s.productSvc.BulkInsertRepos(s.ctx, repos.Data.Search.Edges, tag); err != nil {
		logger.Error("Error bulk inserting repos", "error", err)

		return err
	}

	for _, repo := range repos.Data.Search.Edges {
		dd, err := s.ghClient.GetRepoDetails(owner, repo.Node.Name, ghToken)
		if err != nil {
			logger.Error("Error getting repo details", "repo", repo.Node.Name, "error", err)
		}

		if err = s.productSvc.BulkInsertRepoDetails(s.ctx, dd); err != nil {
			logger.Error("Error bulk inserting repo details", "repo", repo.Node.Name, "error", err)
		}
	}

	return nil
}
