package watchtower

import (
	"database/sql"
	"encoding/json"
	"strings"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/go-toolbox/github"
	"github.com/code-gorilla-au/go-toolbox/logging"
)

// CreateProduct creates a new product and associates it with an organisation
func (s *Service) CreateProduct(name string, description string, tags []string, organisationID int64) (ProductDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Creating product")

	var tagsNS sql.NullString

	tagJson, err := json.Marshal(tags)
	if err != nil {
		logger.Error("Error marshalling tags", err)
		return ProductDTO{}, err
	}

	tagsNS = sql.NullString{String: string(tagJson), Valid: true}

	prod, err := s.db.CreateProduct(s.ctx, database.CreateProductParams{
		Name:        name,
		Tags:        tagsNS,
		Description: description,
	})
	if err != nil {
		logger.Error("Error creating product", err)
		return ProductDTO{}, err
	}

	err = s.db.AddProductToOrganisation(s.ctx, database.AddProductToOrganisationParams{
		ProductID:      sql.NullInt64{Int64: prod.ID, Valid: true},
		OrganisationID: sql.NullInt64{Int64: organisationID, Valid: true},
	})
	if err != nil {
		logger.Error("Error linking product to organisation", err)
		return ProductDTO{}, err
	}

	return ToProductDTO(prod), nil
}

// GetProductByID fetches a product by its ID
func (s *Service) GetProductByID(id int64) (ProductDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Fetching product by ID")

	prod, err := s.db.GetProductByID(s.ctx, id)
	if err != nil {
		logger.Error("Error fetching product by ID", err)
		return ProductDTO{}, err
	}

	return ToProductDTO(prod), nil
}

// GetAllProductsForOrganisation lists products linked to the given organisation
func (s *Service) GetAllProductsForOrganisation(organisationID int64) ([]ProductDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Listing products for organisation")

	models, err := s.db.ListProductsByOrganisation(s.ctx, sql.NullInt64{Int64: organisationID, Valid: true})
	if err != nil {
		logger.Error("Error listing products for organisation", err)
		return nil, err
	}

	result := make([]ProductDTO, 0, len(models))
	for _, m := range models {
		result = append(result, ToProductDTO(m))
	}

	return result, nil
}

// UpdateProduct updates a product and returns the updated entity
func (s *Service) UpdateProduct(id int64, name string, tags *string) (ProductDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Updating product")

	var tagsNS sql.NullString
	if tags != nil {
		tagsNS = sql.NullString{String: *tags, Valid: true}
	}

	err := s.db.UpdateProduct(s.ctx, database.UpdateProductParams{
		Name: name,
		Tags: tagsNS,
		ID:   id,
	})
	if err != nil {
		logger.Error("Error updating product", err)
		return ProductDTO{}, err
	}

	prod, err := s.db.GetProductByID(s.ctx, id)
	if err != nil {
		logger.Error("Error fetching updated product", err)
		return ProductDTO{}, err
	}

	return ToProductDTO(prod), nil
}

func (s *Service) DeleteProduct(id int64) error {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Deleting product")
	if err := s.db.DeleteProduct(s.ctx, id); err != nil {
		logger.Error("Error deleting product", err)
	}

	if err := s.deleteReposByProductID(id); err != nil {
		logger.Error("Error deleting repos for product", err)
	}

	return nil
}

func (s *Service) GetProductRepos(id int64) ([]RepositoryDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Fetching repos for product")

	repos, err := s.db.GetReposByProductID(s.ctx, id)
	if err != nil {
		logger.Error("Error fetching repos for product", err)
		return nil, err
	}

	result := make([]RepositoryDTO, 0, len(repos))
	for _, r := range repos {
		result = append(result, ToRepositoryDTO(r))
	}

	return result, nil
}

func (s *Service) GetProductPullRequests(id int64) ([]PullRequestDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Fetching pull requests for product")
	models, err := s.db.GetPullRequestByProductIDAndState(s.ctx, database.GetPullRequestByProductIDAndStateParams{
		ID:    id,
		State: string(github.PrOpen),
	})
	if err != nil {
		logger.Error("Error fetching pull requests for product", err)
		return nil, err
	}

	return toPullRequestDTOs(models), nil
}

func (s *Service) GetPullRequestByOrganisation(id int64) ([]PullRequestDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Fetching pull requests for product")
	models, err := s.db.GetPullRequestByProductIDAndState(s.ctx, database.GetPullRequestByProductIDAndStateParams{
		ID:    id,
		State: string(github.PrOpen),
	})
	if err != nil {
		logger.Error("Error fetching pull requests for product", err)
		return nil, err
	}

	return toPullRequestDTOs(models), nil
}

func (s *Service) SyncProduct(id int64) error {
	logger := logging.FromContext(s.ctx)

	product, err := s.GetProductByID(id)
	if err != nil {
		logger.Error("Error fetching product", err)
		return err
	}

	org, err := s.db.GetOrganisationForProduct(s.ctx, sql.NullInt64{Int64: product.ID, Valid: true})
	if err != nil {
		logger.Error("Error fetching organisation for product", err)
		return err
	}

	for _, tag := range product.Tags {
		if err = s.syncByTag(tag, org.Namespace, org.Token); err != nil {
			logger.Error("Error syncing repos", "error", err)
			return err
		}

	}

	return nil
}

func (s *Service) syncByTag(tag string, owner string, ghToken string) error {
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
		mergedAt := sql.NullInt64{}

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

func (s *Service) deleteReposByProductID(id int64) error {
	logger := logging.FromContext(s.ctx)
	logger.Debug("Deleting repos for product")
	return s.db.DeleteReposByProductID(s.ctx, id)
}
