package organisations

import (
	"context"
	"database/sql"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestService(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()

	group.BeforeEach(func() {
		s = New(_testDB, _testTxnDB, func(tx *sql.Tx) OrgStore {
			return _testDB.WithTx(tx)
		})
	})

	err := group.
		Test("Create should create an organisation", func(t *testing.T) {
			params := CreateOrgParams{
				FriendlyName: "Test Org",
				Namespace:    "test-ns",
				Token:        "test-token",
				Description:  "test-desc",
			}

			org, err := s.Create(ctx, params)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, params.FriendlyName, org.FriendlyName)
			odize.AssertEqual(t, params.Namespace, org.Namespace)
			odize.AssertEqual(t, params.Description, org.Description)
			odize.AssertTrue(t, org.ID > 0)
			odize.AssertTrue(t, org.DefaultOrg) // First org should be default
		}).
		Test("Get should return organisation by ID", func(t *testing.T) {
			org, err := s.Create(ctx, CreateOrgParams{
				FriendlyName: "Get Org",
				Namespace:    "get-ns",
			})
			odize.AssertNoError(t, err)

			fetched, err := s.Get(ctx, org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, org.ID, fetched.ID)
			odize.AssertEqual(t, "Get Org", fetched.FriendlyName)
		}).
		Test("GetAll should return all organisations", func(t *testing.T) {
			org1, _ := s.Create(ctx, CreateOrgParams{FriendlyName: "Org 1", Namespace: "ns-getall-1"})
			org2, _ := s.Create(ctx, CreateOrgParams{FriendlyName: "Org 2", Namespace: "ns-getall-2"})

			all, err := s.GetAll(ctx)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(all) >= 2)

			found1 := false
			found2 := false
			for _, o := range all {
				if o.ID == org1.ID {
					found1 = true
				}
				if o.ID == org2.ID {
					found2 = true
				}
			}
			odize.AssertTrue(t, found1)
			odize.AssertTrue(t, found2)
		}).
		Test("Update should update organisation details", func(t *testing.T) {
			org, _ := s.Create(ctx, CreateOrgParams{FriendlyName: "Old Name", Namespace: "old-ns-update"})

			updated, err := s.Update(ctx, UpdateOrgParams{
				ID:           org.ID,
				FriendlyName: "New Name",
				Namespace:    "new-ns-update",
				Description:  "new-desc",
				DefaultOrg:   true,
			})
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "New Name", updated.FriendlyName)
			odize.AssertEqual(t, "new-ns-update", updated.Namespace)
			odize.AssertEqual(t, "new-desc", updated.Description)
		}).
		Test("SetDefault should change the default organisation", func(t *testing.T) {
			org1, err := s.Create(ctx, CreateOrgParams{FriendlyName: "Org 1", Namespace: "ns-set-default-1"})
			odize.AssertNoError(t, err)
			org2, err := s.Create(ctx, CreateOrgParams{FriendlyName: "Org 2", Namespace: "ns-set-default-2"})
			odize.AssertNoError(t, err)

			_, err = s.SetDefault(ctx, org2.ID)
			odize.AssertNoError(t, err)

			def, err := s.GetDefault(ctx)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, org2.ID, def.ID)

			o1, _ := s.Get(ctx, org1.ID)
			odize.AssertFalse(t, o1.DefaultOrg)
		}).
		Test("Delete should remove organisation", func(t *testing.T) {
			org, _ := s.Create(ctx, CreateOrgParams{FriendlyName: "To Delete", Namespace: "delete-ns-unique"})

			err := s.Delete(ctx, org.ID)
			odize.AssertNoError(t, err)

			_, err = s.Get(ctx, org.ID)
			odize.AssertTrue(t, err != nil)
		}).
		Test("GetStaleOrgs should return orgs not synced recently", func(t *testing.T) {
			_, err := s.GetStaleOrgs(ctx)
			odize.AssertNoError(t, err)

		}).
		Run()

	odize.AssertNoError(t, err)
}
