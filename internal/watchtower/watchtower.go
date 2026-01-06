package watchtower

import (
	"context"
	"database/sql"
	"strings"
	"watchtower/internal/database"
	"watchtower/internal/organisations"
	"watchtower/internal/products"

	"watchtower/internal/github"
	"watchtower/internal/logging"
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
	}
}

func (s *Service) Startup(ctx context.Context) {
	s.ctx = ctx
}

// CreateOrganisation creates a new organisation in the database using the specified friendly name and namespace.
// It logs the creation process and returns the created organisation DTO or an error if the operation fails.
func (s *Service) CreateOrganisation(friendlyName string, namespace string, token string, description string) (organisations.OrganisationDTO, error) {
	return s.orgSvc.Create(s.ctx, organisations.CreateOrgParams{
		FriendlyName: friendlyName,
		Namespace:    namespace,
		Token:        token,
		Description:  description,
	})
}

// GetDefaultOrganisation returns the default organisation (marked as default_org = 1).
func (s *Service) GetDefaultOrganisation() (organisations.OrganisationDTO, error) {
	return s.orgSvc.GetDefault(s.ctx)
}

func (s *Service) SetDefaultOrg(id int64) (organisations.OrganisationDTO, error) {
	return s.orgSvc.SetDefault(s.ctx, id)
}

func (s *Service) GetOrganisationByID(id int64) (organisations.OrganisationDTO, error) {
	return s.orgSvc.Get(s.ctx, id)
}

// GetAllOrganisations returns all organisations ordered by friendly_name.
func (s *Service) GetAllOrganisations() ([]organisations.OrganisationDTO, error) {
	return s.orgSvc.GetAll(s.ctx)
}

func (s *Service) DeleteAllOrgs() error {
	list, err := s.orgSvc.GetAll(s.ctx)
	if err != nil {
		return err
	}

	for _, org := range list {
		_ = s.DeleteOrganisation(org.ID)
	}

	return nil
}

// DeleteOrganisation deletes an organisation from the database by its ID. Returns an error if the operation fails.
func (s *Service) DeleteOrganisation(id int64) error {
	list, err := s.productSvc.GetByOrg(s.ctx, id)
	if err != nil {
		return err
	}

	for _, p := range list {
		_ = s.productSvc.DeleteProduct(s.ctx, p.ID)
	}

	return s.orgSvc.Delete(s.ctx, id)
}

// UpdateOrganisation updates the details of an organisation based on the provided parameters.
// It sets other organisations' default status to false if a new default organisation is specified.
// Returns the updated organisation details as OrganisationDTO or an error if the update fails.
func (s *Service) UpdateOrganisation(params organisations.UpdateOrgParams) (organisations.OrganisationDTO, error) {
	return s.orgSvc.Update(s.ctx, params)
}

// CreateProduct creates a new product and associates it with an organisation.
func (s *Service) CreateProduct(name string, description string, tags []string, organisationID int64) (products.ProductDTO, error) {

	prod, err := s.productSvc.Create(s.ctx, products.CreateProductParams{
		Name: name,
		Tags: tags,
		Desc: description,
	})
	if err != nil {
		return products.ProductDTO{}, err
	}

	err = s.orgSvc.AssociateProductToOrg(s.ctx, organisationID, prod.ID)
	if err != nil {
		return products.ProductDTO{}, err
	}

	return prod, nil
}

// GetProductByID fetches a product by its ID.
func (s *Service) GetProductByID(id int64) (products.ProductDTO, error) {
	return s.productSvc.Get(s.ctx, id)
}

// GetAllProductsForOrganisation lists products linked to the given organisation.
func (s *Service) GetAllProductsForOrganisation(organisationID int64) ([]products.ProductDTO, error) {
	return s.productSvc.GetByOrg(s.ctx, organisationID)
}

// UpdateProduct updates a product and returns the updated entity.
func (s *Service) UpdateProduct(id int64, name string, tags []string) (products.ProductDTO, error) {
	return s.productSvc.Update(s.ctx, products.UpdateProductParams{
		ID:   id,
		Name: name,
		Tags: tags,
	})
}

// DeleteProduct removes a product from the database by its ID. Returns an error if the operation fails.
func (s *Service) DeleteProduct(id int64) error {
	return s.productSvc.DeleteProduct(s.ctx, id)
}

// GetProductRepos retrieves a list of repositories associated with the product identified by the given ID.
func (s *Service) GetProductRepos(id int64) ([]products.RepositoryDTO, error) {
	return s.productSvc.GetRepos(s.ctx, id)
}

// GetProductPullRequests retrieves a list of pull requests for a specified product ID.
func (s *Service) GetProductPullRequests(id int64) ([]products.PullRequestDTO, error) {
	return s.productSvc.GetPullRequests(s.ctx, id)
}

// GetPullRequestByOrganisation retrieves a list of pull requests for a given organization by its ID.
func (s *Service) GetPullRequestByOrganisation(id int64) ([]products.PullRequestDTO, error) {
	return s.productSvc.GetPullRequestByOrg(s.ctx, id)
}

// GetSecurityByProductID retrieves a list of security details for a given product ID. Returns an error if the operation fails.
func (s *Service) GetSecurityByProductID(productID int64) ([]products.SecurityDTO, error) {
	return s.productSvc.GetSecurity(s.ctx, productID)
}

// GetSecurityByOrganisation retrieves a list of security details associated with a specific organization by its ID.
func (s *Service) GetSecurityByOrganisation(id int64) ([]products.SecurityDTO, error) {
	return s.productSvc.GetSecurityByOrg(s.ctx, id)
}

// SyncOrgs synchronizes stale organizations by retrieving them and invoking the sync process for each.
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
