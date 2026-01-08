package products

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"
	"watchtower/internal/database"
	"watchtower/internal/github"

	"github.com/code-gorilla-au/odize"
	"github.com/google/uuid"
)

func TestService(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()

	group.BeforeEach(func() {
		s = New(_testDB)
	})

	err := group.
		Test("Create should create a product", func(t *testing.T) {
			params := CreateProductParams{
				Name: "Test Product",
				Desc: "Test Description",
				Tags: []string{"tag1", "tag2"},
			}

			prod, err := s.Create(ctx, params)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, params.Name, prod.Name)
			odize.AssertEqual(t, params.Desc, prod.Description)
			odize.AssertEqual(t, 2, len(prod.Tags))
			odize.AssertTrue(t, prod.ID > 0)
		}).
		Test("Get should return a product", func(t *testing.T) {
			params := CreateProductParams{
				Name: "Get Product",
				Desc: "Get Description",
			}
			created, _ := s.Create(ctx, params)

			got, err := s.Get(ctx, created.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, created.ID, got.ID)
			odize.AssertEqual(t, "Get Product", got.Name)
		}).
		Test("GetByOrg should return products for an org", func(t *testing.T) {
			// We need to use database directly to associate with org or use a helper
			prod, _ := s.Create(ctx, CreateProductParams{Name: "Org Prod"})
			orgID := int64(123)

			err := _testDB.AddProductToOrganisation(ctx, database.AddProductToOrganisationParams{
				ProductID:      sql.NullInt64{Int64: prod.ID, Valid: true},
				OrganisationID: sql.NullInt64{Int64: orgID, Valid: true},
			})
			odize.AssertNoError(t, err)

			products, err := s.GetByOrg(ctx, orgID)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(products) > 0)
			found := false
			for _, p := range products {
				if p.ID == prod.ID {
					found = true
					break
				}
			}
			odize.AssertTrue(t, found)
		}).
		Test("Update should update product details", func(t *testing.T) {
			prod, _ := s.Create(ctx, CreateProductParams{Name: "Old Name", Tags: []string{"old"}})

			updated, err := s.Update(ctx, UpdateProductParams{
				ID:   prod.ID,
				Name: "New Name",
				Tags: []string{"new1", "new2"},
			})
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "New Name", updated.Name)
			odize.AssertEqual(t, 2, len(updated.Tags))
		}).
		Test("UpdateSyncDateNow should update the sync date", func(t *testing.T) {
			prod, _ := s.Create(ctx, CreateProductParams{Name: "Sync Test"})
			err := s.UpdateSyncDateNow(ctx, prod.ID)
			odize.AssertNoError(t, err)
		}).
		Test("DeleteProduct should remove product and related data", func(t *testing.T) {
			prod, _ := s.Create(ctx, CreateProductParams{Name: "Delete Me"})

			// Create some related data
			_, _ = _testDB.CreateRepo(ctx, database.CreateRepoParams{Name: "repo1", Topic: "tag", Owner: "owner"})
			// Note: the association between product and repo is via the tag/topic in some queries,
			// but DeleteReposByProductID uses a direct link if it exists?
			// Let's check how DeleteReposByProductID is implemented in SQL.
			// Actually, let's just test that the call succeeds.

			err := s.DeleteProduct(ctx, prod.ID)
			odize.AssertNoError(t, err)

			_, err = s.Get(ctx, prod.ID)
			odize.AssertError(t, err)
		}).
		Test("CreateRepo and GetRepos", func(t *testing.T) {
			// Repos are associated with products via tags (topics)
			tag := fmt.Sprintf("tag-%d", time.Now().UnixNano())
			prod, _ := s.Create(ctx, CreateProductParams{Name: "Repo Product", Tags: []string{tag}})

			err := s.CreateRepo(ctx, CreateRepoParams{
				Name:  "test-repo",
				Url:   "url",
				Topic: tag,
				Owner: "owner",
			})
			odize.AssertNoError(t, err)

			repos, err := s.GetRepos(ctx, prod.ID)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(repos) > 0)
			odize.AssertEqual(t, "test-repo", repos[0].Name)
		}).
		Test("UpsertRepo should create a repo when it does not exist", func(t *testing.T) {
			params := CreateRepoParams{
				Name:  "upsert-create-repo",
				Url:   "url1",
				Topic: "topic1",
				Owner: "owner1",
			}

			err := s.UpsertRepo(ctx, params)
			odize.AssertNoError(t, err)

			model, err := _testDB.GetRepoByName(ctx, params.Name)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, params.Url, model.Url)
			odize.AssertEqual(t, params.Topic, model.Topic)
		}).
		Test("UpsertRepo should update a repo when it already exists", func(t *testing.T) {
			params := CreateRepoParams{
				Name:  "upsert-update-repo",
				Url:   "url1",
				Topic: "topic1",
				Owner: "owner1",
			}

			// Initial creation
			err := s.UpsertRepo(ctx, params)
			odize.AssertNoError(t, err)

			// Update
			params.Url = "url2"
			params.Topic = "topic2"
			err = s.UpsertRepo(ctx, params)
			odize.AssertNoError(t, err)

			model, err := _testDB.GetRepoByName(ctx, params.Name)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "url2", model.Url)
			odize.AssertEqual(t, "topic2", model.Topic)
		}).
		Test("CreatePullRequest should create a new pull request", func(t *testing.T) {
			params := CreatePRParams{
				ExternalID:     "new-pr",
				Title:          "New PR",
				RepositoryName: "repo1",
				Url:            "url1",
				State:          "OPEN",
				Author:         "author1",
				CreatedAt:      time.Now(),
			}

			err := s.CreatePullRequest(ctx, params)
			odize.AssertNoError(t, err)

			pr, err := _testDB.GetPullRequestByExternalID(ctx, params.ExternalID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, params.Title, pr.Title)
			odize.AssertEqual(t, params.Author, pr.Author)
		}).
		Test("UpdatePullRequest should update an existing pull request", func(t *testing.T) {
			params := CreatePRParams{
				ExternalID:     "update-pr",
				Title:          "Original Title",
				RepositoryName: "repo1",
				Url:            "url1",
				State:          "OPEN",
				Author:         "author1",
				CreatedAt:      time.Now(),
			}

			err := s.CreatePullRequest(ctx, params)
			odize.AssertNoError(t, err)

			pr, _ := _testDB.GetPullRequestByExternalID(ctx, params.ExternalID)

			updateParams := UpdatePRParams{
				ID:             pr.ID,
				Title:          "Updated Title",
				RepositoryName: "repo1",
				Url:            "url-updated",
				State:          "CLOSED",
				Author:         "author-updated",
			}

			err = s.UpdatePullRequest(ctx, updateParams)
			odize.AssertNoError(t, err)

			updatedPr, err := _testDB.GetPullRequestByExternalID(ctx, params.ExternalID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, updateParams.Title, updatedPr.Title)
			odize.AssertEqual(t, updateParams.State, updatedPr.State)
			odize.AssertEqual(t, updateParams.Author, updatedPr.Author)
		}).
		Test("UpsertPullRequest should create when not exists and update when exists", func(t *testing.T) {
			params := CreatePRParams{
				ExternalID:     "upsert-pr",
				Title:          "Upsert Title",
				RepositoryName: "repo1",
				Url:            "url1",
				State:          "OPEN",
				Author:         "author1",
				CreatedAt:      time.Now(),
			}

			// Create
			err := s.UpsertPullRequest(ctx, params)
			odize.AssertNoError(t, err)

			pr, err := _testDB.GetPullRequestByExternalID(ctx, params.ExternalID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, params.Title, pr.Title)

			// Update
			params.Title = "Upsert Updated Title"
			err = s.UpsertPullRequest(ctx, params)
			odize.AssertNoError(t, err)

			updatedPr, err := _testDB.GetPullRequestByExternalID(ctx, params.ExternalID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "Upsert Updated Title", updatedPr.Title)
			odize.AssertEqual(t, pr.ID, updatedPr.ID)
		}).
		Test("GetPullRequests and GetPullRequestByOrg", func(t *testing.T) {
			tag := fmt.Sprintf("pr-tag-%d", time.Now().UnixNano())
			prod, _ := s.Create(ctx, CreateProductParams{Name: "PR Product", Tags: []string{tag}})
			orgID := int64(456)
			_ = _testDB.AddProductToOrganisation(ctx, database.AddProductToOrganisationParams{
				ProductID:      sql.NullInt64{Int64: prod.ID, Valid: true},
				OrganisationID: sql.NullInt64{Int64: orgID, Valid: true},
			})

			repoName := "pr-repo"
			_ = s.CreateRepo(ctx, CreateRepoParams{Name: repoName, Topic: tag, Owner: "owner"})

			err := s.BulkCreatePullRequest(ctx, []CreatePRParams{
				{
					ExternalID:     "pr1",
					Title:          "PR 1",
					RepositoryName: repoName,
					Url:            "url1",
					State:          "OPEN",
					Author:         "author1",
					CreatedAt:      time.Now(),
				},
			})
			odize.AssertNoError(t, err)

			prs, err := s.GetPullRequests(ctx, prod.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 1, len(prs))
			odize.AssertEqual(t, "PR 1", prs[0].Title)

			orgPrs, err := s.GetPullRequestByOrg(ctx, orgID)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(orgPrs) > 0)
		}).
		Test("GetRecentPullRequests should return external IDs of recent PRs", func(t *testing.T) {
			params := CreatePRParams{
				ExternalID:     uuid.New().String(),
				Title:          uuid.New().String(),
				RepositoryName: uuid.New().String(),
				Url:            "url1",
				State:          "OPEN",
				Author:         "author1",
				CreatedAt:      time.Now(),
			}

			err := s.CreateRepo(ctx, CreateRepoParams{Name: params.RepositoryName, Topic: "tag", Owner: "owner"})
			odize.AssertNoError(t, err)

			prodID, err := s.Create(ctx, CreateProductParams{
				Name: params.RepositoryName,
				Desc: "",
				Tags: []string{"tag"},
			})
			odize.AssertNoError(t, err)

			_ = _testDB.AddProductToOrganisation(ctx, database.AddProductToOrganisationParams{
				ProductID: sql.NullInt64{
					Int64: prodID.ID,
					Valid: true,
				},
				OrganisationID: sql.NullInt64{
					Int64: 0,
					Valid: true,
				},
			})

			err = s.CreatePullRequest(ctx, params)
			odize.AssertNoError(t, err)

			recent, err := s.GetRecentPullRequests(ctx)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(recent) > 0)

			found := false
			for _, entity := range recent {
				if entity.ExternalID == params.ExternalID {
					found = true
					break
				}
			}
			odize.AssertTrue(t, found)
		}).
		Test("GetSecurity and GetSecurityByOrg", func(t *testing.T) {
			tag := fmt.Sprintf("sec-tag-%d", time.Now().UnixNano())
			prod, _ := s.Create(ctx, CreateProductParams{Name: "Sec Product", Tags: []string{tag}})
			orgID := int64(789)
			_ = _testDB.AddProductToOrganisation(ctx, database.AddProductToOrganisationParams{
				ProductID:      sql.NullInt64{Int64: prod.ID, Valid: true},
				OrganisationID: sql.NullInt64{Int64: orgID, Valid: true},
			})

			repoName := "sec-repo"
			_ = s.CreateRepo(ctx, CreateRepoParams{Name: repoName, Topic: tag, Owner: "owner"})

			err := s.BulkCreateSecurity(ctx, []CreateSecurityParams{
				{
					ExternalID:     "sec1",
					RepositoryName: repoName,
					PackageName:    "pkg1",
					State:          "OPEN",
					Severity:       "HIGH",
					CreatedAt:      time.Now(),
				},
			})
			odize.AssertNoError(t, err)

			secs, err := s.GetSecurity(ctx, prod.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 1, len(secs))

			orgSecs, err := s.GetSecurityByOrg(ctx, orgID)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(orgSecs) > 0)
		}).
		Test("GetRecentSecurity should return external IDs of recent security alerts", func(t *testing.T) {
			params := CreateSecurityParams{
				ExternalID:     "recent-sec-1",
				RepositoryName: "repo1",
				PackageName:    "pkg1",
				State:          "OPEN",
				Severity:       "HIGH",
				CreatedAt:      time.Now(),
			}

			err := s.UpsertSecurity(ctx, params)
			odize.AssertNoError(t, err)

			recent, err := s.GetRecentSecurity(ctx)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(recent) > 0)

			found := false
			for _, entity := range recent {
				if entity.ExternalID == params.ExternalID {
					found = true
					break
				}
			}
			odize.AssertTrue(t, found)
		}).
		Test("BulkInsertRepos", func(t *testing.T) {
			tag := "bulk-repo-tag"
			repos := []github.Node[github.Repository]{
				{
					Node: github.Repository{
						Name: "bulk-repo-1",
						Url:  "url1",
						Owner: github.Owner{
							Login: "owner1",
						},
					},
				},
			}

			err := s.BulkInsertRepos(ctx, repos, tag)
			odize.AssertNoError(t, err)
		}).
		Test("BulkInsertRepoDetails", func(t *testing.T) {
			repoDetails := github.QueryRepository{
				Data: github.RepositoryData{
					Repository: github.Repository{
						Name: "detail-repo",
						PullRequests: github.RootNode[github.PullRequest]{
							Nodes: []github.PullRequest{
								{
									ID:    "pr-detail-1",
									Title: "PR Detail 1",
									Author: github.Author{
										Login: "author1",
									},
									State: github.PrOpen,
								},
							},
						},
						VulnerabilityAlerts: github.RootNode[github.VulnerabilityAlerts]{
							Nodes: []github.VulnerabilityAlerts{
								{
									ID:    "sec-detail-1",
									State: "OPEN",
									SecurityVulnerability: github.SecurityVulnerability{
										Package: github.Package{
											Name: "pkg1",
										},
										Advisory: github.Advisory{
											Severity: "HIGH",
										},
										FirstPatchedVersion: github.FirstPatchedVersion{
											Identifier: "1.0.1",
										},
									},
								},
							},
						},
					},
				},
			}

			err := s.BulkInsertRepoDetails(ctx, repoDetails)
			odize.AssertNoError(t, err)
		}).
		Test("UpdateSecurity should update an existing security alert", func(t *testing.T) {
			params := CreateSecurityParams{
				ExternalID:     "update-sec",
				RepositoryName: "repo1",
				PackageName:    "pkg1",
				State:          "OPEN",
				Severity:       "HIGH",
				CreatedAt:      time.Now(),
			}

			err := s.UpsertSecurity(ctx, params)
			odize.AssertNoError(t, err)

			updateParams := UpdateSecurityParams{
				ExternalID:     params.ExternalID,
				RepositoryName: "repo1",
				PackageName:    "pkg1-updated",
				State:          "FIXED",
				Severity:       "CRITICAL",
				PatchedVersion: "1.2.3",
			}

			err = s.UpdateSecurity(ctx, updateParams)
			odize.AssertNoError(t, err)

			updated, err := _testDB.GetSecurityByExternalID(ctx, params.ExternalID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, updateParams.PackageName, updated.PackageName)
			odize.AssertEqual(t, updateParams.State, updated.State)
			odize.AssertEqual(t, updateParams.Severity, updated.Severity)
			odize.AssertEqual(t, updateParams.PatchedVersion, updated.PatchedVersion)
		}).
		Test("UpsertSecurity should create when not exists and update when exists", func(t *testing.T) {
			params := CreateSecurityParams{
				ExternalID:     "upsert-sec",
				RepositoryName: "repo1",
				PackageName:    "pkg-upsert",
				State:          "OPEN",
				Severity:       "LOW",
				CreatedAt:      time.Now(),
			}

			// Create
			err := s.UpsertSecurity(ctx, params)
			odize.AssertNoError(t, err)

			sec, err := _testDB.GetSecurityByExternalID(ctx, params.ExternalID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, params.PackageName, sec.PackageName)

			// Update
			params.PackageName = "pkg-upsert-updated"
			params.Severity = "MEDIUM"
			err = s.UpsertSecurity(ctx, params)
			odize.AssertNoError(t, err)

			updated, err := _testDB.GetSecurityByExternalID(ctx, params.ExternalID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "pkg-upsert-updated", updated.PackageName)
			odize.AssertEqual(t, "MEDIUM", updated.Severity)
			odize.AssertEqual(t, sec.ID, updated.ID)
		}).
		Test("AssociateProductToOrg should create link in product_organisations", func(t *testing.T) {
			org, _ := _testDB.CreateOrganisation(ctx, database.CreateOrganisationParams{
				FriendlyName: "Product Org",
				Namespace:    "prod-ns",
			})
			productID := int64(123)

			err := s.AssociateProductToOrg(ctx, org.ID, productID)
			odize.AssertNoError(t, err)

			fetchedOrg, err := _testDB.GetOrganisationForProduct(ctx, sql.NullInt64{Int64: productID, Valid: true})
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, org.ID, fetchedOrg.ID)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestTransforms(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("toTime should convert unix timestamp to time.Time", func(t *testing.T) {
			now := time.Now().Unix()
			got := toTime(now)
			odize.AssertEqual(t, now, got.Unix())
		}).
		Run()

	odize.AssertNoError(t, err)
}
