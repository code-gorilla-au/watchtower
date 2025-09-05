package github

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func makeResponse(status int, body string, headers map[string]string) *http.Response {
	h := http.Header{}
	for k, v := range headers {
		h.Set(k, v)
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     h,
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
}

func TestGithubClient(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var mock *fetcherMock
	var client *Client

	group.BeforeEach(func() {
		mock = &fetcherMock{}
		client = &Client{fetch: mock, logger: testLogger()}
	})

	err := group.
		Test("SearchRepos success with rate limits", func(t *testing.T) {
			mock.PostFunc = func(url string, body io.Reader, headers map[string]string) (*http.Response, error) {
				if url != ghUrl {
					t.Fatalf("unexpected URL: %s", url)
				}
				respBody := `{
					"data": {
						"search": {
							"edges": [
								{ "node": { "name": "repo1", "url": "https://example.com/repo1", "owner": {"login": "ow"} } }
							]
						}
					}
				}`
				hds := map[string]string{
					"x-ratelimit-limit":     "5000",
					"x-ratelimit-remaining": "4999",
					"x-ratelimit-used":      "1",
				}
				return makeResponse(200, respBody, hds), nil
			}

			res, err := client.SearchRepos("ow", "topic", "token")
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 1, len(res.Data.Search.Edges))
			odize.AssertEqual(t, "5000", res.RateLimits.Limit)
			odize.AssertEqual(t, "4999", res.RateLimits.Remaining)
			odize.AssertEqual(t, "1", res.RateLimits.Used)
			odize.AssertEqual(t, 1, len(mock.PostCalls()))
		}).
		Test("GetRepoDetails returns joined GraphQL errors", func(t *testing.T) {
			mock.PostFunc = func(url string, body io.Reader, headers map[string]string) (*http.Response, error) {
				respBody := `{
					"errors": [
						{"message": "boom"},
						{"message": "bang"}
					]
				}`
				return makeResponse(200, respBody, map[string]string{}), nil
			}

			_, err := client.GetRepoDetails("ow", "repo", "token")
			odize.AssertError(t, err)
			odize.AssertTrue(t, containsAll(err.Error(), []string{"boom", "bang"}))
		}).
		Test("SearchRepos propagates Post error", func(t *testing.T) {
			expected := errors.New("network down")
			mock.PostFunc = func(url string, body io.Reader, headers map[string]string) (*http.Response, error) {
				return nil, expected
			}
			_, err := client.SearchRepos("ow", "topic", "token")
			odize.AssertTrue(t, errors.Is(err, expected))
		}).
		Test("GetRepoDetails invalid JSON returns error", func(t *testing.T) {
			mock.PostFunc = func(url string, body io.Reader, headers map[string]string) (*http.Response, error) {
				return makeResponse(200, "{invalid", map[string]string{}), nil
			}
			_, err := client.GetRepoDetails("ow", "repo", "token")
			odize.AssertError(t, err)
		}).
		Run()

	odize.AssertNoError(t, err)
}

// containsAll checks all substrings are in s
func containsAll(s string, subs []string) bool {
	for _, sub := range subs {
		if !bytes.Contains([]byte(s), []byte(sub)) {
			return false
		}
	}
	return true
}
