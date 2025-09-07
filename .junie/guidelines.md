

## SQLC queries

When writing SQL queries ensure you annotate your queries

following are examples of correct annotations

```sql
-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = ? LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateAuthor :exec
UPDATE authors
set name = ?,
bio = ?
WHERE id = ?;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;
```

## Frontend

## Golang

- use Odize as the test framework
- Do not add useless comments

writing test example

```golang
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


func TestService_GetAllOrganisations(t *testing.T) {
    group := odize.NewGroup(t, nil)
    
    var s *Service
    
    ctx := context.Background()
    
    group.BeforeEach(func() {
        s = NewService(ctx, _testDB, _testTxnDB)
    })
    
    err := group.
        Test("should return all existing organisations", func(t *testing.T) {
        initialOrgs, err := s.GetAllOrganisations()
        odize.AssertNoError(t, err)
        initialCount := len(initialOrgs)
        odize.AssertTrue(t, initialCount > 0)
        }).
    Run()
    odize.AssertNoError(t, err)
}
```