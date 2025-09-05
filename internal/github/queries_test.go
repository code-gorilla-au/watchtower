package github

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestQueries(t *testing.T) {
	group := odize.NewGroup(t, nil)

	owner := "acme"
	repo := "super-repo"
	topic := "tooling"

	err := group.
		Test("queryReposByTopic builds expected search query", func(t *testing.T) {
			q := queryReposByTopic(owner, topic)

			odize.AssertTrue(t, containsAll(q, []string{
				"search(",
				"type: REPOSITORY",
				"first: 100",
				"pageInfo",
				"hasNextPage",
				"endCursor",
				"repositoryCount",
				"edges",
				"... on Repository",
				"name",
				"url",
				"owner",
				"login",
			}))

			// Interpolated parameters
			odize.AssertTrue(t, containsAll(q, []string{
				"owner:" + owner,
				"topic:" + topic,
			}))
		}).
		Test("queryGetRepoDetails builds expected repository details query", func(t *testing.T) {
			q := queryGetRepoDetails(owner, repo)

			odize.AssertTrue(t, containsAll(q, []string{
				"repository(owner: \"" + owner + "\"",
				", name: \"" + repo + "\")",
			}))

			odize.AssertTrue(t, containsAll(q, []string{
				"name",
				"url",
				"owner",
				"login",
			}))

			odize.AssertTrue(t, containsAll(q, []string{
				"pullRequests",
				"orderBy: {field: CREATED_AT, direction: ASC}",
				"totalCount",
				"nodes",
				"id",
				"state",
				"title",
				"createdAt",
				"mergedAt",
				"permalink",
				"author",
				"login",
			}))

			odize.AssertTrue(t, containsAll(q, []string{
				"vulnerabilityAlerts",
				"pageInfo",
				"hasNextPage",
				"endCursor",
				"nodes",
				"securityVulnerability",
				"package",
				"name",
				"advisory",
				"severity",
				"firstPatchedVersion",
				"identifier",
				"updatedAt",
			}))
		}).
		Run()

	odize.AssertNoError(t, err)
}
