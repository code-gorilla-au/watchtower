// Package github provides simple queries for retrieving data from repositories, pull requests and security vulnerabilities
package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/code-gorilla-au/fetch"
)

const ghUrl = "https://api.github.com/graphql"

type Client struct {
	fetch  fetcher
	logger *slog.Logger
}

// New initialises and returns a new Client instance with the provided token and logger.
func New(logger *slog.Logger) *Client {
	return &Client{
		fetch: fetch.New(
			fetch.WithOpts(),
		),
		logger: logger,
	}
}

// SearchRepos searches for repositories by owner and topic
func (c *Client) SearchRepos(owner string, topic string, token string) (QuerySearch[Repository], error) {
	log := c.logger.With("owner", owner, "topic", topic)
	log.Info("Searching repos by owner and topic")

	var result QuerySearch[Repository]

	err := executeQuery(c.fetch, queryReposByTopic(owner, topic), &result, token)

	return result, err
}

// GetRepoDetails fetches repository details for a specified owner and repository name and returns the result or an error.
func (c *Client) GetRepoDetails(owner, repo string, token string) (QueryRepository, error) {
	log := c.logger.With("owner", owner, "repo", repo)
	log.Info("Getting repo details")

	var result QueryRepository

	err := executeQuery(c.fetch, queryGetRepoDetails(owner, repo), &result, token)

	return result, err
}

// extractRateLimits extracts rate limit information from the HTTP response headers and returns it as a RateLimits struct.
func extractRateLimits(resp *http.Response) RateLimits {
	return RateLimits{
		Limit:     resp.Header.Get("x-ratelimit-limit"),
		Remaining: resp.Header.Get("x-ratelimit-remaining"),
		Used:      resp.Header.Get("x-ratelimit-used"),
	}
}

func executeQuery(client fetcher, query string, dest Querier, token string) error {

	q, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(q)

	resp, err := client.Post(ghUrl, buf, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if err = json.NewDecoder(resp.Body).Decode(&dest); err != nil {
		return err
	}

	rateLimits := extractRateLimits(resp)
	dest.SetLimits(rateLimits)

	if dest.HasErrors() {
		var errs []error
		for _, e := range dest.GetErrors() {
			errs = append(errs, errors.New(e.Message))
		}

		return errors.Join(errs...)
	}

	return nil
}
