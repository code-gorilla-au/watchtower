package insights

import (
	"context"
	"database/sql"
	"fmt"
	"watchtower/internal/database"
	"watchtower/internal/logging"
)

func New(db Store) *Service {
	return &Service{
		store: db,
	}
}

func (s *Service) GetPullRequestInsightsByOrg(ctx context.Context, orgID int64, filterDays FilterDateDays) (PullRequestInsights, error) {
	logger := logging.FromContext(ctx).With("orgID", orgID, "service", "insights")
	logger.Debug("Fetching pull request insights for org")

	if !validateFilterDateDays(filterDays) {
		logger.Warn("Invalid filter date days", "days", filterDays)
		filterDays = Last30Days
	}

	filterDate := fmt.Sprintf("-%s days", filterDays)

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

func (s *Service) GetSecurityInsightsByOrg(ctx context.Context, orgID int64, filterDays FilterDateDays) (SecurityInsights, error) {
	logger := logging.FromContext(ctx).With("orgID", orgID, "service", "insights")
	logger.Debug("Fetching security insights for org")

	if !validateFilterDateDays(filterDays) {
		logger.Warn("Invalid filter date days", "days", filterDays)
		filterDays = Last30Days
	}

	filterDate := fmt.Sprintf("-%s days", filterDays)

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

func validateFilterDateDays(date FilterDateDays) bool {
	switch date {
	case Last30Days, Last90Days, Last180Days:
		return true
	default:
		return false
	}
}
