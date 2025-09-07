package github

import (
	"io"
	"net/http"
)

//go:generate moq -rm -stub -out mocks.gen.go . fetcher ClientGithub

type ClientGithub interface {
	ClientGetter
	ClientSearcher
}

type ClientSearcher interface {
	SearchRepos(owner string, topic string, token string) (QuerySearch[Repository], error)
}

type ClientGetter interface {
	GetRepoDetails(owner, repo string, token string) (QueryRepository, error)
}

var _ ClientGithub = (*Client)(nil)

type fetcher interface {
	Post(url string, body io.Reader, headers map[string]string) (*http.Response, error)
}
