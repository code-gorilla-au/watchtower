package insights

import (
	"watchtower/internal/database"
)

func fromPullRequestInsightsModel(m database.GetPullRequestInsightsByOrgRow) PullRequestInsights {
	return PullRequestInsights{
		MinDaysToMerge: m.MinDaysToMerge,
		MaxDaysToMerge: m.MaxDaysToMerge,
		AvgDaysToMerge: m.AvgDaysToMerge,
		Merged:         m.Merged,
		Closed:         m.Closed,
		Open:           m.Open,
	}
}

func fromSecurityInsightsModel(m database.GetSecuritiesInsightsByOrgRow) SecurityInsights {
	return SecurityInsights{
		MinDaysToFix: m.MinDaysToFix,
		MaxDaysToFix: m.MaxDaysToFix,
		AvgDaysToFix: m.AvgDaysToFix,
		Fixed:        m.Fixed,
		Open:         m.Open,
	}
}
