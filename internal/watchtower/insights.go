package watchtower

import "watchtower/internal/insights"

// GetPullRequestInsightsByOrg returns the pull request insights for the specified organisation and filter date.
func (s *Service) GetPullRequestInsightsByOrg(orgID int64, filterDate string) (insights.PullRequestInsights, error) {
	return s.insightSvc.GetPullRequestInsightsByOrg(s.ctx, orgID, filterDate)
}

// GetSecurityInsightsByOrg returns the security insights for the specified organisation and filter date.
func (s *Service) GetSecurityInsightsByOrg(orgID int64, filterDate string) (insights.SecurityInsights, error) {
	return s.insightSvc.GetSecurityInsightsByOrg(s.ctx, orgID, filterDate)
}
