package watchtower

import (
	"context"
	"testing"
	"time"

	"github.com/code-gorilla-au/odize"
)

func TestService_CreateProduct(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should create a product with valid inputs", func(t *testing.T) {
			org, err := s.CreateOrganisation("test_org", "test_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"web", "api", "microservice"}
			product, err := s.CreateProduct("Test Product", "A test product description", tags, org.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product.Name, "Test Product")
			odize.AssertEqual(t, product.Description, "A test product description")
			odize.AssertEqual(t, len(product.Tags), 3)
			odize.AssertEqual(t, product.Tags[0], "web")
			odize.AssertEqual(t, product.Tags[1], "api")
			odize.AssertEqual(t, product.Tags[2], "microservice")
			odize.AssertFalse(t, product.CreatedAt == time.Time{})
			odize.AssertTrue(t, product.ID > 0)
		}).
		Test("should create a product with empty tags array", func(t *testing.T) {
			org, err := s.CreateOrganisation("empty_tags_org", "empty_tags_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{}
			product, err := s.CreateProduct("Empty Tags Product", "Product with no tags", tags, org.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product.Name, "Empty Tags Product")
			odize.AssertEqual(t, product.Description, "Product with no tags")
			odize.AssertEqual(t, len(product.Tags), 0)
			odize.AssertTrue(t, product.ID > 0)
		}).
		Test("should create a product with single tag", func(t *testing.T) {
			org, err := s.CreateOrganisation("single_tag_org", "single_tag_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"backend"}
			product, err := s.CreateProduct("Single Tag Product", "Product with one tag", tags, org.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product.Name, "Single Tag Product")
			odize.AssertEqual(t, product.Description, "Product with one tag")
			odize.AssertEqual(t, len(product.Tags), 1)
			odize.AssertEqual(t, product.Tags[0], "backend")
			odize.AssertTrue(t, product.ID > 0)
		}).
		Test("should create product even with non-existent organisation ID", func(t *testing.T) {
			tags := []string{"test"}
			product, err := s.CreateProduct("Test Product", "Test description", tags, 99999)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, product.Name, "Test Product")
			odize.AssertEqual(t, product.Description, "Test description")
			odize.AssertTrue(t, product.ID > 0)
		}).
		Test("should create product with empty name", func(t *testing.T) {
			org, err := s.CreateOrganisation("empty_name_org", "empty_name_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"test"}
			product, err := s.CreateProduct("", "Product with empty name", tags, org.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product.Name, "")
			odize.AssertEqual(t, product.Description, "Product with empty name")
			odize.AssertEqual(t, len(product.Tags), 1)
			odize.AssertEqual(t, product.Tags[0], "test")
			odize.AssertTrue(t, product.ID > 0)
		}).
		Test("should create product with empty description", func(t *testing.T) {
			org, err := s.CreateOrganisation("empty_desc_org", "empty_desc_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"test"}
			product, err := s.CreateProduct("Product Name", "", tags, org.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product.Name, "Product Name")
			odize.AssertEqual(t, product.Description, "")
			odize.AssertEqual(t, len(product.Tags), 1)
			odize.AssertEqual(t, product.Tags[0], "test")
			odize.AssertTrue(t, product.ID > 0)
		}).
		Test("should create multiple products for same organisation", func(t *testing.T) {
			org, err := s.CreateOrganisation("multi_prod_org", "multi_prod_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags1 := []string{"frontend", "react"}
			product1, err := s.CreateProduct("Frontend Product", "React frontend", tags1, org.ID)
			odize.AssertNoError(t, err)

			tags2 := []string{"backend", "go"}
			product2, err := s.CreateProduct("Backend Product", "Go backend", tags2, org.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product1.Name, "Frontend Product")
			odize.AssertEqual(t, product2.Name, "Backend Product")
			odize.AssertTrue(t, product1.ID != product2.ID)

			products, err := s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(products) >= 2)
		}).
		Test("should handle complex tags with special characters", func(t *testing.T) {
			org, err := s.CreateOrganisation("complex_tags_org", "complex_tags_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"tag-with-dash", "tag_with_underscore", "tag.with.dots", "tag with spaces"}
			product, err := s.CreateProduct("Complex Tags Product", "Product with complex tags", tags, org.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product.Name, "Complex Tags Product")
			odize.AssertEqual(t, len(product.Tags), 4)
			odize.AssertEqual(t, product.Tags[0], "tag-with-dash")
			odize.AssertEqual(t, product.Tags[1], "tag_with_underscore")
			odize.AssertEqual(t, product.Tags[2], "tag.with.dots")
			odize.AssertEqual(t, product.Tags[3], "tag with spaces")
		}).
		Run()
	odize.AssertNoError(t, err)
}
