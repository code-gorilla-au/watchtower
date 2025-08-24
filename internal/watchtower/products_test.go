package watchtower

import (
	"context"
	"fmt"
	"testing"
	"time"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/odize"
)

func TestService_CreateProduct(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	var orgID int64

	group.BeforeAll(func() {
		s = NewService(ctx, _testDB)

		org, err := s.CreateOrganisation("test_org_for_product", "test_org_namespace_for_product", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}
		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should create a product with valid inputs", func(t *testing.T) {

			tags := []string{"web", "api", "microservice"}
			product, err := s.CreateProduct("Test Product", "A test product description", tags, orgID)
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

			tags := []string{}
			product, err := s.CreateProduct("Empty Tags Product", "Product with no tags", tags, orgID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product.Name, "Empty Tags Product")
			odize.AssertEqual(t, product.Description, "Product with no tags")
			odize.AssertEqual(t, len(product.Tags), 0)
			odize.AssertTrue(t, product.ID > 0)
		}).
		Test("should create a product with single tag", func(t *testing.T) {

			tags := []string{"backend"}
			product, err := s.CreateProduct("Single Tag Product", "Product with one tag", tags, orgID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product.Name, "Single Tag Product")
			odize.AssertEqual(t, product.Description, "Product with one tag")
			odize.AssertEqual(t, len(product.Tags), 1)
			odize.AssertEqual(t, product.Tags[0], "backend")
			odize.AssertTrue(t, product.ID > 0)
		}).
		Test("should create product even with non-existent organisation ID", func(t *testing.T) {
			tags := []string{"test"}
			product, err := s.CreateProduct("test non existent org product", "Test description", tags, 99999000000)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, product.Name, "test non existent org product")
			odize.AssertEqual(t, product.Description, "Test description")
			odize.AssertTrue(t, product.ID > 0)
		}).
		Test("should create product with empty name", func(t *testing.T) {

			tags := []string{"test"}
			product, err := s.CreateProduct("", "Product with empty name", tags, orgID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product.Name, "")
			odize.AssertEqual(t, product.Description, "Product with empty name")
			odize.AssertEqual(t, len(product.Tags), 1)
			odize.AssertEqual(t, product.Tags[0], "test")
			odize.AssertTrue(t, product.ID > 0)
		}).
		Test("should create product with empty description", func(t *testing.T) {

			tags := []string{"test"}
			product, err := s.CreateProduct("Product Name", "", tags, orgID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product.Name, "Product Name")
			odize.AssertEqual(t, product.Description, "")
			odize.AssertEqual(t, len(product.Tags), 1)
			odize.AssertEqual(t, product.Tags[0], "test")
			odize.AssertTrue(t, product.ID > 0)
		}).
		Test("should create multiple products for same organisation", func(t *testing.T) {

			tags1 := []string{"frontend", "react"}
			product1, err := s.CreateProduct("Frontend Product", "React frontend", tags1, orgID)
			odize.AssertNoError(t, err)

			tags2 := []string{"backend", "go"}
			product2, err := s.CreateProduct("Backend Product", "Go backend", tags2, orgID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, product1.Name, "Frontend Product")
			odize.AssertEqual(t, product2.Name, "Backend Product")
			odize.AssertTrue(t, product1.ID != product2.ID)

			products, err := s.GetAllProductsForOrganisation(orgID)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(products) >= 2)
		}).
		Test("should handle complex tags with special characters", func(t *testing.T) {

			tags := []string{"tag-with-dash", "tag_with_underscore", "tag.with.dots", "tag with spaces"}
			product, err := s.CreateProduct("Complex Tags Product", "Product with complex tags", tags, orgID)
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

func TestService_GetProductByID(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	var orgID int64

	group.BeforeAll(func() {
		s = NewService(ctx, _testDB)

		org, err := s.CreateOrganisation("test_org_for_get_product", "test_org_namespace_for_get_product", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}
		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should return error when product does not exist", func(t *testing.T) {
			_, err := s.GetProductByID(999)
			odize.AssertError(t, err)
		}).
		Test("should return product when it exists", func(t *testing.T) {

			org, err := s.CreateOrganisation("get_test_org", "get_test_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"web", "api"}
			createdProduct, err := s.CreateProduct("Get Test Product", "Product for get test", tags, org.ID)
			odize.AssertNoError(t, err)

			fetchedProduct, err := s.GetProductByID(createdProduct.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, fetchedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, fetchedProduct.Name, "Get Test Product")
			odize.AssertEqual(t, fetchedProduct.Description, "Product for get test")
			odize.AssertEqual(t, len(fetchedProduct.Tags), 2)
			odize.AssertEqual(t, fetchedProduct.Tags[0], "web")
			odize.AssertEqual(t, fetchedProduct.Tags[1], "api")
			odize.AssertFalse(t, fetchedProduct.CreatedAt == time.Time{})
		}).
		Test("should return correct product when multiple products exist", func(t *testing.T) {
			org, err := s.CreateOrganisation("multi_get_org", "multi_get_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags1 := []string{"frontend"}
			product1, err := s.CreateProduct("First Product", "First product description", tags1, org.ID)
			odize.AssertNoError(t, err)

			tags2 := []string{"backend"}
			product2, err := s.CreateProduct("Second Product", "Second product description", tags2, org.ID)
			odize.AssertNoError(t, err)

			fetchedProduct1, err := s.GetProductByID(product1.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, fetchedProduct1.ID, product1.ID)
			odize.AssertEqual(t, fetchedProduct1.Name, "First Product")
			odize.AssertEqual(t, fetchedProduct1.Description, "First product description")
			odize.AssertEqual(t, len(fetchedProduct1.Tags), 1)
			odize.AssertEqual(t, fetchedProduct1.Tags[0], "frontend")

			fetchedProduct2, err := s.GetProductByID(product2.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, fetchedProduct2.ID, product2.ID)
			odize.AssertEqual(t, fetchedProduct2.Name, "Second Product")
			odize.AssertEqual(t, fetchedProduct2.Description, "Second product description")
			odize.AssertEqual(t, len(fetchedProduct2.Tags), 1)
			odize.AssertEqual(t, fetchedProduct2.Tags[0], "backend")

			odize.AssertTrue(t, fetchedProduct1.ID != fetchedProduct2.ID)
		}).
		Test("should return product with empty tags", func(t *testing.T) {
			org, err := s.CreateOrganisation("empty_tags_get_org", "empty_tags_get_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{}
			createdProduct, err := s.CreateProduct("No Tags Product", "Product without tags", tags, org.ID)
			odize.AssertNoError(t, err)

			fetchedProduct, err := s.GetProductByID(createdProduct.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, fetchedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, fetchedProduct.Name, "No Tags Product")
			odize.AssertEqual(t, fetchedProduct.Description, "Product without tags")
			odize.AssertEqual(t, len(fetchedProduct.Tags), 0)
		}).
		Test("should return product with complex data", func(t *testing.T) {
			org, err := s.CreateOrganisation("complex_get_org", "complex_get_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"tag-with-dash", "tag_with_underscore", "tag.with.dots", "tag with spaces"}
			createdProduct, err := s.CreateProduct("Complex Product", "Product with complex tags and special characters", tags, org.ID)
			odize.AssertNoError(t, err)

			fetchedProduct, err := s.GetProductByID(createdProduct.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, fetchedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, fetchedProduct.Name, "Complex Product")
			odize.AssertEqual(t, fetchedProduct.Description, "Product with complex tags and special characters")
			odize.AssertEqual(t, len(fetchedProduct.Tags), 4)
			odize.AssertEqual(t, fetchedProduct.Tags[0], "tag-with-dash")
			odize.AssertEqual(t, fetchedProduct.Tags[1], "tag_with_underscore")
			odize.AssertEqual(t, fetchedProduct.Tags[2], "tag.with.dots")
			odize.AssertEqual(t, fetchedProduct.Tags[3], "tag with spaces")
		}).
		Test("should return product with empty name and description", func(t *testing.T) {

			tags := []string{"test"}
			createdProduct, err := s.CreateProduct("", "", tags, orgID)
			odize.AssertNoError(t, err)

			fetchedProduct, err := s.GetProductByID(createdProduct.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, fetchedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, fetchedProduct.Name, "")
			odize.AssertEqual(t, fetchedProduct.Description, "")
			odize.AssertEqual(t, len(fetchedProduct.Tags), 1)
			odize.AssertEqual(t, fetchedProduct.Tags[0], "test")
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_GetAllProductsForOrganisation(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should return empty slice when organisation has no products", func(t *testing.T) {
			org, err := s.CreateOrganisation("empty_org", "empty_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			products, err := s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(products), 0)
		}).
		Test("should return single product when organisation has one product", func(t *testing.T) {
			org, err := s.CreateOrganisation("single_prod_org", "single_prod_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"web", "api"}
			createdProduct, err := s.CreateProduct("Single Product", "Single product description", tags, org.ID)
			odize.AssertNoError(t, err)

			products, err := s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(products), 1)
			odize.AssertEqual(t, products[0].ID, createdProduct.ID)
			odize.AssertEqual(t, products[0].Name, "Single Product")
			odize.AssertEqual(t, products[0].Description, "Single product description")
			odize.AssertEqual(t, len(products[0].Tags), 2)
			odize.AssertEqual(t, products[0].Tags[0], "web")
			odize.AssertEqual(t, products[0].Tags[1], "api")
		}).
		Test("should return all products when organisation has multiple products", func(t *testing.T) {
			org, err := s.CreateOrganisation("multi_prod_org", "multi_prod_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags1 := []string{"frontend"}
			product1, err := s.CreateProduct("Frontend Product", "Frontend description", tags1, org.ID)
			odize.AssertNoError(t, err)

			tags2 := []string{"backend"}
			product2, err := s.CreateProduct("Backend Product", "Backend description", tags2, org.ID)
			odize.AssertNoError(t, err)

			tags3 := []string{"mobile"}
			product3, err := s.CreateProduct("Mobile Product", "Mobile description", tags3, org.ID)
			odize.AssertNoError(t, err)

			products, err := s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(products), 3)

			productIDs := make(map[int64]bool)
			for _, p := range products {
				productIDs[p.ID] = true
			}
			odize.AssertTrue(t, productIDs[product1.ID])
			odize.AssertTrue(t, productIDs[product2.ID])
			odize.AssertTrue(t, productIDs[product3.ID])
		}).
		Test("should only return products for specified organisation", func(t *testing.T) {
			org1, err := s.CreateOrganisation("org1", "namespace1", "token1", "description1")
			odize.AssertNoError(t, err)

			org2, err := s.CreateOrganisation("org2", "namespace2", "token2", "description2")
			odize.AssertNoError(t, err)

			tags1 := []string{"org1-tag"}
			product1, err := s.CreateProduct("Org1 Product", "Product for org1", tags1, org1.ID)
			odize.AssertNoError(t, err)

			tags2 := []string{"org2-tag"}
			product2, err := s.CreateProduct("Org2 Product", "Product for org2", tags2, org2.ID)
			odize.AssertNoError(t, err)

			products1, err := s.GetAllProductsForOrganisation(org1.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(products1), 1)
			odize.AssertEqual(t, products1[0].ID, product1.ID)
			odize.AssertEqual(t, products1[0].Name, "Org1 Product")

			products2, err := s.GetAllProductsForOrganisation(org2.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(products2), 1)
			odize.AssertEqual(t, products2[0].ID, product2.ID)
			odize.AssertEqual(t, products2[0].Name, "Org2 Product")
		}).
		Test("should return empty slice for non-existent organisation", func(t *testing.T) {
			products, err := s.GetAllProductsForOrganisation(99999)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(products), 0)
		}).
		Test("should return products with all data fields correctly populated", func(t *testing.T) {
			org, err := s.CreateOrganisation("data_test_org", "data_test_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"tag1", "tag2", "tag3"}
			createdProduct, err := s.CreateProduct("Data Test Product", "Product with full data", tags, org.ID)
			odize.AssertNoError(t, err)

			products, err := s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(products), 1)

			product := products[0]
			odize.AssertEqual(t, product.ID, createdProduct.ID)
			odize.AssertEqual(t, product.Name, "Data Test Product")
			odize.AssertEqual(t, product.Description, "Product with full data")
			odize.AssertEqual(t, len(product.Tags), 3)
			odize.AssertEqual(t, product.Tags[0], "tag1")
			odize.AssertEqual(t, product.Tags[1], "tag2")
			odize.AssertEqual(t, product.Tags[2], "tag3")
			odize.AssertFalse(t, product.CreatedAt == time.Time{})
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_UpdateProduct(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	var orgID int64

	group.BeforeAll(func() {
		s = NewService(ctx, _testDB)

		org, err := s.CreateOrganisation("test_org_for_update_product", "test_org_namespace_for_update_product", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}
		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should return error when trying to update non-existent product", func(t *testing.T) {
			tags := `["updated"]`
			_, err := s.UpdateProduct(999, "Updated Name", &tags)
			odize.AssertError(t, err)
		}).
		Test("should successfully update product with new name and tags", func(t *testing.T) {

			initialTags := []string{"initial", "tag"}
			createdProduct, err := s.CreateProduct("Initial Name", "Initial description", initialTags, orgID)
			odize.AssertNoError(t, err)

			updatedTags := `["updated", "new", "tags"]`
			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "Updated Name", &updatedTags)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, updatedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, updatedProduct.Name, "Updated Name")
			odize.AssertEqual(t, updatedProduct.Description, "Initial description")
			odize.AssertEqual(t, len(updatedProduct.Tags), 3)
			odize.AssertEqual(t, updatedProduct.Tags[0], "updated")
			odize.AssertEqual(t, updatedProduct.Tags[1], "new")
			odize.AssertEqual(t, updatedProduct.Tags[2], "tags")
		}).
		Test("should successfully update product with only name change (nil tags)", func(t *testing.T) {

			initialTags := []string{"keep", "these", "tags"}
			createdProduct, err := s.CreateProduct("Old Name", "Keep description", initialTags, orgID)
			odize.AssertNoError(t, err)

			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "New Name", nil)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, updatedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, updatedProduct.Name, "New Name")
			odize.AssertEqual(t, updatedProduct.Description, "Keep description")
			odize.AssertEqual(t, len(updatedProduct.Tags), 0)
		}).
		Test("should successfully update product with only tags change", func(t *testing.T) {

			initialTags := []string{"old", "tags"}
			createdProduct, err := s.CreateProduct("Keep Name", "Keep description", initialTags, orgID)
			odize.AssertNoError(t, err)

			newTags := `["completely", "different", "tags"]`
			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "Keep Name", &newTags)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, updatedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, updatedProduct.Name, "Keep Name")
			odize.AssertEqual(t, updatedProduct.Description, "Keep description")
			odize.AssertEqual(t, len(updatedProduct.Tags), 3)
			odize.AssertEqual(t, updatedProduct.Tags[0], "completely")
			odize.AssertEqual(t, updatedProduct.Tags[1], "different")
			odize.AssertEqual(t, updatedProduct.Tags[2], "tags")
		}).
		Test("should successfully update product with empty tags", func(t *testing.T) {

			initialTags := []string{"remove", "these"}
			createdProduct, err := s.CreateProduct("Product Name", "Product description", initialTags, orgID)
			odize.AssertNoError(t, err)

			emptyTags := `[]`
			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "Updated Product Name", &emptyTags)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, updatedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, updatedProduct.Name, "Updated Product Name")
			odize.AssertEqual(t, updatedProduct.Description, "Product description")
			odize.AssertEqual(t, len(updatedProduct.Tags), 0)
		}).
		Test("should successfully update product with complex tags containing special characters", func(t *testing.T) {

			initialTags := []string{"simple"}
			createdProduct, err := s.CreateProduct("Test Product", "Test description", initialTags, orgID)
			odize.AssertNoError(t, err)

			complexTags := `["tag-with-dash", "tag_with_underscore", "tag.with.dots", "tag with spaces"]`
			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "complex tags product", &complexTags)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, updatedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, updatedProduct.Name, "complex tags product")
			odize.AssertEqual(t, updatedProduct.Description, "Test description")
			odize.AssertEqual(t, len(updatedProduct.Tags), 4)
			odize.AssertEqual(t, updatedProduct.Tags[0], "tag-with-dash")
			odize.AssertEqual(t, updatedProduct.Tags[1], "tag_with_underscore")
			odize.AssertEqual(t, updatedProduct.Tags[2], "tag.with.dots")
			odize.AssertEqual(t, updatedProduct.Tags[3], "tag with spaces")
		}).
		Test("should verify updated product can be retrieved correctly", func(t *testing.T) {

			initialTags := []string{"before"}
			createdProduct, err := s.CreateProduct("Before Update", "Before description", initialTags, orgID)
			odize.AssertNoError(t, err)

			afterTags := `["after", "update"]`
			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "After Update", &afterTags)
			odize.AssertNoError(t, err)

			retrievedProduct, err := s.GetProductByID(createdProduct.ID)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, retrievedProduct.ID, updatedProduct.ID)
			odize.AssertEqual(t, retrievedProduct.Name, "After Update")
			odize.AssertEqual(t, retrievedProduct.Description, "Before description")
			odize.AssertEqual(t, len(retrievedProduct.Tags), 2)
			odize.AssertEqual(t, retrievedProduct.Tags[0], "after")
			odize.AssertEqual(t, retrievedProduct.Tags[1], "update")
			odize.AssertEqual(t, retrievedProduct.CreatedAt, updatedProduct.CreatedAt)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_DeleteProduct(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	var orgID int64

	group.BeforeAll(func() {
		s = NewService(ctx, _testDB)

		org, err := s.CreateOrganisation("test_org_for_delete_product", "test_org_namespace_for_delete_product", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}
		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should not return error when trying to delete non-existent product", func(t *testing.T) {
			err := s.DeleteProduct(999)
			odize.AssertNoError(t, err)
		}).
		Test("should successfully delete existing product", func(t *testing.T) {

			tags := []string{"delete", "test"}
			createdProduct, err := s.CreateProduct("Delete Test Product", "Product to be deleted", tags, orgID)
			odize.AssertNoError(t, err)

			err = s.DeleteProduct(createdProduct.ID)
			odize.AssertNoError(t, err)
		}).
		Test("should not be able to retrieve deleted product", func(t *testing.T) {

			tags := []string{"verify", "delete"}
			createdProduct, err := s.CreateProduct("Verify Delete Product", "Product to verify deletion", tags, orgID)
			odize.AssertNoError(t, err)

			err = s.DeleteProduct(createdProduct.ID)
			odize.AssertNoError(t, err)

			_, err = s.GetProductByID(createdProduct.ID)
			odize.AssertError(t, err)
		}).
		Test("should successfully delete product with empty tags", func(t *testing.T) {

			tags := []string{}
			createdProduct, err := s.CreateProduct("Empty Tags Delete", "Product with empty tags to delete", tags, orgID)
			odize.AssertNoError(t, err)

			err = s.DeleteProduct(createdProduct.ID)
			odize.AssertNoError(t, err)

			_, err = s.GetProductByID(createdProduct.ID)
			odize.AssertError(t, err)
		}).
		Test("should successfully delete product with complex data", func(t *testing.T) {

			tags := []string{"tag-with-dash", "tag_with_underscore", "tag.with.dots", "tag with spaces"}
			createdProduct, err := s.CreateProduct("Complex Delete Product", "Product with complex tags to delete", tags, orgID)
			odize.AssertNoError(t, err)

			err = s.DeleteProduct(createdProduct.ID)
			odize.AssertNoError(t, err)

			_, err = s.GetProductByID(createdProduct.ID)
			odize.AssertError(t, err)
		}).
		Test("should successfully delete product with empty name and description", func(t *testing.T) {

			tags := []string{"test"}
			createdProduct, err := s.CreateProduct("", "", tags, orgID)
			odize.AssertNoError(t, err)

			err = s.DeleteProduct(createdProduct.ID)
			odize.AssertNoError(t, err)

			_, err = s.GetProductByID(createdProduct.ID)
			odize.AssertError(t, err)
		}).
		Test("should delete multiple products independently", func(t *testing.T) {

			tags1 := []string{"product1"}
			product1, err := s.CreateProduct("Product 1", "First product", tags1, orgID)
			odize.AssertNoError(t, err)

			tags2 := []string{"product2"}
			product2, err := s.CreateProduct("Product 2", "Second product", tags2, orgID)
			odize.AssertNoError(t, err)

			err = s.DeleteProduct(product1.ID)
			odize.AssertNoError(t, err)

			_, err = s.GetProductByID(product1.ID)
			odize.AssertError(t, err)

			fetchedProduct2, err := s.GetProductByID(product2.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, fetchedProduct2.ID, product2.ID)
			odize.AssertEqual(t, fetchedProduct2.Name, "Product 2")
		}).
		Test("should verify product is removed from organisation product list after deletion", func(t *testing.T) {
			org, err := s.CreateOrganisation("delete_verify_org", "delete_verify_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"list", "test"}
			createdProduct, err := s.CreateProduct("List Test Product", "Product for list test", tags, org.ID)
			odize.AssertNoError(t, err)

			products, err := s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(products), 1)
			odize.AssertEqual(t, products[0].ID, createdProduct.ID)

			err = s.DeleteProduct(createdProduct.ID)
			odize.AssertNoError(t, err)

			products, err = s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(products), 0)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_GetProductRepos(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	ctx := context.Background()
	var orgID int64

	group.BeforeAll(func() {
		s = NewService(ctx, _testDB)

		org, err := s.CreateOrganisation("test_org_for_get_repos", "test_org_namespace_for_get_repos", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}
		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB)
	})

	err := group.
		Test("should return empty slice when product has no matching repositories", func(t *testing.T) {
			tags := []string{"no-matching-repos"}
			product, err := s.CreateProduct("No Repos Product", "Product with no matching repositories", tags, orgID)
			odize.AssertNoError(t, err)

			repos, err := s.GetProductRepos(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(repos), 0)
		}).
		Test("should return single repository when product has one matching tag", func(t *testing.T) {
			tags := []string{"single-repo-tag"}
			product, err := s.CreateProduct("Single Repo Product", "Product with one repository", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = s.db.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "single-test-repo",
				Url:   "https://github.com/test/single-test-repo",
				Topic: "single-repo-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			repos, err := s.GetProductRepos(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(repos), 1)
			odize.AssertEqual(t, repos[0].Name, "single-test-repo")
			odize.AssertEqual(t, repos[0].URL, "https://github.com/test/single-test-repo")
			odize.AssertEqual(t, repos[0].Topic, "single-repo-tag")
			odize.AssertEqual(t, repos[0].Owner, "test-owner")
			odize.AssertTrue(t, repos[0].ID > 0)
			odize.AssertFalse(t, repos[0].CreatedAt == time.Time{})
			odize.AssertFalse(t, repos[0].UpdatedAt == time.Time{})
		}).
		Test("should return multiple repositories when product has multiple matching tags", func(t *testing.T) {
			tags := []string{"multi-repo-tag1", "multi-repo-tag2"}
			product, err := s.CreateProduct("Multi Repo Product", "Product with multiple repositories", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = s.db.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "multi-test-repo-1",
				Url:   "https://github.com/test/multi-test-repo-1",
				Topic: "multi-repo-tag1",
				Owner: "test-owner-1",
			})
			odize.AssertNoError(t, err)

			_, err = s.db.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "multi-test-repo-2",
				Url:   "https://github.com/test/multi-test-repo-2",
				Topic: "multi-repo-tag2",
				Owner: "test-owner-2",
			})
			odize.AssertNoError(t, err)

			repos, err := s.GetProductRepos(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(repos), 2)

			repoNames := make(map[string]bool)
			for _, repo := range repos {
				repoNames[repo.Name] = true
			}
			odize.AssertTrue(t, repoNames["multi-test-repo-1"])
			odize.AssertTrue(t, repoNames["multi-test-repo-2"])
		}).
		Test("should return only repositories with matching topics", func(t *testing.T) {
			tags := []string{"matching-tag"}
			product, err := s.CreateProduct("Filtering Product", "Product to test filtering", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = s.db.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "matching-repo",
				Url:   "https://github.com/test/matching-repo",
				Topic: "matching-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = s.db.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "non-matching-repo",
				Url:   "https://github.com/test/non-matching-repo",
				Topic: "different-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			repos, err := s.GetProductRepos(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(repos), 1)
			odize.AssertEqual(t, repos[0].Name, "matching-repo")
			odize.AssertEqual(t, repos[0].Topic, "matching-tag")
		}).
		Test("should handle product with empty tags", func(t *testing.T) {
			tags := []string{}
			product, err := s.CreateProduct("Empty Tags Product", "Product with empty tags", tags, orgID)
			odize.AssertNoError(t, err)

			repos, err := s.GetProductRepos(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(repos), 0)
		}).
		Test("should return empty slice for non-existent product", func(t *testing.T) {
			repos, err := s.GetProductRepos(99999)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(repos), 0)
		}).
		Test("should handle repositories with special characters in names and URLs", func(t *testing.T) {
			tags := []string{"special-chars-tag"}
			product, err := s.CreateProduct("Special Chars Product", "Product with special character repos", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = s.db.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "repo-with-dash_underscore.dot",
				Url:   "https://github.com/test-owner/repo-with-dash_underscore.dot",
				Topic: "special-chars-tag",
				Owner: "test-owner-with-dash",
			})
			odize.AssertNoError(t, err)

			repos, err := s.GetProductRepos(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(repos), 1)
			odize.AssertEqual(t, repos[0].Name, "repo-with-dash_underscore.dot")
			odize.AssertEqual(t, repos[0].URL, "https://github.com/test-owner/repo-with-dash_underscore.dot")
			odize.AssertEqual(t, repos[0].Owner, "test-owner-with-dash")
		}).
		Test("should handle multiple repositories with same topic", func(t *testing.T) {
			tags := []string{"same-topic"}
			product, err := s.CreateProduct("Same Topic Product", "Product with multiple repos having same topic", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = s.db.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "same-topic-repo-1",
				Url:   "https://github.com/test/same-topic-repo-1",
				Topic: "same-topic",
				Owner: "owner-1",
			})
			odize.AssertNoError(t, err)

			_, err = s.db.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "same-topic-repo-2",
				Url:   "https://github.com/test/same-topic-repo-2",
				Topic: "same-topic",
				Owner: "owner-2",
			})
			odize.AssertNoError(t, err)

			repos, err := s.GetProductRepos(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(repos), 2)

			repoNames := make(map[string]string)
			for _, repo := range repos {
				repoNames[repo.Name] = repo.Owner
				odize.AssertEqual(t, repo.Topic, "same-topic")
				odize.AssertTrue(t, repo.ID > 0)
			}
			odize.AssertEqual(t, repoNames["same-topic-repo-1"], "owner-1")
			odize.AssertEqual(t, repoNames["same-topic-repo-2"], "owner-2")
		}).
		Run()
	odize.AssertNoError(t, err)
}
