package watchtower

import (
	"context"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestService_Organisations(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
	})

	err := group.
		Test("DeleteAllOrgs should delete all organisations and their products", func(t *testing.T) {
			// 1. Create orgs
			org1, err := s.CreateOrganisation("Org 1", "ns-1", "token-1", "desc-1")
			odize.AssertNoError(t, err)

			org2, err := s.CreateOrganisation("Org 2", "ns-2", "token-2", "desc-2")
			odize.AssertNoError(t, err)

			// 2. Create products for orgs
			_, err = s.CreateProduct("Prod 1", "desc-p1", []string{"tag1"}, org1.ID)
			odize.AssertNoError(t, err)

			_, err = s.CreateProduct("Prod 2", "desc-p2", []string{"tag2"}, org2.ID)
			odize.AssertNoError(t, err)

			// 3. Verify they exist
			allOrgs, err := s.GetAllOrganisations()
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(allOrgs) >= 2)

			prods1, err := s.GetAllProductsForOrganisation(org1.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 1, len(prods1))

			prods2, err := s.GetAllProductsForOrganisation(org2.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 1, len(prods2))

			// 4. Delete all
			err = s.DeleteAllOrgs()
			odize.AssertNoError(t, err)

			// 5. Verify they are gone
			remainingOrgs, err := s.GetAllOrganisations()
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 0, len(remainingOrgs))

			// Check if products are also gone
			remainingProds1, _ := s.productSvc.GetByOrg(ctx, org1.ID)
			odize.AssertEqual(t, 0, len(remainingProds1))

			remainingProds2, _ := s.productSvc.GetByOrg(ctx, org2.ID)
			odize.AssertEqual(t, 0, len(remainingProds2))
		}).
		Run()

	odize.AssertNoError(t, err)
}
