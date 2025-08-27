package watchtower

import (
	"database/sql"
	"strings"
	"time"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/go-toolbox/github"
	"github.com/code-gorilla-au/go-toolbox/logging"
)

// CreateProduct creates a new product and associates it with an organisation.
func (s *Service) CreateProduct(name string, description string, tags []string, organisationID int64) (ProductDTO, error) {

	prod, err := s.productSvc.Create(s.ctx, CreateProductParams{
		Name: name,
		Tags: tags,
		Desc: description,
	})
	if err != nil {
		return ProductDTO{}, err
	}

	err = s.orgSvc.AssociateProductToOrg(s.ctx, organisationID, prod.ID)
	if err != nil {
		return ProductDTO{}, err
	}

	return prod, nil
}

// GetProductByID fetches a product by its ID.
func (s *Service) GetProductByID(id int64) (ProductDTO, error) {
	return s.productSvc.Get(s.ctx, id)
}

// GetAllProductsForOrganisation lists products linked to the given organisation.
func (s *Service) GetAllProductsForOrganisation(organisationID int64) ([]ProductDTO, error) {
	return s.productSvc.GetByOrg(s.ctx, organisationID)
}

// UpdateProduct updates a product and returns the updated entity.
func (s *Service) UpdateProduct(id int64, name string, tags []string) (ProductDTO, error) {
	return s.productSvc.Update(s.ctx, UpdateProductParams{
		ID:   id,
		Name: name,
		Tags: tags,
	})
}

func (s *Service) DeleteProduct(id int64) error {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Deleting product")

	if err := s.deleteReposByProductID(id); err != nil {
		logger.Error("Error deleting repos for product", "error", err)
	}

	if err := s.deleteSecurityByProductID(id); err != nil {
		logger.Error("Error deleting security for product", "error", err)
	}

	if err := s.deletePullRequestsByProductID(id); err != nil {
		logger.Error("Error deleting pull requests for product", "error", err)
	}

	if err := s.db.DeleteProduct(s.ctx, id); err != nil {
		logger.Error("Error deleting product", "error", err)

		return err
	}

	return nil
}

func (s *Service) GetProductRepos(id int64) ([]RepositoryDTO, error) {
	return s.productSvc.GetRepos(s.ctx, id)
}

func (s *Service) GetProductPullRequests(id int64) ([]PullRequestDTO, error) {
	return s.productSvc.GetPullRequests(s.ctx, id)
}

func (s *Service) GetPullRequestByOrganisation(id int64) ([]PullRequestDTO, error) {
	return s.productSvc.GetPullRequestByOrg(s.ctx, id)
}

func (s *Service) deletePullRequestsByProductID(id int64) error {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Deleting PRs for product")

	return s.db.DeletePullRequestsByProductID(s.ctx, id)
}

func (s *Service) GetSecurityByProductID(productID int64) ([]SecurityDTO, error) {
	return s.productSvc.GetSecurity(s.ctx, productID)
}

func (s *Service) GetSecurityByOrganisation(id int64) ([]SecurityDTO, error) {
	return s.productSvc.GetSecurityByOrg(s.ctx, id)
}

func (s *Service) deleteSecurityByProductID(id int64) error {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Deleting security for product")

	return s.db.DeleteSecurityByProductID(s.ctx, id)
}

func (s *Service) SyncOrgs() error {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Syncing orgs")

	fiveMinutesAgo := time.Now().Add(-5 * time.Minute).Unix()

	orgs, err := s.db.ListOrgsOlderThanUpdatedAt(s.ctx, fiveMinutesAgo)
	if err != nil {
		logger.Error("Error fetching orgs", "error", err)

		return err
	}

	logger.Info("syncing number of orgs", "count", len(orgs))

	for _, org := range orgs {
		if err = s.SyncOrg(org.ID); err != nil {
			logger.Error("Error syncing org", "error", err)

			continue
		}
	}

	return nil
}

func (s *Service) SyncOrg(orgId int64) error {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Syncing org", "org", orgId)

	products, err := s.GetAllProductsForOrganisation(orgId)
	if err != nil {
		logger.Error("Error fetching products for org", "error", err)

		return err
	}

	if len(products) == 0 {
		logger.Debug("No products found for org", "org", orgId)
		return nil
	}

	org, err := s.db.GetOrganisationForProduct(s.ctx, sql.NullInt64{Int64: products[0].ID, Valid: true})
	if err != nil {
		logger.Error("Error fetching organisation for product", "error", err)

		return err
	}

	for _, p := range products {
		if err = s.syncProductFromGithub(p, org); err != nil {
			logger.Error("Error syncing product", "error", err)

			return err
		}
	}

	if err = s.db.UpdateOrganisationSync(s.ctx, org.ID); err != nil {
		logger.Error("Error updating organisation sync", "error", err)

		return err
	}

	return nil
}

func (s *Service) SyncProduct(id int64) error {
	logger := logging.FromContext(s.ctx)

	product, err := s.GetProductByID(id)
	if err != nil {
		logger.Error("Error fetching product", "error", err)

		return err
	}

	org, err := s.db.GetOrganisationForProduct(s.ctx, sql.NullInt64{Int64: product.ID, Valid: true})
	if err != nil {
		logger.Error("Error fetching organisation for product", "error", err)

		return err
	}

	return s.syncProductFromGithub(product, org)
}

func (s *Service) syncProductFromGithub(product ProductDTO, org database.Organisation) error {
	logger := logging.FromContext(s.ctx)

	for _, tag := range product.Tags {
		if err := s.syncRepoDataByTag(tag, org.Namespace, org.Token); err != nil {
			logger.Error("Error syncing repos", "error", err)

			return err
		}
	}

	if err := s.db.UpdateProductSync(s.ctx, product.ID); err != nil {
		logger.Error("Error updating product sync", "error", err)

		return err
	}

	return nil
}

func (s *Service) syncRepoDataByTag(tag string, owner string, ghToken string) error {
	logger := logging.FromContext(s.ctx)

	logger.Debug("Searching for repo with tag", "tag", tag)

	repos, apiErr := s.ghClient.SearchRepos(owner, strings.TrimSpace(tag), ghToken)
	if apiErr != nil {
		logger.Error("Error searching for repos", "error", apiErr)

		return apiErr
	}

	if err := s.bulkInsertRepos(repos.Data.Search.Edges, tag); err != nil {
		logger.Error("Error bulk inserting repos", "error", err)

		return err
	}

	for _, repo := range repos.Data.Search.Edges {
		dd, err := s.ghClient.GetRepoDetails(owner, repo.Node.Name, ghToken)
		if err != nil {
			logger.Error("Error getting repo details", "repo", repo.Node.Name, "error", err)
		}

		if err = s.bulkInsertPullRequests(dd.Data.Repository.PullRequests, repo.Node.Name); err != nil {
			logger.Error("Error bulk inserting pull requests", "error", err)

			return err
		}

		if err = s.bulkInsertSecurity(repo.Node.VulnerabilityAlerts, repo.Node.Name); err != nil {
			logger.Error("Error bulk inserting security", "error", err)

			return err
		}
	}

	return nil
}

func (s *Service) bulkInsertRepos(repos []github.Node[github.Repository], tag string) error {
	logger := logging.FromContext(s.ctx)

	for _, repo := range repos {
		_, err := s.db.CreateRepo(s.ctx, database.CreateRepoParams{
			Name:  repo.Node.Name,
			Url:   repo.Node.Url,
			Topic: tag,
			Owner: repo.Node.Owner.Login,
		})
		if err != nil {
			logger.Error("Error creating repo", "error", err)

			return err
		}
	}

	return nil
}

func (s *Service) bulkInsertPullRequests(prs github.RootNode[github.PullRequest], repoName string) error {
	logger := logging.FromContext(s.ctx)

	for _, pr := range prs.Nodes {
		mergedAt := sql.NullInt64{
			Valid: false,
			Int64: 0,
		}

		if pr.MergedAt != nil {
			mergedAt.Int64 = pr.MergedAt.Unix()
			mergedAt.Valid = true
		}

		_, err := s.db.CreatePullRequest(s.ctx, database.CreatePullRequestParams{
			ExternalID:     pr.ID,
			Title:          pr.Title,
			RepositoryName: repoName,
			Url:            pr.Permalink,
			State:          string(pr.State),
			Author:         pr.Author.Login,
			MergedAt:       mergedAt,
		})
		if err != nil {
			logger.Error("Error creating pull request", "error", err)

			return err
		}
	}

	return nil
}

func (s *Service) bulkInsertSecurity(secs github.RootNode[github.VulnerabilityAlerts], repoName string) error {
	logger := logging.FromContext(s.ctx)

	for _, sec := range secs.Nodes {
		_, err := s.db.CreateSecurity(s.ctx, database.CreateSecurityParams{
			ExternalID:     sec.ID,
			RepositoryName: repoName,
			PackageName:    sec.SecurityVulnerability.Package.Name,
			State:          string(sec.State),
			Severity:       string(sec.SecurityVulnerability.Advisory.Severity),
			PatchedVersion: sec.SecurityVulnerability.FirstPatchedVersion.Identifier,
		})
		if err != nil {
			logger.Error("Error creating security", "error", err)

			return err
		}
	}

	return nil
}

func (s *Service) deleteReposByProductID(id int64) error {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Deleting repos for product")

	return s.db.DeleteReposByProductID(s.ctx, id)
}
