package github

import (
	"fmt"
)

func queryReposByTopic(owner string, topic string) string {
	return fmt.Sprintf(`
{
  search(query: "owner:%s topic:%s", type: REPOSITORY, first: 100) {
    pageInfo {
      hasNextPage
      endCursor
    }
    repositoryCount
    edges {
      node {
        ... on Repository {
          name
          url
          owner {
            login
          }
        }
      }
    }
  }
}
  `, owner, topic)
}

func queryGetRepoDetails(owner, repo string) string {
	return fmt.Sprintf(`
{
  repository(owner: "%s", name: "%s") {
    name
    url
    owner {
      login
    }
    pullRequests(last: 100, orderBy: {field: CREATED_AT, direction: ASC}) {
      totalCount
      nodes {
        id
        state
        title
        createdAt
        closedAt
        mergedAt
        permalink
        author {
          login
        }
      }
    }
    vulnerabilityAlerts(last: 100) {
      pageInfo {
        hasNextPage
        endCursor
      }
      nodes {
        state
        id
        number
        createdAt
        fixedAt
        securityVulnerability {
          package {
            name
          }
          advisory {
            severity
          }
          firstPatchedVersion {
            identifier
          }
          updatedAt
        }
        dependabotUpdate {
          pullRequest {
            id
            state
            createdAt
            closedAt
            mergedAt
            permalink
          }
          error {
            title
            body
            errorType
          }
        }
      }
    }
  }
}  
  `,
		owner,
		repo,
	)
}
