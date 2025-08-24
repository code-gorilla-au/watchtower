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
			odize.AssertFalse(t, org.CreatedAt == time.Time{})
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
			odize.AssertFalse(t, defaultOrg.CreatedAt == time.Time{})
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
			odize.AssertFalse(t, fetchedOrg.CreatedAt == time.Time{})
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
