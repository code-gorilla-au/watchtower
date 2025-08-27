package watchtower

import (
	"context"
	"database/sql"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/go-toolbox/github"
	"github.com/code-gorilla-au/go-toolbox/logging"
)

type CreateRepoParams struct {
	Name  string
	Url   string
	Topic string
	Owner string
}

func (r *repoService) CreateRepo(ctx context.Context, params CreateRepoParams) (RepositoryDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Creating repo", "repo", params.Name)

	model, err := r.db.CreateRepo(ctx, database.CreateRepoParams{
		Name:  params.Name,
		Url:   params.Url,
		Topic: params.Topic,
		Owner: params.Owner,
	})
	if err != nil {
		logger.Error("Error creating repo", "error", err)
		return RepositoryDTO{}, err
	}

	return ToRepositoryDTO(model), nil
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

	return toPullRequestDTOs(models), nil
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

	return ToSecurityDTOs(model), nil
}
