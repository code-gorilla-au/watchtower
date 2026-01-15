package insights

type PullRequestInsights struct {
	MinDaysToMerge float64 `json:"minDaysToMerge"`
	MaxDaysToMerge float64 `json:"maxDaysToMerge"`
	AvgDaysToMerge float64 `json:"avgDaysToMerge"`
	Merged         int64   `json:"merged"`
	Closed         int64   `json:"closed"`
	Open           int64   `json:"open"`
}

type SecurityInsights struct {
	MinDaysToFix float64 `json:"minDaysToFix"`
	MaxDaysToFix float64 `json:"maxDaysToFix"`
	AvgDaysToFix float64 `json:"avgDaysToFix"`
	Fixed        int64   `json:"fixed"`
	Open         int64   `json:"open"`
}

type Service struct {
	store Store
}
