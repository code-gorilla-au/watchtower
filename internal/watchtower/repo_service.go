package watchtower

import (
	"context"
	"database/sql"
	"time"
	"watchtower/internal/database"

	"watchtower/internal/github"
	"watchtower/internal/logging"
)

type CreateRepoParams struct {
	Name  string
	Url   string
	Topic string
	Owner string
}

func (r *repoService) CreateRepo(ctx context.Context, params CreateRepoParams) error {
	logger := logging.FromContext(ctx)
	logger.Debug("Creating repo", "repo", params.Name)

	_, err := r.db.CreateRepo(ctx, database.CreateRepoParams{
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

func (r *repoService) GetRepos(ctx context.Context, productID int64) ([]RepositoryDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Fetching repos for product")

	repos, err := r.db.GetReposByProductID(ctx, productID)
	if err != nil {
		logger.Error("Error fetching repos for product", "error", err)

		return nil, err
	}

	result := make([]RepositoryDTO, 0, len(repos))
	for _, r := range repos {
		result = append(result, ToRepositoryDTO(r))
	}

	return result, nil
}

func (r *repoService) GetPullRequests(ctx context.Context, productID int64) ([]PullRequestDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Fetching pull requests for product")

	models, err := r.db.GetPullRequestByProductIDAndState(ctx, database.GetPullRequestByProductIDAndStateParams{
		ID:    productID,
		State: string(github.PrOpen),
	})
	if err != nil {
		logger.Error("Error fetching pull requests for product", "error", err)

		return nil, err
	}

	return toPullRequestDTOs(models), nil
}

func (r *repoService) GetPullRequestByOrg(ctx context.Context, orgID int64) ([]PullRequestDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Fetching pull requests for organisation", "org", orgID)

	models, err := r.db.GetPullRequestsByOrganisationAndState(ctx, database.GetPullRequestsByOrganisationAndStateParams{
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

func (r *repoService) BulkCreateRepos(ctx context.Context, paramsList []CreateRepoParams) error {
	for _, params := range paramsList {
		err := r.CreateRepo(ctx, params)
		if err != nil {
			return err
		}
	}

	return nil
}

type CreatePRParams struct {
	ExternalID     string
	Title          string
	RepositoryName string
	Url            string
	State          string
	Author         string
	MergedAt       *time.Time
}

func (r *repoService) BulkCreatePullRequest(ctx context.Context, paramsList []CreatePRParams) error {
	logger := logging.FromContext(ctx)

	for _, params := range paramsList {
		mergedAt := sql.NullInt64{
			Valid: false,
			Int64: 0,
		}

		if params.MergedAt != nil {
			mergedAt.Int64 = params.MergedAt.Unix()
			mergedAt.Valid = true
		}
		_, err := r.db.CreatePullRequest(ctx, database.CreatePullRequestParams{
			ExternalID:     params.ExternalID,
			Title:          params.Title,
			RepositoryName: params.RepositoryName,
			Url:            params.Url,
			State:          params.State,
			Author:         params.Author,
			MergedAt:       mergedAt,
		})
		if err != nil {
			logger.Error("Error creating pull request", "error", err)

			return err
		}
	}

	return nil
}

func (r *repoService) GetSecurity(ctx context.Context, productID int64) ([]SecurityDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("getting security by product id")

	model, err := r.db.GetSecurityByProductIDAndState(ctx, database.GetSecurityByProductIDAndStateParams{
		ID:    productID,
		State: "OPEN",
	})
	if err != nil {
		logger.Error("Error fetching security by product id", "error", err)

		return []SecurityDTO{}, err
	}

	return ToSecurityDTOs(model), nil
}

func (r *repoService) GetSecurityByOrg(ctx context.Context, orgID int64) ([]SecurityDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("getting security by organisation", "org", orgID)

	model, err := r.db.GetSecurityByOrganisationAndState(ctx, database.GetSecurityByOrganisationAndStateParams{
		OrganisationID: sql.NullInt64{Int64: orgID, Valid: true},
		State:          "OPEN",
	})
	if err != nil {
		logger.Error("Error fetching security by organisation", "error", err)

		return []SecurityDTO{}, err
	}

	return OrgToSecurityDTOs(model), nil
}

type CreateSecurityParams struct {
	ExternalID     string
	RepositoryName string
	PackageName    string
	State          string
	Severity       string
	PatchedVersion string
}

func (r *repoService) BulkCreateSecurity(ctx context.Context, paramsList []CreateSecurityParams) error {
	logger := logging.FromContext(ctx)

	for _, params := range paramsList {
		_, err := r.db.CreateSecurity(ctx, database.CreateSecurityParams{
			ExternalID:     params.ExternalID,
			RepositoryName: params.RepositoryName,
			PackageName:    params.PackageName,
			State:          params.State,
			Severity:       params.Severity,
			PatchedVersion: params.PatchedVersion,
		})
		if err != nil {
			logger.Error("Error creating security", "error", err)
			return err
		}
	}

	return nil
}
