package insights

import (
	"context"
	"database/sql"
	"watchtower/internal/database"
	"watchtower/internal/logging"
)

func New(db Store) *Service {
	return &Service{
		store: db,
	}
}

func (s *Service) GetPullRequestInsightsByOrg(ctx context.Context, orgID int64, filterDate string) (PullRequestInsights, error) {
	logger := logging.FromContext(ctx).With("orgID", orgID, "service", "insights")
	logger.Debug("Fetching pull request insights for org")

	row, err := s.store.GetPullRequestInsightsByOrg(ctx, database.GetPullRequestInsightsByOrgParams{
		OrganisationID: sql.NullInt64{Int64: orgID, Valid: true},
		Strftime:       filterDate,
		Strftime_2:     filterDate,
		Strftime_3:     filterDate,
	})
	if err != nil {
		logger.Error("Error fetching pull request insights", "error", err)
		return PullRequestInsights{}, err
	}

	return fromPullRequestInsightsModel(row), nil
}

func (s *Service) GetSecurityInsightsByOrg(ctx context.Context, orgID int64, filterDate string) (SecurityInsights, error) {
	logger := logging.FromContext(ctx).With("orgID", orgID, "service", "insights")
	logger.Debug("Fetching security insights for org")

	row, err := s.store.GetSecuritiesInsightsByOrg(ctx, database.GetSecuritiesInsightsByOrgParams{
		OrganisationID: sql.NullInt64{Int64: orgID, Valid: true},
		Strftime:       filterDate,
		Strftime_2:     filterDate,
		Strftime_3:     filterDate,
	})
	if err != nil {
		logger.Error("Error fetching security insights", "error", err)
		return SecurityInsights{}, err
	}

	return fromSecurityInsightsModel(row), nil
}
