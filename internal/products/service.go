package products

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"watchtower/internal/database"
	"watchtower/internal/github"
	"watchtower/internal/logging"
)

type Service struct {
	store ProductStore
}

func New(store ProductStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Create(ctx context.Context, params CreateProductParams) (ProductDTO, error) {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Creating product")

	var tagsNS sql.NullString

	tagJson, err := json.Marshal(params.Tags)
	if err != nil {
		logger.Error("Error marshalling tags", "error", err)

		return ProductDTO{}, err
	}

	tagsNS = sql.NullString{String: string(tagJson), Valid: true}

	prod, err := s.store.CreateProduct(ctx, database.CreateProductParams{
		Name:        params.Name,
		Tags:        tagsNS,
		Description: params.Desc,
	})
	if err != nil {
		logger.Error("Error creating product", "error", err)

		return ProductDTO{}, err
	}

	return toProductDTO(prod), nil
}

func (s *Service) Get(ctx context.Context, id int64) (ProductDTO, error) {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Fetching product by ID")

	prod, err := s.store.GetProductByID(ctx, id)
	if err != nil {
		logger.Error("Error fetching product by ID", "error", err)

		return ProductDTO{}, err
	}

	return toProductDTO(prod), nil
}

func (s *Service) GetByOrg(ctx context.Context, orgID int64) ([]ProductDTO, error) {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Listing products for organisation")

	models, err := s.store.ListProductsByOrganisation(ctx, sql.NullInt64{Int64: orgID, Valid: true})
	if err != nil {
		logger.Error("Error listing products for organisation", "error", err)

		return nil, err
	}

	return toProductDTOs(models), nil
}

func (s *Service) Update(ctx context.Context, params UpdateProductParams) (ProductDTO, error) {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Updating product")

	data, err := json.Marshal(params.Tags)
	if err != nil {
		logger.Error("Error marshalling tags", "error", err)
	}

	var tagsNS sql.NullString
	if params.Tags != nil {
		tagsNS = sql.NullString{String: string(data), Valid: true}
	}

	model, err := s.store.UpdateProduct(ctx, database.UpdateProductParams{
		Name:        params.Name,
		Tags:        tagsNS,
		ID:          params.ID,
		Description: params.Description,
	})
	if err != nil {
		logger.Error("Error updating product", "error", err)

		return ProductDTO{}, err
	}

	return toProductDTO(model), nil
}

func (s *Service) UpdateSyncDateNow(ctx context.Context, id int64) error {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Updating sync date")

	return s.store.UpdateProductSync(ctx, id)
}

func (s *Service) DeleteProduct(ctx context.Context, id int64) error {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Deleting product")

	if err := s.store.DeleteSecurityByProductID(ctx, id); err != nil {
		logger.Error("Error deleting security for product", "error", err)
	}

	if err := s.store.DeletePullRequestsByProductID(ctx, id); err != nil {
		logger.Error("Error deleting pull requests for product", "error", err)
	}

	if err := s.store.DeleteReposByProductID(ctx, id); err != nil {
		logger.Error("Error deleting repos for product", "error", err)
	}

	return s.store.DeleteProduct(ctx, id)
}

func (s *Service) CreateRepo(ctx context.Context, params CreateRepoParams) error {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Creating repo", "repo", params.Name)

	_, err := s.store.CreateRepo(ctx, database.CreateRepoParams{
		Name:  params.Name,
		Url:   params.Url,
		Topic: params.Topic,
		Owner: params.Owner,
	})
	if err != nil {
		logger.Error("Error creating repo", "error", err)
		return err
	}

	return nil
}

func (s *Service) UpdateRepo(ctx context.Context, params UpdateRepoParams) error {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Updating repo", "repo", params.Name)

	_, err := s.store.UpdateRepo(ctx, database.UpdateRepoParams{
		Name:  params.Name,
		Url:   params.Url,
		Topic: params.Topic,
		Owner: params.Owner,
		ID:    params.ID,
	})
	if err != nil {
		logger.Error("Error updating repo", "error", err)
		return err
	}

	return nil
}

func (s *Service) UpsertRepo(ctx context.Context, params CreateRepoParams) error {
	logger := logging.FromContext(ctx).With("service", "products")

	createErr := s.CreateRepo(ctx, params)
	if createErr == nil {
		return nil
	}

	if !strings.Contains(createErr.Error(), "constraint failed: UNIQUE constraint failed") {
		logger.Error("Error creating repo", "error", createErr)
		return createErr
	}

	model, getErr := s.store.GetRepoByName(ctx, params.Name)
	if getErr != nil {
		logger.Error("Error fetching repo", "error", getErr)
		return getErr
	}

	if updateErr := s.UpdateRepo(ctx, UpdateRepoParams{
		ID:    model.ID,
		Name:  params.Name,
		Url:   params.Url,
		Topic: params.Topic,
		Owner: params.Owner,
	}); updateErr != nil {
		logger.Error("Error updating repo", "error", updateErr)
		return updateErr
	}

	return nil
}
func (s *Service) GetRepos(ctx context.Context, productID int64) ([]RepositoryDTO, error) {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Fetching repos for product")

	repos, err := s.store.GetReposByProductID(ctx, productID)
	if err != nil {
		logger.Error("Error fetching repos for product", "error", err)

		return nil, err
	}

	result := make([]RepositoryDTO, 0, len(repos))
	for _, r := range repos {
		result = append(result, toRepositoryDTO(r))
	}

	return result, nil
}

func (s *Service) GetPullRequests(ctx context.Context, productID int64) ([]PullRequestDTO, error) {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Fetching pull requests for product")

	models, err := s.store.GetPullRequestByProductIDAndState(ctx, database.GetPullRequestByProductIDAndStateParams{
		ID:    productID,
		State: string(github.PrOpen),
	})
	if err != nil {
		logger.Error("Error fetching pull requests for product", "error", err)

		return nil, err
	}

	return toPullRequestDTOs(models), nil
}

func (s *Service) GetPullRequestByOrg(ctx context.Context, orgID int64) ([]PullRequestDTO, error) {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Fetching pull requests for organisation", "org", orgID)

	models, err := s.store.GetPullRequestsByOrganisationAndState(ctx, database.GetPullRequestsByOrganisationAndStateParams{
		OrganisationID: sql.NullInt64{
			Int64: orgID,
			Valid: true,
		},
		State: "OPEN",
	})
	if err != nil {
		logger.Error("Error fetching pull requests for product", "error", err)

		return nil, err
	}

	logger.Debug("Found", "count", len(models))

	return orgToPullRequestDTOs(models), nil
}

func (s *Service) BulkUpsertRepos(ctx context.Context, paramsList []CreateRepoParams) error {
	for _, params := range paramsList {
		err := s.UpsertRepo(ctx, params)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreatePullRequest creates a new pull request entry in the database using the provided parameters.
func (s *Service) CreatePullRequest(ctx context.Context, params CreatePRParams) error {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Creating pull request")

	var mergedAt sql.NullInt64
	if params.MergedAt != nil {
		mergedAt.Valid = true
		mergedAt.Int64 = params.MergedAt.Unix()
	}

	_, err := s.store.CreatePullRequest(ctx, database.CreatePullRequestParams{
		ExternalID:     params.ExternalID,
		Title:          params.Title,
		RepositoryName: params.RepositoryName,
		Url:            params.Url,
		State:          params.State,
		Author:         params.Author,
		MergedAt:       mergedAt,
		CreatedAt:      params.CreatedAt.Unix(),
	})

	if err != nil {
		logger.Error("Error creating pull request", "error", err)
		return err
	}

	return nil
}

func (s *Service) UpdatePullRequest(ctx context.Context, params UpdatePRParams) error {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("Updating pull request")

	var mergedAt sql.NullInt64
	if params.MergedAt != nil {
		mergedAt.Valid = true
		mergedAt.Int64 = params.MergedAt.Unix()
	}

	_, err := s.store.UpdatePullRequest(ctx, database.UpdatePullRequestParams{
		Title:          params.Title,
		RepositoryName: params.RepositoryName,
		Url:            params.Url,
		State:          params.State,
		Author:         params.Author,
		MergedAt:       mergedAt,
		ID:             params.ID,
	})

	if err != nil {
		logger.Error("Error updating pull request", "error", err)
		return err
	}

	return nil
}

func (s *Service) UpsertPullRequest(ctx context.Context, params CreatePRParams) error {
	logger := logging.FromContext(ctx).With("service", "products")

	createErr := s.CreatePullRequest(ctx, params)
	if createErr == nil {
		return nil
	}

	if !strings.Contains(createErr.Error(), "constraint failed: UNIQUE constraint failed") {
		logger.Error("Error creating pull request", "error", createErr)
		return createErr
	}

	pr, getErr := s.store.GetPullRequestByExternalID(ctx, params.ExternalID)
	if getErr != nil {
		logger.Error("Error fetching pull request", "error", getErr)
		return getErr
	}

	return s.UpdatePullRequest(ctx, UpdatePRParams{
		ID:             pr.ID,
		ExternalID:     pr.ExternalID,
		Title:          params.Title,
		RepositoryName: params.RepositoryName,
		Url:            params.Url,
		State:          params.State,
		Author:         params.Author,
		MergedAt:       params.MergedAt,
	})
}

func (s *Service) BulkCreatePullRequest(ctx context.Context, paramsList []CreatePRParams) error {
	logger := logging.FromContext(ctx)

	for _, params := range paramsList {

		if err := s.UpsertPullRequest(ctx, params); err != nil {
			logger.Error("Error creating pull request", "error", err)
			return err
		}

	}

	return nil
}

func (s *Service) GetSecurity(ctx context.Context, productID int64) ([]SecurityDTO, error) {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("getting security by product id")

	model, err := s.store.GetSecurityByProductIDAndState(ctx, database.GetSecurityByProductIDAndStateParams{
		ID:    productID,
		State: "OPEN",
	})
	if err != nil {
		logger.Error("Error fetching security by product id", "error", err)

		return []SecurityDTO{}, err
	}

	return toSecurityDTOs(model), nil
}

func (s *Service) GetSecurityByOrg(ctx context.Context, orgID int64) ([]SecurityDTO, error) {
	logger := logging.FromContext(ctx).With("service", "products")

	logger.Debug("getting security by organisation", "org", orgID)

	model, err := s.store.GetSecurityByOrganisationAndState(ctx, database.GetSecurityByOrganisationAndStateParams{
		OrganisationID: sql.NullInt64{Int64: orgID, Valid: true},
		State:          "OPEN",
	})
	if err != nil {
		logger.Error("Error fetching security by organisation", "error", err)

		return []SecurityDTO{}, err
	}

	return orgToSecurityDTOs(model), nil
}

func (s *Service) BulkCreateSecurity(ctx context.Context, paramsList []CreateSecurityParams) error {
	logger := logging.FromContext(ctx)

	for _, params := range paramsList {
		fixedAt := sql.NullInt64{}
		if params.FixedAt != nil {
			fixedAt.Int64 = params.FixedAt.Unix()
			fixedAt.Valid = true
		}

		_, err := s.store.CreateSecurity(ctx, database.CreateSecurityParams{
			ExternalID:     params.ExternalID,
			RepositoryName: params.RepositoryName,
			PackageName:    params.PackageName,
			State:          params.State,
			Severity:       params.Severity,
			PatchedVersion: params.PatchedVersion,
			FixedAt:        fixedAt,
		})
		if err != nil {
			logger.Error("Error creating security", "error", err)
			return err
		}
	}

	return nil
}

func (s *Service) BulkInsertRepos(ctx context.Context, reposList []github.Node[github.Repository], tag string) error {
	params := toCreateRepoFromGithub(reposList, tag)

	return s.BulkUpsertRepos(ctx, params)
}

func (s *Service) BulkInsertRepoDetails(ctx context.Context, repoDetails github.QueryRepository) error {
	logger := logging.FromContext(ctx)

	prs := toCreatePRsFromGithubRepos(repoDetails.Data.Repository.PullRequests, repoDetails.Data.Repository.Name)
	secs := toSecParamsFromGithubVulnerabilities(repoDetails.Data.Repository.VulnerabilityAlerts, repoDetails.Data.Repository.Name)

	if err := s.BulkCreatePullRequest(ctx, prs); err != nil {
		logger.Error("Error creating pull requests", "error", err)
		return err
	}

	if err := s.BulkCreateSecurity(ctx, secs); err != nil {
		logger.Error("Error creating security", "error", err)
		return err
	}

	return nil
}
