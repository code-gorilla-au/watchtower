package insights

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/odize"
)

func TestService_Insights(t *testing.T) {
	group := odize.NewGroup(t, nil)

	ctx := context.Background()
	svc := New(_testDB)

	group.BeforeEach(func() {
		cleanup(t)
	})

	err := group.
		Test("should return empty insights when no data exists", func(t *testing.T) {
			orgID := seedOrg(t, "test-org")
			filterDate := "-30 days"

			prInsights, err := svc.GetPullRequestInsightsByOrg(ctx, orgID, filterDate)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, prInsights.Merged == 0)

			secInsights, err := svc.GetSecurityInsightsByOrg(ctx, orgID, filterDate)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, secInsights.Fixed == 0)
		}).
		Test("should return pull request insights for an organisation", func(t *testing.T) {
			orgID := seedOrg(t, "test-org")
			productID := seedProduct(t, "test-product", "[\"topic-a\"]")
			seedProductOrg(t, productID, orgID)
			seedRepo(t, "repo-a", "topic-a")

			now := time.Now().Unix()
			twoDaysAgo := time.Now().AddDate(0, 0, -2).Unix()

			// Merged PR: 2 days to merge
			seedPullRequest(t, "pr-1", "repo-a", "MERGED", twoDaysAgo, now)
			// Open PR
			seedPullRequest(t, "pr-2", "repo-a", "OPEN", twoDaysAgo, 0)

			filterDate := "-30 days"
			insights, err := svc.GetPullRequestInsightsByOrg(ctx, orgID, filterDate)
			odize.AssertNoError(t, err)

			odize.AssertTrue(t, insights.Merged == 1)
			odize.AssertTrue(t, insights.Open == 1)
			odize.AssertTrue(t, insights.AvgDaysToMerge == 2.0)
		}).
		Test("should return security insights for an organisation", func(t *testing.T) {
			orgID := seedOrg(t, "test-org")
			productID := seedProduct(t, "test-product", "[\"topic-s\"]")
			seedProductOrg(t, productID, orgID)
			seedRepo(t, "repo-s", "topic-s")

			now := time.Now().Unix()
			threeDaysAgo := time.Now().AddDate(0, 0, -3).Unix()

			// Fixed security issue: 3 days to fix
			seedSecurity(t, "sec-1", "repo-s", "FIXED", threeDaysAgo, now)
			// Open security issue
			seedSecurity(t, "sec-2", "repo-s", "OPEN", threeDaysAgo, 0)

			filterDate := "-30 days"
			insights, err := svc.GetSecurityInsightsByOrg(ctx, orgID, filterDate)
			odize.AssertNoError(t, err)

			odize.AssertTrue(t, insights.Fixed == 1)
			odize.AssertTrue(t, insights.Open == 1)
			odize.AssertTrue(t, insights.AvgDaysToFix == 3.0)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func cleanup(t *testing.T) {
	_, err := _testRawDB.Exec("DELETE FROM pull_requests")
	odize.AssertNoError(t, err)
	_, err = _testRawDB.Exec("DELETE FROM securities")
	odize.AssertNoError(t, err)
	_, err = _testRawDB.Exec("DELETE FROM repositories")
	odize.AssertNoError(t, err)
	_, err = _testRawDB.Exec("DELETE FROM product_organisations")
	odize.AssertNoError(t, err)
	_, err = _testRawDB.Exec("DELETE FROM products")
	odize.AssertNoError(t, err)
	_, err = _testRawDB.Exec("DELETE FROM organisations")
	odize.AssertNoError(t, err)
}

func seedOrg(t *testing.T, name string) int64 {
	org, err := _testDB.CreateOrganisation(context.Background(), database.CreateOrganisationParams{
		FriendlyName: name,
		Namespace:    name,
		Token:        "token",
		Description:  "desc",
	})
	odize.AssertNoError(t, err)
	return org.ID
}

func seedProduct(t *testing.T, name string, tags string) int64 {
	product, err := _testDB.CreateProduct(context.Background(), database.CreateProductParams{
		Name:        name,
		Description: "desc",
		Tags:        sql.NullString{String: tags, Valid: true},
	})
	odize.AssertNoError(t, err)
	return product.ID
}

func seedProductOrg(t *testing.T, productID, orgID int64) {
	err := _testDB.AddProductToOrganisation(context.Background(), database.AddProductToOrganisationParams{
		ProductID:      sql.NullInt64{Int64: productID, Valid: true},
		OrganisationID: sql.NullInt64{Int64: orgID, Valid: true},
	})
	odize.AssertNoError(t, err)
}

func seedRepo(t *testing.T, name, topic string) {
	_, err := _testDB.CreateRepo(context.Background(), database.CreateRepoParams{
		Name:  name,
		Url:   "url",
		Topic: topic,
		Owner: "owner",
	})
	odize.AssertNoError(t, err)
}

func seedPullRequest(t *testing.T, externalID, repo, state string, createdAt, mergedAt int64) {
	var mAt sql.NullInt64
	if mergedAt > 0 {
		mAt = sql.NullInt64{Int64: mergedAt, Valid: true}
	}
	_, err := _testDB.CreatePullRequest(context.Background(), database.CreatePullRequestParams{
		ExternalID:     externalID,
		Title:          "title",
		RepositoryName: repo,
		Url:            "url",
		State:          state,
		Author:         "author",
		CreatedAt:      createdAt,
		MergedAt:       mAt,
	})
	odize.AssertNoError(t, err)
}

func seedSecurity(t *testing.T, externalID, repo, state string, createdAt, fixedAt int64) {
	var fAt sql.NullInt64
	if fixedAt > 0 {
		fAt = sql.NullInt64{Int64: fixedAt, Valid: true}
	}
	_, err := _testDB.CreateSecurity(context.Background(), database.CreateSecurityParams{
		ExternalID:     externalID,
		RepositoryName: repo,
		PackageName:    "pkg",
		State:          state,
		Severity:       "high",
		PatchedVersion: "1.0.0",
		CreatedAt:      createdAt,
		FixedAt:        fAt,
	})
	odize.AssertNoError(t, err)
}
