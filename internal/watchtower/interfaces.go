package watchtower

import "github.com/code-gorilla-au/go-toolbox/github"

//go:generate moq -rm -stub -out mocks.gen.go . ghClient

type ghClient interface {
	SearchRepos(owner string, topic string, token string) (github.QuerySearch[github.Repository], error)
	GetRepoDetails(owner string, repo string, token string) (github.QueryRepository, error)
}
