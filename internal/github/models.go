package github

import "time"

type PageInfo struct {
	EndCursor   *string `json:"endCursor,omitempty"`
	HasNextPage bool    `json:"hasNextPage,omitempty"`
}

type Owner struct {
	Login string `json:"login,omitempty"`
}

type Package struct {
	Name string `json:"name,omitempty"`
}

type Advisory struct {
	Severity AdvisorySeverity `json:"severity,omitempty"`
}

type FirstPatchedVersion struct {
	Identifier string `json:"identifier,omitempty"`
}

type Author struct {
	Login string `json:"login,omitempty"`
}

type RepositoryData struct {
	Repository Repository `json:"repository,omitempty"`
}

type Querier interface {
	HasErrors() bool
	GetErrors() []Errors
	SetLimits(limits RateLimits)
}

type QueryModel struct {
	Errors     []Errors   `json:"errors,omitempty"`
	RateLimits RateLimits `json:"rateLimits,omitempty"`
}

func (q *QueryModel) HasErrors() bool {
	return len(q.Errors) > 0
}

func (q *QueryModel) SetLimits(limits RateLimits) {
	q.RateLimits = limits
}

func (q *QueryModel) GetErrors() []Errors {
	return q.Errors
}

type QueryData[T any] struct {
	Search Search[T] `json:"search,omitempty"`
}

type QuerySearch[T any] struct {
	QueryModel
	Data QueryData[T] `json:"data,omitempty"`
}

type QueryRepository struct {
	QueryModel
	Data RepositoryData `json:"data,omitempty"`
}

type Search[T any] struct {
	PageInfo PageInfo  `json:"pageInfo,omitempty"`
	Edges    []Node[T] `json:"edges,omitempty"`
}

type RootNode[T any] struct {
	PageInfo PageInfo `json:"pageInfo,omitempty"`
	Nodes    []T      `json:"nodes,omitempty"`
}

type Node[T any] struct {
	Node T `json:"node,omitempty"`
}

type Repository struct {
	Url                 string                        `json:"url,omitempty"`
	Name                string                        `json:"name,omitempty"`
	Owner               Owner                         `json:"owner,omitempty"`
	VulnerabilityAlerts RootNode[VulnerabilityAlerts] `json:"vulnerabilityAlerts,omitempty"`
	PullRequests        RootNode[PullRequest]         `json:"pullRequests,omitempty"`
}

type VulnerabilityAlertState string

const (
	Open          VulnerabilityAlertState = "OPEN"
	Fixed         VulnerabilityAlertState = "FIXED"
	Dismissed     VulnerabilityAlertState = "DISMISSED"
	AutoDismissed VulnerabilityAlertState = "AUTO_DISMISSED"
)

type VulnerabilityAlerts struct {
	State                 VulnerabilityAlertState `json:"state,omitempty"`
	ID                    string                  `json:"id,omitempty"`
	Number                int                     `json:"number,omitempty"`
	SecurityVulnerability SecurityVulnerability   `json:"securityVulnerability,omitempty"`
	CreatedAt             time.Time               `json:"createdAt,omitempty"`
	FixedAt               *time.Time              `json:"fixedAt,omitempty"`
}

type AdvisorySeverity string

const (
	Low      AdvisorySeverity = "LOW"
	Moderate AdvisorySeverity = "MODERATE"
	High     AdvisorySeverity = "HIGH"
	Critical AdvisorySeverity = "CRITICAL"
)

type SecurityVulnerability struct {
	Package             Package             `json:"package,omitempty"`
	Advisory            Advisory            `json:"advisory,omitempty"`
	FirstPatchedVersion FirstPatchedVersion `json:"firstPatchedVersion,omitempty"`
	UpdatedAt           time.Time           `json:"updatedAt,omitempty"`
}

type PullRequestState string

const (
	PrOpen   PullRequestState = "OPEN"
	PrClosed PullRequestState = "CLOSED"
	PrMerged PullRequestState = "MERGED"
)

type PullRequest struct {
	ID        string           `json:"id,omitempty"`
	Title     string           `json:"title,omitempty"`
	State     PullRequestState `json:"state,omitempty"`
	CreatedAt time.Time        `json:"createdAt,omitempty"`
	MergedAt  *time.Time       `json:"mergedAt,omitempty"`
	Permalink string           `json:"permalink,omitempty"`
	Author    Author           `json:"author,omitempty"`
}

type RateLimits struct {
	Limit     string `json:"limit,omitempty"`
	Remaining string `json:"remaining,omitempty"`
	Used      string `json:"used,omitempty"`
}

type Errors struct {
	Message string `json:"message,omitempty"`
}
