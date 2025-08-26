package watchtower

import (
	"context"
	"testing"
	"time"

	"github.com/code-gorilla-au/odize"
)

func TestService_CreateOrganisation(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should create an organisation", func(t *testing.T) {
			org, err := s.CreateOrganisation("friendly_name", "namespace", "token", "description")
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, org.FriendlyName, "friendly_name")
			odize.AssertEqual(t, org.Namespace, "namespace")
			odize.AssertEqual(t, org.Description, "description")
			odize.AssertEqual(t, org.DefaultOrg, true)
			odize.AssertFalse(t, org.CreatedAt.Equal(time.Time{}))
		}).
		Test("should return error if trying to create same org", func(t *testing.T) {
			_, err := s.CreateOrganisation("friendly_name", "namespace", "token", "description")
			odize.AssertError(t, err)

		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_GetDefaultOrganisation(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should return error when no default organisation exists", func(t *testing.T) {
			_, err := s.GetDefaultOrganisation()
			odize.AssertError(t, err)
		}).
		Test("should return default organisation when it exists", func(t *testing.T) {

			createdOrg, err := s.CreateOrganisation("default_org", "default_namespace", "token", "default description")
			odize.AssertNoError(t, err)

			defaultOrg, err := s.GetDefaultOrganisation()
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, defaultOrg.ID, createdOrg.ID)
			odize.AssertEqual(t, defaultOrg.FriendlyName, "default_org")
			odize.AssertEqual(t, defaultOrg.Namespace, "default_namespace")
			odize.AssertEqual(t, defaultOrg.Description, "default description")
			odize.AssertEqual(t, defaultOrg.DefaultOrg, true)
			odize.AssertFalse(t, defaultOrg.CreatedAt.Equal(time.Time{}))
		}).
		Test("should return correct default organisation when multiple orgs exist", func(t *testing.T) {
			firstOrg, err := s.CreateOrganisation("first_org", "first_namespace", "token1", "first description")
			odize.AssertNoError(t, err)

			secondOrg, err := s.CreateOrganisation("second_org", "second_namespace", "token2", "second description")
			odize.AssertNoError(t, err)

			defaultOrg, err := s.GetDefaultOrganisation()
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, defaultOrg.ID, secondOrg.ID)
			odize.AssertEqual(t, defaultOrg.FriendlyName, "second_org")
			odize.AssertEqual(t, defaultOrg.Namespace, "second_namespace")
			odize.AssertEqual(t, defaultOrg.Description, "second description")
			odize.AssertEqual(t, defaultOrg.DefaultOrg, true)

			odize.AssertFalse(t, defaultOrg.ID == firstOrg.ID)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_SetDefaultOrg(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should return error when trying to set non-existent org as default", func(t *testing.T) {
			_, err := s.SetDefaultOrg(999)
			odize.AssertError(t, err)
		}).
		Test("should successfully set existing org as default", func(t *testing.T) {
			createdOrg, err := s.CreateOrganisation("test_org", "test_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			setOrg, err := s.SetDefaultOrg(createdOrg.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, setOrg.ID, createdOrg.ID)
			odize.AssertEqual(t, setOrg.FriendlyName, "test_org")
			odize.AssertEqual(t, setOrg.Namespace, "test_namespace")
			odize.AssertEqual(t, setOrg.Description, "test description")
			odize.AssertEqual(t, setOrg.DefaultOrg, true)
		}).
		Test("should unset previous default when setting new default", func(t *testing.T) {
			firstOrg, err := s.CreateOrganisation("first_org", "first_namespace_unique", "token1", "first description")
			odize.AssertNoError(t, err)

			_, err = s.CreateOrganisation("second_org", "second_namespace_unique", "token2", "second description")
			odize.AssertNoError(t, err)

			setOrg, err := s.SetDefaultOrg(firstOrg.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, setOrg.ID, firstOrg.ID)
			odize.AssertEqual(t, setOrg.DefaultOrg, true)

			defaultOrg, err := s.GetDefaultOrganisation()
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, defaultOrg.ID, firstOrg.ID)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_GetOrganisationByID(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should return error when organisation does not exist", func(t *testing.T) {
			_, err := s.GetOrganisationByID(999)
			odize.AssertError(t, err)
		}).
		Test("should return organisation when it exists", func(t *testing.T) {
			createdOrg, err := s.CreateOrganisation("test_org_byid", "test_namespace_byid", "token", "test description")
			odize.AssertNoError(t, err)

			fetchedOrg, err := s.GetOrganisationByID(createdOrg.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, fetchedOrg.ID, createdOrg.ID)
			odize.AssertEqual(t, fetchedOrg.FriendlyName, "test_org_byid")
			odize.AssertEqual(t, fetchedOrg.Namespace, "test_namespace_byid")
			odize.AssertEqual(t, fetchedOrg.Description, "test description")
			odize.AssertEqual(t, fetchedOrg.DefaultOrg, true)
			odize.AssertFalse(t, fetchedOrg.CreatedAt.Equal(time.Time{}))
		}).
		Test("should return correct organisation when multiple orgs exist", func(t *testing.T) {
			firstOrg, err := s.CreateOrganisation("first_org", "first_namespace_get", "token1", "first description")
			odize.AssertNoError(t, err)

			secondOrg, err := s.CreateOrganisation("second_org", "second_namespace_get", "token2", "second description")
			odize.AssertNoError(t, err)

			fetchedFirstOrg, err := s.GetOrganisationByID(firstOrg.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, fetchedFirstOrg.ID, firstOrg.ID)
			odize.AssertEqual(t, fetchedFirstOrg.FriendlyName, "first_org")
			odize.AssertEqual(t, fetchedFirstOrg.Namespace, "first_namespace_get")
			odize.AssertEqual(t, fetchedFirstOrg.Description, "first description")
			odize.AssertEqual(t, fetchedFirstOrg.DefaultOrg, false)

			fetchedSecondOrg, err := s.GetOrganisationByID(secondOrg.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, fetchedSecondOrg.ID, secondOrg.ID)
			odize.AssertEqual(t, fetchedSecondOrg.FriendlyName, "second_org")
			odize.AssertEqual(t, fetchedSecondOrg.Namespace, "second_namespace_get")
			odize.AssertEqual(t, fetchedSecondOrg.Description, "second description")
			odize.AssertEqual(t, fetchedSecondOrg.DefaultOrg, true)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_GetAllOrganisations(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
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

func TestService_DeleteOrganisation(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should not return error when trying to delete non-existent organisation", func(t *testing.T) {
			err := s.DeleteOrganisation(999)
			odize.AssertNoError(t, err)
		}).
		Test("should successfully delete existing organisation", func(t *testing.T) {
			createdOrg, err := s.CreateOrganisation("delete_test_org", "delete_test_namespace", "token", "delete test description")
			odize.AssertNoError(t, err)

			err = s.DeleteOrganisation(createdOrg.ID)
			odize.AssertNoError(t, err)
		}).
		Test("should not be able to retrieve deleted organisation", func(t *testing.T) {
			createdOrg, err := s.CreateOrganisation("verify_delete_org", "verify_delete_namespace", "token", "verify delete description")
			odize.AssertNoError(t, err)

			err = s.DeleteOrganisation(createdOrg.ID)
			odize.AssertNoError(t, err)

			_, err = s.GetOrganisationByID(createdOrg.ID)
			odize.AssertError(t, err)
		}).
		Test("should successfully delete default organisation", func(t *testing.T) {
			createdOrg, err := s.CreateOrganisation("default_delete_org", "default_delete_namespace", "token", "default delete description")
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, createdOrg.DefaultOrg, true)

			err = s.DeleteOrganisation(createdOrg.ID)
			odize.AssertNoError(t, err)

			_, err = s.GetOrganisationByID(createdOrg.ID)
			odize.AssertError(t, err)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_UpdateOrganisation(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should return error when trying to update non-existent organisation", func(t *testing.T) {
			params := UpdateOrgParams{
				ID:           999,
				FriendlyName: "updated_name",
				Namespace:    "updated_namespace",
				Description:  "updated_description",
				DefaultOrg:   false,
			}
			_, err := s.UpdateOrganisation(params)
			odize.AssertError(t, err)
		}).
		Test("should successfully update organisation fields", func(t *testing.T) {
			createdOrg, err := s.CreateOrganisation("original_org", "original_namespace", "token", "original description")
			odize.AssertNoError(t, err)

			params := UpdateOrgParams{
				ID:           createdOrg.ID,
				FriendlyName: "updated_friendly_name",
				Namespace:    "updated_namespace",
				Description:  "updated_description",
				DefaultOrg:   false,
			}
			updatedOrg, err := s.UpdateOrganisation(params)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, updatedOrg.ID, createdOrg.ID)
			odize.AssertEqual(t, updatedOrg.FriendlyName, "updated_friendly_name")
			odize.AssertEqual(t, updatedOrg.Namespace, "updated_namespace")
			odize.AssertEqual(t, updatedOrg.DefaultOrg, false)
			odize.AssertEqual(t, updatedOrg.Description, "updated_description")
		}).
		Test("should successfully set organisation as default", func(t *testing.T) {
			firstOrg, err := s.CreateOrganisation("first_update_org", "first_update_namespace", "token1", "first description")
			odize.AssertNoError(t, err)

			secondOrg, err := s.CreateOrganisation("second_update_org", "second_update_namespace", "token2", "second description")
			odize.AssertNoError(t, err)

			params := UpdateOrgParams{
				ID:           firstOrg.ID,
				FriendlyName: "updated_first_org",
				Namespace:    "updated_first_namespace",
				Description:  "updated_first_description",
				DefaultOrg:   true,
			}
			updatedOrg, err := s.UpdateOrganisation(params)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, updatedOrg.ID, firstOrg.ID)
			odize.AssertEqual(t, updatedOrg.FriendlyName, "updated_first_org")
			odize.AssertEqual(t, updatedOrg.Namespace, "updated_first_namespace")
			odize.AssertEqual(t, updatedOrg.DefaultOrg, true)

			fetchedSecondOrg, err := s.GetOrganisationByID(secondOrg.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, fetchedSecondOrg.DefaultOrg, false)
		}).
		Test("should unset previous defaults when setting new default", func(t *testing.T) {
			firstOrg, err := s.CreateOrganisation("prev_default_org", "prev_default_namespace", "token1", "prev description")
			odize.AssertNoError(t, err)

			secondOrg, err := s.CreateOrganisation("new_default_org", "new_default_namespace", "token2", "new description")
			odize.AssertNoError(t, err)

			params := UpdateOrgParams{
				ID:           firstOrg.ID,
				FriendlyName: firstOrg.FriendlyName,
				Namespace:    firstOrg.Namespace,
				Description:  firstOrg.Description,
				DefaultOrg:   true,
			}
			_, err = s.UpdateOrganisation(params)
			odize.AssertNoError(t, err)

			defaultOrg, err := s.GetDefaultOrganisation()
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, defaultOrg.ID, firstOrg.ID)

			fetchedSecondOrg, err := s.GetOrganisationByID(secondOrg.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, fetchedSecondOrg.DefaultOrg, false)
		}).
		Test("should update without changing default status when DefaultOrg is false", func(t *testing.T) {
			createdOrg, err := s.CreateOrganisation("maintain_default_org", "maintain_default_namespace", "token", "maintain description")
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, createdOrg.DefaultOrg, true)

			params := UpdateOrgParams{
				ID:           createdOrg.ID,
				FriendlyName: "updated_maintain_org",
				Namespace:    "updated_maintain_namespace",
				Description:  "updated_maintain_description",
				DefaultOrg:   false,
			}
			updatedOrg, err := s.UpdateOrganisation(params)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, updatedOrg.ID, createdOrg.ID)
			odize.AssertEqual(t, updatedOrg.FriendlyName, "updated_maintain_org")
			odize.AssertEqual(t, updatedOrg.Namespace, "updated_maintain_namespace")
			odize.AssertEqual(t, updatedOrg.DefaultOrg, false)
		}).
		Run()
	odize.AssertNoError(t, err)
}
