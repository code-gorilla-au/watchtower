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
