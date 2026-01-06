package watchtower

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"
	"watchtower/internal/database"
	"watchtower/internal/organisations"
	"watchtower/internal/products"

	"watchtower/internal/github"

	"github.com/code-gorilla-au/odize"
)

func TestService_CreateProduct(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service

	ctx := context.Background()

	var orgID int64

	group.BeforeAll(func() {
		s = NewService(ctx, _testDB, _testTxnDB)

		org, err := s.CreateOrganisation("test_org_for_product", "test_org_namespace_for_product", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}

		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
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
			odize.AssertFalse(t, product.CreatedAt.Equal(time.Time{}))
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

			productList, err := s.GetAllProductsForOrganisation(orgID)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(productList) >= 2)
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

	group.BeforeAll(func() {
		s = NewService(ctx, _testDB, _testTxnDB)

		_, err := s.CreateOrganisation("test_org_for_get_product", "test_org_namespace_for_get_product", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}

		odize.AssertNoError(t, err)

	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
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
			odize.AssertFalse(t, fetchedProduct.CreatedAt.Equal(time.Time{}))
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
		Run()
	odize.AssertNoError(t, err)
}

func TestService_GetAllProductsForOrganisation(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service

	ctx := context.Background()

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
	})

	err := group.
		Test("should return empty slice when organisation has no products", func(t *testing.T) {
			org, err := s.CreateOrganisation("empty_org", "empty_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			productList, err := s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(productList), 0)
		}).
		Test("should return single product when organisation has one product", func(t *testing.T) {
			org, err := s.CreateOrganisation("single_prod_org", "single_prod_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"web", "api"}
			createdProduct, err := s.CreateProduct("Single Product", "Single product description", tags, org.ID)
			odize.AssertNoError(t, err)

			productList, err := s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(productList), 1)
			odize.AssertEqual(t, productList[0].ID, createdProduct.ID)
			odize.AssertEqual(t, productList[0].Name, "Single Product")
			odize.AssertEqual(t, productList[0].Description, "Single product description")
			odize.AssertEqual(t, len(productList[0].Tags), 2)
			odize.AssertEqual(t, productList[0].Tags[0], "web")
			odize.AssertEqual(t, productList[0].Tags[1], "api")
		}).
		Test("should return all products when organisation has multiple products", func(t *testing.T) {
			org, err := s.CreateOrganisation("multi_prod_org", "multi_prod_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags1 := []string{"frontend"}
			product1, err := s.CreateProduct("Frontend Product1", "Frontend description", tags1, org.ID)
			odize.AssertNoError(t, err)

			tags2 := []string{"backend"}
			product2, err := s.CreateProduct("Backend Product2", "Backend description", tags2, org.ID)
			odize.AssertNoError(t, err)

			tags3 := []string{"mobile"}
			product3, err := s.CreateProduct("Mobile Product3", "Mobile description", tags3, org.ID)
			odize.AssertNoError(t, err)

			productList, err := s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(productList), 3)

			productIDs := make(map[int64]bool)
			for _, p := range productList {
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

			productList1, err := s.GetAllProductsForOrganisation(org1.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(productList1), 1)
			odize.AssertEqual(t, productList1[0].ID, product1.ID)
			odize.AssertEqual(t, productList1[0].Name, "Org1 Product")

			products2, err := s.GetAllProductsForOrganisation(org2.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(products2), 1)
			odize.AssertEqual(t, products2[0].ID, product2.ID)
			odize.AssertEqual(t, products2[0].Name, "Org2 Product")
		}).
		Test("should return empty slice for non-existent organisation", func(t *testing.T) {
			productList, err := s.GetAllProductsForOrganisation(99999)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(productList), 0)
		}).
		Test("should return products with all data fields correctly populated", func(t *testing.T) {
			org, err := s.CreateOrganisation("data_test_org", "data_test_namespace", "token", "test description")
			odize.AssertNoError(t, err)

			tags := []string{"tag1", "tag2", "tag3"}
			createdProduct, err := s.CreateProduct("Data Test Product", "Product with full data", tags, org.ID)
			odize.AssertNoError(t, err)

			productList, err := s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(productList), 1)

			product := productList[0]
			odize.AssertEqual(t, product.ID, createdProduct.ID)
			odize.AssertEqual(t, product.Name, "Data Test Product")
			odize.AssertEqual(t, product.Description, "Product with full data")
			odize.AssertEqual(t, len(product.Tags), 3)
			odize.AssertEqual(t, product.Tags[0], "tag1")
			odize.AssertEqual(t, product.Tags[1], "tag2")
			odize.AssertEqual(t, product.Tags[2], "tag3")
			odize.AssertFalse(t, product.CreatedAt.Equal(time.Time{}))
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
		s = NewService(ctx, _testDB, _testTxnDB)

		org, err := s.CreateOrganisation("test_org_for_update_product", "test_org_namespace_for_update_product", "token", "test description")
		if err != nil {
			fmt.Println("create org error", err)
		}

		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
	})

	err := group.
		Test("should return error when trying to update non-existent product", func(t *testing.T) {
			tags := []string{"updated", "test"}
			_, err := s.UpdateProduct(999, "Updated Name", "updated description", tags)
			odize.AssertError(t, err)
		}).
		Test("should successfully update product with new name and tags", func(t *testing.T) {
			initialTags := []string{"initial", "tag"}
			createdProduct, err := s.CreateProduct("Initial Name", "Initial description", initialTags, orgID)
			odize.AssertNoError(t, err)

			updatedTags := []string{"updated", "new", "tags"}
			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "Updated Name", "updated description", updatedTags)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, updatedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, updatedProduct.Name, "Updated Name")
			odize.AssertEqual(t, updatedProduct.Description, "updated description")
			odize.AssertEqual(t, len(updatedProduct.Tags), 3)
			odize.AssertEqual(t, updatedProduct.Tags[0], "updated")
			odize.AssertEqual(t, updatedProduct.Tags[1], "new")
			odize.AssertEqual(t, updatedProduct.Tags[2], "tags")
		}).
		Test("should successfully update product with only name change (nil tags)", func(t *testing.T) {
			initialTags := []string{"keep", "these", "tags"}
			createdProduct, err := s.CreateProduct("Old Name", "Keep description", initialTags, orgID)
			odize.AssertNoError(t, err)

			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "New Name", "Keep description", nil)
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

			newTags := []string{"completely", "different", "tags"}
			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "Keep Name", "Keep description", newTags)
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
			createdProduct, err := s.CreateProduct("Product Name empty tags", "Product description", initialTags, orgID)
			odize.AssertNoError(t, err)

			emptyTags := []string{}
			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "Updated Product Name empty tags", "Product description", emptyTags)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, updatedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, updatedProduct.Name, "Updated Product Name empty tags")
			odize.AssertEqual(t, updatedProduct.Description, "Product description")
			odize.AssertEqual(t, len(updatedProduct.Tags), 0)
		}).
		Test("should successfully update product with complex tags containing special characters", func(t *testing.T) {
			initialTags := []string{"simple"}
			createdProduct, err := s.CreateProduct("Test Product special characters", "Test description", initialTags, orgID)
			odize.AssertNoError(t, err)

			complexTags := []string{"tag-with-dash", "tag_with_underscore", "tag.with.dots", "tag with spaces"}
			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "complex tags product special characters", "Test description", complexTags)
			odize.AssertNoError(t, err)

			odize.AssertEqual(t, updatedProduct.ID, createdProduct.ID)
			odize.AssertEqual(t, updatedProduct.Name, "complex tags product special characters")
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

			afterTags := []string{"after", "update"}
			updatedProduct, err := s.UpdateProduct(createdProduct.ID, "After Update", "Before description", afterTags)
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
		s = NewService(ctx, _testDB, _testTxnDB)

		org, err := s.CreateOrganisation("test_org_for_delete_product", "test_org_namespace_for_delete_product", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}

		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
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

			productList, err := s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(productList), 1)
			odize.AssertEqual(t, productList[0].ID, createdProduct.ID)

			err = s.DeleteProduct(createdProduct.ID)
			odize.AssertNoError(t, err)

			productList, err = s.GetAllProductsForOrganisation(org.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(productList), 0)
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
		s = NewService(ctx, _testDB, _testTxnDB)

		org, err := s.CreateOrganisation("test_org_for_get_repos", "test_org_namespace_for_get_repos", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}

		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
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

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
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
			odize.AssertFalse(t, product.CreatedAt.Equal(time.Time{}))
			odize.AssertFalse(t, repos[0].UpdatedAt.Equal(time.Time{}))
		}).
		Test("should return multiple repositories when product has multiple matching tags", func(t *testing.T) {
			tags := []string{"multi-repo-tag1", "multi-repo-tag2"}
			product, err := s.CreateProduct("Multi Repo Product", "Product with multiple repositories", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "multi-test-repo-1",
				Url:   "https://github.com/test/multi-test-repo-1",
				Topic: "multi-repo-tag1",
				Owner: "test-owner-1",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
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

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "matching-repo",
				Url:   "https://github.com/test/matching-repo",
				Topic: "matching-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
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
			product, err := s.CreateProduct("Empty Tags Product1", "Product with empty tags", tags, orgID)
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

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
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

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "same-topic-repo-1",
				Url:   "https://github.com/test/same-topic-repo-1",
				Topic: "same-topic",
				Owner: "owner-1",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
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

func TestService_GetProductPullRequests(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service

	ctx := context.Background()

	var orgID int64

	group.BeforeAll(func() {
		s = NewService(ctx, _testDB, _testTxnDB)

		org, err := s.CreateOrganisation("test_org_for_get_prs", "test_org_namespace_for_get_prs", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}

		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
	})

	err := group.
		Test("should return empty slice when product has no matching pull requests", func(t *testing.T) {
			tags := []string{"no-matching-prs"}
			product, err := s.CreateProduct("No PRs Product", "Product with no matching pull requests", tags, orgID)
			odize.AssertNoError(t, err)

			prs, err := s.GetProductPullRequests(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(prs), 0)
		}).
		Test("should return single pull request when product has one matching repository", func(t *testing.T) {
			tags := []string{"single-pr-tag"}
			product, err := s.CreateProduct("Single PR Product", "Product with one pull request", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "single-pr-repo",
				Url:   "https://github.com/test/single-pr-repo",
				Topic: "single-pr-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "pr-external-1",
				Title:          "Single Test PR",
				RepositoryName: "single-pr-repo",
				Url:            "https://github.com/test/single-pr-repo/pull/1",
				State:          "OPEN",
				Author:         "test-author",
				MergedAt: sql.NullInt64{
					Valid: true,
					Int64: 1610000000,
				},
			})
			odize.AssertNoError(t, err)

			prs, err := s.GetProductPullRequests(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(prs), 1)
			odize.AssertEqual(t, prs[0].ExternalID, "pr-external-1")
			odize.AssertEqual(t, prs[0].Title, "Single Test PR")
			odize.AssertEqual(t, prs[0].RepositoryName, "single-pr-repo")
			odize.AssertEqual(t, prs[0].URL, "https://github.com/test/single-pr-repo/pull/1")
			odize.AssertEqual(t, prs[0].State, "OPEN")
			odize.AssertEqual(t, prs[0].Author, "test-author")
			odize.AssertTrue(t, prs[0].ID > 0)
			odize.AssertFalse(t, prs[0].CreatedAt.Equal(time.Time{}))
			odize.AssertFalse(t, prs[0].UpdatedAt.Equal(time.Time{}))
		}).
		Test("should return multiple pull requests when product has multiple matching repositories", func(t *testing.T) {
			tags := []string{"multi-pr-tag1", "multi-pr-tag2"}
			product, err := s.CreateProduct("Multi PR Product", "Product with multiple pull requests", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "multi-pr-repo-1",
				Url:   "https://github.com/test/multi-pr-repo-1",
				Topic: "multi-pr-tag1",
				Owner: "test-owner-1",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "multi-pr-repo-2",
				Url:   "https://github.com/test/multi-pr-repo-2",
				Topic: "multi-pr-tag2",
				Owner: "test-owner-2",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "pr-external-2",
				Title:          "Multi Test PR 1",
				RepositoryName: "multi-pr-repo-1",
				Url:            "https://github.com/test/multi-pr-repo-1/pull/1",
				State:          "OPEN",
				Author:         "test-author-1",
				MergedAt: sql.NullInt64{
					Valid: true,
					Int64: 1610000000,
				},
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "pr-external-3",
				Title:          "Multi Test PR 2",
				RepositoryName: "multi-pr-repo-2",
				Url:            "https://github.com/test/multi-pr-repo-2/pull/1",
				State:          "OPEN",
				Author:         "test-author-2",
				MergedAt: sql.NullInt64{
					Int64: 0,
					Valid: false,
				},
			})
			odize.AssertNoError(t, err)

			prs, err := s.GetProductPullRequests(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(prs), 2)

			prTitles := make(map[string]bool)
			for _, pr := range prs {
				prTitles[pr.Title] = true
				odize.AssertEqual(t, pr.State, "OPEN")
			}
			odize.AssertTrue(t, prTitles["Multi Test PR 1"])
			odize.AssertTrue(t, prTitles["Multi Test PR 2"])
		}).
		Test("should return only OPEN pull requests and filter out closed ones", func(t *testing.T) {
			tags := []string{"state-filter-tag"}
			product, err := s.CreateProduct("State Filter Product", "Product to test state filtering", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "state-filter-repo",
				Url:   "https://github.com/test/state-filter-repo",
				Topic: "state-filter-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "pr-external-open",
				Title:          "Open PR",
				RepositoryName: "state-filter-repo",
				Url:            "https://github.com/test/state-filter-repo/pull/1",
				State:          "OPEN",
				Author:         "test-author",
				MergedAt: sql.NullInt64{
					Valid: true,
					Int64: 1610000000,
				},
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "pr-external-closed",
				Title:          "Closed PR",
				RepositoryName: "state-filter-repo",
				Url:            "https://github.com/test/state-filter-repo/pull/2",
				State:          "CLOSED",
				Author:         "test-author",
				MergedAt: sql.NullInt64{
					Valid: true,
					Int64: 1610000000,
				},
			})
			odize.AssertNoError(t, err)

			prs, err := s.GetProductPullRequests(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(prs), 1)
			odize.AssertEqual(t, prs[0].Title, "Open PR")
			odize.AssertEqual(t, prs[0].State, "OPEN")
		}).
		Test("should handle product with empty tags", func(t *testing.T) {
			tags := []string{}
			product, err := s.CreateProduct("Empty Tags PR Product", "Product with empty tags", tags, orgID)
			odize.AssertNoError(t, err)

			prs, err := s.GetProductPullRequests(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(prs), 0)
		}).
		Test("should return empty slice for non-existent product", func(t *testing.T) {
			prs, err := s.GetProductPullRequests(99999)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(prs), 0)
		}).
		Test("should handle pull requests with special characters and complete data", func(t *testing.T) {
			tags := []string{"special-chars-pr-tag"}
			product, err := s.CreateProduct("Special Chars PR Product", "Product with special character PRs", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "special-chars-pr-repo",
				Url:   "https://github.com/test-owner/special-chars-pr-repo",
				Topic: "special-chars-pr-tag",
				Owner: "test-owner-with-dash",
			})
			odize.AssertNoError(t, err)

			mergedAt := sql.NullInt64{Int64: time.Now().Unix(), Valid: true}
			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "pr-with-special-chars",
				Title:          "Fix: bug with special-chars & symbols",
				RepositoryName: "special-chars-pr-repo",
				Url:            "https://github.com/test-owner/special-chars-pr-repo/pull/1",
				State:          "OPEN",
				Author:         "author-with_underscore",
				MergedAt:       mergedAt,
			})
			odize.AssertNoError(t, err)

			prs, err := s.GetProductPullRequests(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(prs), 1)
			odize.AssertEqual(t, prs[0].Title, "Fix: bug with special-chars & symbols")
			odize.AssertEqual(t, prs[0].Author, "author-with_underscore")
			odize.AssertEqual(t, prs[0].ExternalID, "pr-with-special-chars")
		}).
		Test("should handle multiple pull requests from same repository", func(t *testing.T) {
			tags := []string{"same-repo-pr-tag"}
			product, err := s.CreateProduct("Same Repo PR Product", "Product with multiple PRs from same repo", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "same-repo-multiple-prs",
				Url:   "https://github.com/test/same-repo-multiple-prs",
				Topic: "same-repo-pr-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "pr-same-repo-1",
				Title:          "First PR from same repo",
				RepositoryName: "same-repo-multiple-prs",
				Url:            "https://github.com/test/same-repo-multiple-prs/pull/1",
				State:          "OPEN",
				Author:         "author-1",
				MergedAt: sql.NullInt64{
					Valid: true,
					Int64: 1610000000,
				},
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "pr-same-repo-2",
				Title:          "Second PR from same repo",
				RepositoryName: "same-repo-multiple-prs",
				Url:            "https://github.com/test/same-repo-multiple-prs/pull/2",
				State:          "OPEN",
				Author:         "author-2",
				MergedAt: sql.NullInt64{
					Valid: true,
					Int64: 1610000000,
				},
			})
			odize.AssertNoError(t, err)

			prs, err := s.GetProductPullRequests(product.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(prs), 2)

			prTitles := make(map[string]string)
			for _, pr := range prs {
				prTitles[pr.Title] = pr.Author
				odize.AssertEqual(t, pr.RepositoryName, "same-repo-multiple-prs")
				odize.AssertEqual(t, pr.State, "OPEN")
				odize.AssertTrue(t, pr.ID > 0)
			}
			odize.AssertEqual(t, prTitles["First PR from same repo"], "author-1")
			odize.AssertEqual(t, prTitles["Second PR from same repo"], "author-2")
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_GetPullRequestByOrganisation(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service

	ctx := context.Background()

	var orgID int64

	group.BeforeAll(func() {
		s = NewService(ctx, _testDB, _testTxnDB)

		org, err := s.CreateOrganisation("test_org_for_pr_by_org", "test_org_namespace_for_pr_by_org", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}

		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
	})

	err := group.
		Test("should return error when organisation does not exist", func(t *testing.T) {
			prs, err := s.GetPullRequestByOrganisation(99999)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(prs), 0)
		}).
		Test("should return empty slice when organisation has no pull requests", func(t *testing.T) {
			prs, err := s.GetPullRequestByOrganisation(orgID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(prs), 0)
		}).
		Test("should return pull requests for organisation with products", func(t *testing.T) {
			tags := []string{"org-pr-test-tag"}
			_, err := s.CreateProduct("Org PR Test Product", "Product for org PR testing", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "org-pr-test-repo",
				Url:   "https://github.com/test/org-pr-test-repo",
				Topic: "org-pr-test-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "org-pr-external-1",
				Title:          "First Org PR",
				RepositoryName: "org-pr-test-repo",
				Url:            "https://github.com/test/org-pr-test-repo/pull/1",
				State:          "OPEN",
				Author:         "test-author-1",
				MergedAt: sql.NullInt64{
					Valid: true,
					Int64: 1610000000,
				},
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "org-pr-external-2",
				Title:          "Second Org PR",
				RepositoryName: "org-pr-test-repo",
				Url:            "https://github.com/test/org-pr-test-repo/pull/2",
				State:          "OPEN",
				Author:         "test-author-2",
				MergedAt: sql.NullInt64{
					Valid: true,
					Int64: 1610000000,
				},
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "org-pr-external-closed",
				Title:          "Closed Org PR",
				RepositoryName: "org-pr-test-repo",
				Url:            "https://github.com/test/org-pr-test-repo/pull/3",
				State:          "CLOSED",
				Author:         "test-author-3",
				MergedAt: sql.NullInt64{
					Int64: 0,
					Valid: false,
				},
			})
			odize.AssertNoError(t, err)

			prs, err := s.GetPullRequestByOrganisation(orgID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(prs), 2)

			prTitles := make(map[string]string)
			for _, pr := range prs {
				prTitles[pr.Title] = pr.Author
				odize.AssertEqual(t, pr.State, "OPEN")
				odize.AssertTrue(t, pr.ID > 0)
			}
			odize.AssertEqual(t, prTitles["First Org PR"], "test-author-1")
			odize.AssertEqual(t, prTitles["Second Org PR"], "test-author-2")
		}).
		Test("should only return OPEN pull requests for organisation", func(t *testing.T) {
			tags := []string{"org-pr-state-test-tag"}
			_, err := s.CreateProduct("Org PR State Test Product", "Product for org PR state testing", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "org-pr-state-test-repo",
				Url:   "https://github.com/test/org-pr-state-test-repo",
				Topic: "org-pr-state-test-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			// Create pull requests with different states
			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "org-pr-open-state",
				Title:          "Open State PR",
				RepositoryName: "org-pr-state-test-repo",
				Url:            "https://github.com/test/org-pr-state-test-repo/pull/1",
				State:          "OPEN",
				Author:         "state-test-author",
				MergedAt: sql.NullInt64{
					Valid: true,
					Int64: 1610000000,
				},
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "org-pr-closed-state",
				Title:          "Closed State PR",
				RepositoryName: "org-pr-state-test-repo",
				Url:            "https://github.com/test/org-pr-state-test-repo/pull/2",
				State:          "CLOSED",
				Author:         "state-test-author",
				MergedAt: sql.NullInt64{
					Int64: 0,
					Valid: false,
				},
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "org-pr-merged-state",
				Title:          "Merged State PR",
				RepositoryName: "org-pr-state-test-repo",
				Url:            "https://github.com/test/org-pr-state-test-repo/pull/3",
				State:          "MERGED",
				Author:         "state-test-author",
				MergedAt:       sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
			})
			odize.AssertNoError(t, err)

			prs, err := s.GetPullRequestByOrganisation(orgID)
			odize.AssertNoError(t, err)

			// Should only return OPEN PRs, including from previous test
			openPRsFound := 0
			for _, pr := range prs {
				if pr.Title == "Open State PR" {
					openPRsFound++
					odize.AssertEqual(t, pr.State, "OPEN")
				}
				odize.AssertEqual(t, pr.State, "OPEN")
			}
			odize.AssertTrue(t, openPRsFound == 1)
		}).
		Test("should handle organisation with multiple products and repositories", func(t *testing.T) {
			tags1 := []string{"multi-prod-tag-1"}
			_, err := s.CreateProduct("Multi Prod 1", "First product for multi test", tags1, orgID)
			odize.AssertNoError(t, err)

			tags2 := []string{"multi-prod-tag-2"}
			_, err = s.CreateProduct("Multi Prod 2", "Second product for multi test", tags2, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "multi-prod-repo-1",
				Url:   "https://github.com/test/multi-prod-repo-1",
				Topic: "multi-prod-tag-1",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "multi-prod-repo-2",
				Url:   "https://github.com/test/multi-prod-repo-2",
				Topic: "multi-prod-tag-2",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "multi-prod-pr-1",
				Title:          "PR from Product 1",
				RepositoryName: "multi-prod-repo-1",
				Url:            "https://github.com/test/multi-prod-repo-1/pull/1",
				State:          "OPEN",
				Author:         "multi-test-author-1",
				MergedAt: sql.NullInt64{
					Int64: 1610000000,
					Valid: true,
				},
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreatePullRequest(ctx, database.CreatePullRequestParams{
				ExternalID:     "multi-prod-pr-2",
				Title:          "PR from Product 2",
				RepositoryName: "multi-prod-repo-2",
				Url:            "https://github.com/test/multi-prod-repo-2/pull/1",
				State:          "OPEN",
				Author:         "multi-test-author-2",
				MergedAt: sql.NullInt64{
					Int64: 0,
					Valid: false,
				},
			})
			odize.AssertNoError(t, err)

			prs, err := s.GetPullRequestByOrganisation(orgID)
			odize.AssertNoError(t, err)

			// Should return PRs from both products, plus any from previous tests
			multiProdPRsFound := 0
			for _, pr := range prs {
				if pr.Title == "PR from Product 1" || pr.Title == "PR from Product 2" {
					multiProdPRsFound++
					odize.AssertEqual(t, pr.State, "OPEN")
				}
			}
			odize.AssertTrue(t, multiProdPRsFound == 2)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_GetSecurityByProductID(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service

	ctx := context.Background()

	var orgID int64

	var productID int64

	group.BeforeAll(func() {
		s = NewService(ctx, _testDB, _testTxnDB)

		org, err := s.CreateOrganisation("test_org_for_security", "test_org_namespace_for_security", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}

		odize.AssertNoError(t, err)

		orgID = org.ID

		tags := []string{"security-test-tag"}

		product, err := s.CreateProduct("Security Test Product", "Product for security testing", tags, orgID)
		if err != nil {
			fmt.Print("create product error", err)
		}

		odize.AssertNoError(t, err)

		productID = product.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
	})

	err := group.
		Test("should return empty slice for non-existent product", func(t *testing.T) {
			securities, err := s.GetSecurityByProductID(99999)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(securities), 0)
		}).
		Test("should return empty slice when product has no security data", func(t *testing.T) {
			securities, err := s.GetSecurityByProductID(productID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(securities), 0)
		}).
		Test("should return security data for product with repositories", func(t *testing.T) {
			_, err := _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "security-test-repo",
				Url:   "https://github.com/test/security-test-repo",
				Topic: "security-test-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "security-external-1",
				RepositoryName: "security-test-repo",
				PackageName:    "test-package-1",
				State:          "OPEN",
				Severity:       "HIGH",
				PatchedVersion: "1.2.3",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "security-external-2",
				RepositoryName: "security-test-repo",
				PackageName:    "test-package-2",
				State:          "OPEN",
				Severity:       "MEDIUM",
				PatchedVersion: "2.3.4",
			})
			odize.AssertNoError(t, err)

			securities, err := s.GetSecurityByProductID(productID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(securities), 2)

			securityMap := make(map[string]products.SecurityDTO)
			for _, sec := range securities {
				securityMap[sec.ExternalID] = sec
				odize.AssertEqual(t, sec.State, "OPEN")
				odize.AssertEqual(t, sec.RepositoryName, "security-test-repo")
				odize.AssertTrue(t, sec.ID > 0)
			}

			sec1 := securityMap["security-external-1"]
			odize.AssertEqual(t, sec1.PackageName, "test-package-1")
			odize.AssertEqual(t, sec1.Severity, "HIGH")
			odize.AssertEqual(t, sec1.PatchedVersion, "1.2.3")

			sec2 := securityMap["security-external-2"]
			odize.AssertEqual(t, sec2.PackageName, "test-package-2")
			odize.AssertEqual(t, sec2.Severity, "MEDIUM")
			odize.AssertEqual(t, sec2.PatchedVersion, "2.3.4")
		}).
		Test("should only return OPEN security items", func(t *testing.T) {
			_, err := _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "security-state-test-repo",
				Url:   "https://github.com/test/security-state-test-repo",
				Topic: "security-test-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "security-open",
				RepositoryName: "security-state-test-repo",
				PackageName:    "open-package",
				State:          "OPEN",
				Severity:       "CRITICAL",
				PatchedVersion: "3.4.5",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "security-closed",
				RepositoryName: "security-state-test-repo",
				PackageName:    "closed-package",
				State:          "CLOSED",
				Severity:       "HIGH",
				PatchedVersion: "4.5.6",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "security-fixed",
				RepositoryName: "security-state-test-repo",
				PackageName:    "fixed-package",
				State:          "FIXED",
				Severity:       "LOW",
				PatchedVersion: "5.6.7",
			})
			odize.AssertNoError(t, err)

			securities, err := s.GetSecurityByProductID(productID)
			odize.AssertNoError(t, err)

			openSecuritiesFound := 0
			for _, sec := range securities {
				if sec.ExternalID == "security-open" {
					openSecuritiesFound++
					odize.AssertEqual(t, sec.State, "OPEN")
					odize.AssertEqual(t, sec.PackageName, "open-package")
					odize.AssertEqual(t, sec.Severity, "CRITICAL")
				}
				odize.AssertEqual(t, sec.State, "OPEN")
			}
			odize.AssertTrue(t, openSecuritiesFound == 1)
		}).
		Test("should handle multiple repositories with security data", func(t *testing.T) {
			_, err := _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "multi-security-repo-1",
				Url:   "https://github.com/test/multi-security-repo-1",
				Topic: "security-test-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "multi-security-repo-2",
				Url:   "https://github.com/test/multi-security-repo-2",
				Topic: "security-test-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "multi-security-1",
				RepositoryName: "multi-security-repo-1",
				PackageName:    "multi-package-1",
				State:          "OPEN",
				Severity:       "HIGH",
				PatchedVersion: "6.7.8",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "multi-security-2",
				RepositoryName: "multi-security-repo-2",
				PackageName:    "multi-package-2",
				State:          "OPEN",
				Severity:       "MEDIUM",
				PatchedVersion: "7.8.9",
			})
			odize.AssertNoError(t, err)

			securities, err := s.GetSecurityByProductID(productID)
			odize.AssertNoError(t, err)

			multiSecuritiesFound := 0
			repoNames := make(map[string]bool)
			for _, sec := range securities {
				if sec.ExternalID == "multi-security-1" || sec.ExternalID == "multi-security-2" {
					multiSecuritiesFound++
					repoNames[sec.RepositoryName] = true
					odize.AssertEqual(t, sec.State, "OPEN")
				}
			}
			odize.AssertTrue(t, multiSecuritiesFound == 2)
			odize.AssertTrue(t, repoNames["multi-security-repo-1"])
			odize.AssertTrue(t, repoNames["multi-security-repo-2"])
		}).
		Test("should handle security items with special characters and edge cases", func(t *testing.T) {
			_, err := _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "special-chars-security-repo",
				Url:   "https://github.com/test/special-chars-security-repo",
				Topic: "security-test-tag",
				Owner: "test-owner-with-dash",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "security-with-special-chars",
				RepositoryName: "special-chars-security-repo",
				PackageName:    "@scope/package-name",
				State:          "OPEN",
				Severity:       "CRITICAL",
				PatchedVersion: "1.0.0-beta.1",
			})
			odize.AssertNoError(t, err)

			securities, err := s.GetSecurityByProductID(productID)
			odize.AssertNoError(t, err)

			specialSecurityFound := false
			for _, sec := range securities {
				if sec.ExternalID == "security-with-special-chars" {
					specialSecurityFound = true
					odize.AssertEqual(t, sec.PackageName, "@scope/package-name")
					odize.AssertEqual(t, sec.Severity, "CRITICAL")
					odize.AssertEqual(t, sec.PatchedVersion, "1.0.0-beta.1")
					odize.AssertEqual(t, sec.State, "OPEN")
				}
			}
			odize.AssertTrue(t, specialSecurityFound)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_GetSecurityByOrganisation(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service

	ctx := context.Background()

	var orgID int64

	group.BeforeAll(func() {
		s = NewService(ctx, _testDB, _testTxnDB)

		org, err := s.CreateOrganisation("test_org_for_security_by_org", "test_org_namespace_for_security_by_org", "token", "test description")
		if err != nil {
			fmt.Print("create org error", err)
		}

		odize.AssertNoError(t, err)

		orgID = org.ID
	})

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
	})

	err := group.
		Test("should return empty slice for non-existent organisation", func(t *testing.T) {
			securities, err := s.GetSecurityByOrganisation(99999)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(securities), 0)
		}).
		Test("should return empty slice when organisation has no security data", func(t *testing.T) {
			securities, err := s.GetSecurityByOrganisation(orgID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(securities), 0)
		}).
		Test("should return security data for organisation with products and repositories", func(t *testing.T) {
			tags := []string{"org-security-test-tag"}
			_, err := s.CreateProduct("Org Security Test Product", "Product for org security testing", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "org-security-test-repo",
				Url:   "https://github.com/test/org-security-test-repo",
				Topic: "org-security-test-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "org-security-external-1",
				RepositoryName: "org-security-test-repo",
				PackageName:    "org-test-package-1",
				State:          "OPEN",
				Severity:       "HIGH",
				PatchedVersion: "1.2.3",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "org-security-external-2",
				RepositoryName: "org-security-test-repo",
				PackageName:    "org-test-package-2",
				State:          "OPEN",
				Severity:       "MEDIUM",
				PatchedVersion: "2.3.4",
			})
			odize.AssertNoError(t, err)

			securities, err := s.GetSecurityByOrganisation(orgID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, len(securities), 2)

			securityMap := make(map[string]products.SecurityDTO)
			for _, sec := range securities {
				securityMap[sec.ExternalID] = sec
				odize.AssertEqual(t, sec.State, "OPEN")
				odize.AssertEqual(t, sec.RepositoryName, "org-security-test-repo")
				odize.AssertTrue(t, sec.ID > 0)
			}

			sec1 := securityMap["org-security-external-1"]
			odize.AssertEqual(t, sec1.PackageName, "org-test-package-1")
			odize.AssertEqual(t, sec1.Severity, "HIGH")
			odize.AssertEqual(t, sec1.PatchedVersion, "1.2.3")

			sec2 := securityMap["org-security-external-2"]
			odize.AssertEqual(t, sec2.PackageName, "org-test-package-2")
			odize.AssertEqual(t, sec2.Severity, "MEDIUM")
			odize.AssertEqual(t, sec2.PatchedVersion, "2.3.4")
		}).
		Test("should only return OPEN security items for organisation", func(t *testing.T) {
			tags := []string{"org-security-state-test-tag"}
			_, err := s.CreateProduct("Org Security State Test Product", "Product for org security state testing", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "org-security-state-test-repo",
				Url:   "https://github.com/test/org-security-state-test-repo",
				Topic: "org-security-state-test-tag",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "org-security-open",
				RepositoryName: "org-security-state-test-repo",
				PackageName:    "open-package",
				State:          "OPEN",
				Severity:       "CRITICAL",
				PatchedVersion: "3.4.5",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "org-security-closed",
				RepositoryName: "org-security-state-test-repo",
				PackageName:    "closed-package",
				State:          "CLOSED",
				Severity:       "HIGH",
				PatchedVersion: "4.5.6",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "org-security-fixed",
				RepositoryName: "org-security-state-test-repo",
				PackageName:    "fixed-package",
				State:          "FIXED",
				Severity:       "LOW",
				PatchedVersion: "5.6.7",
			})
			odize.AssertNoError(t, err)

			securities, err := s.GetSecurityByOrganisation(orgID)
			odize.AssertNoError(t, err)

			openSecuritiesFound := 0
			for _, sec := range securities {
				if sec.ExternalID == "org-security-open" {
					openSecuritiesFound++
					odize.AssertEqual(t, sec.State, "OPEN")
					odize.AssertEqual(t, sec.PackageName, "open-package")
					odize.AssertEqual(t, sec.Severity, "CRITICAL")
				}
				odize.AssertEqual(t, sec.State, "OPEN")
			}
			odize.AssertTrue(t, openSecuritiesFound == 1)
		}).
		Test("should handle organisation with multiple products and repositories", func(t *testing.T) {
			tags1 := []string{"org-multi-security-tag-1"}
			_, err := s.CreateProduct("Org Multi Security Prod 1", "First product for multi security test", tags1, orgID)
			odize.AssertNoError(t, err)

			tags2 := []string{"org-multi-security-tag-2"}
			_, err = s.CreateProduct("Org Multi Security Prod 2", "Second product for multi security test", tags2, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "org-multi-security-repo-1",
				Url:   "https://github.com/test/org-multi-security-repo-1",
				Topic: "org-multi-security-tag-1",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "org-multi-security-repo-2",
				Url:   "https://github.com/test/org-multi-security-repo-2",
				Topic: "org-multi-security-tag-2",
				Owner: "test-owner",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "org-multi-security-1",
				RepositoryName: "org-multi-security-repo-1",
				PackageName:    "multi-package-1",
				State:          "OPEN",
				Severity:       "HIGH",
				PatchedVersion: "6.7.8",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "org-multi-security-2",
				RepositoryName: "org-multi-security-repo-2",
				PackageName:    "multi-package-2",
				State:          "OPEN",
				Severity:       "MEDIUM",
				PatchedVersion: "7.8.9",
			})
			odize.AssertNoError(t, err)

			securities, err := s.GetSecurityByOrganisation(orgID)
			odize.AssertNoError(t, err)

			multiSecuritiesFound := 0
			repoNames := make(map[string]bool)
			for _, sec := range securities {
				if sec.ExternalID == "org-multi-security-1" || sec.ExternalID == "org-multi-security-2" {
					multiSecuritiesFound++
					repoNames[sec.RepositoryName] = true
					odize.AssertEqual(t, sec.State, "OPEN")
				}
			}
			odize.AssertTrue(t, multiSecuritiesFound == 2)
			odize.AssertTrue(t, repoNames["org-multi-security-repo-1"])
			odize.AssertTrue(t, repoNames["org-multi-security-repo-2"])
		}).
		Test("should handle security items with special characters and edge cases", func(t *testing.T) {
			tags := []string{"org-special-chars-security-tag"}
			_, err := s.CreateProduct("Org Special Chars Security Product", "Product with special character security", tags, orgID)
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateRepo(ctx, database.CreateRepoParams{
				Name:  "org-special-chars-security-repo",
				Url:   "https://github.com/test/org-special-chars-security-repo",
				Topic: "org-special-chars-security-tag",
				Owner: "test-owner-with-dash",
			})
			odize.AssertNoError(t, err)

			_, err = _testDB.CreateSecurity(ctx, database.CreateSecurityParams{
				ExternalID:     "org-security-with-special-chars",
				RepositoryName: "org-special-chars-security-repo",
				PackageName:    "@scope/package-name",
				State:          "OPEN",
				Severity:       "CRITICAL",
				PatchedVersion: "1.0.0-beta.1",
			})
			odize.AssertNoError(t, err)

			securities, err := s.GetSecurityByOrganisation(orgID)
			odize.AssertNoError(t, err)

			specialSecurityFound := false
			for _, sec := range securities {
				if sec.ExternalID == "org-security-with-special-chars" {
					specialSecurityFound = true
					odize.AssertEqual(t, sec.PackageName, "@scope/package-name")
					odize.AssertEqual(t, sec.Severity, "CRITICAL")
					odize.AssertEqual(t, sec.PatchedVersion, "1.0.0-beta.1")
					odize.AssertEqual(t, sec.State, "OPEN")
				}
			}
			odize.AssertTrue(t, specialSecurityFound)
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_SyncOrgs(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service
	var ctx context.Context
	var ghMock *ghClientMock
	var org1ID, org2ID int64

	group.BeforeAll(func() {
		ctx = context.Background()
	})

	group.BeforeEach(func() {
		ghMock = &ghClientMock{
			SearchReposFunc: func(owner string, topic string, token string) (github.QuerySearch[github.Repository], error) {
				return github.QuerySearch[github.Repository]{
					Data: github.QueryData[github.Repository]{
						Search: github.Search[github.Repository]{
							PageInfo: github.PageInfo{},
							Edges: []github.Node[github.Repository]{
								{
									Node: github.Repository{
										Url:  "",
										Name: "flash",
										Owner: github.Owner{
											Login: "gordon",
										},
										VulnerabilityAlerts: github.RootNode[github.VulnerabilityAlerts]{},
										PullRequests:        github.RootNode[github.PullRequest]{},
									},
								},
							},
						},
					},
				}, nil
			},
			GetRepoDetailsFunc: func(owner, repo, token string) (github.QueryRepository, error) {
				return github.QueryRepository{
					Data: github.RepositoryData{
						Repository: github.Repository{
							Url:  "",
							Name: "flash",
							Owner: github.Owner{
								Login: "gordon",
							},
							VulnerabilityAlerts: github.RootNode[github.VulnerabilityAlerts]{
								PageInfo: github.PageInfo{},
								Nodes: []github.VulnerabilityAlerts{
									{
										State:  "OPEN",
										ID:     "1",
										Number: 1,
										SecurityVulnerability: github.SecurityVulnerability{
											Package: github.Package{
												Name: "1",
											},
											Advisory: github.Advisory{
												Severity: "22",
											},
											FirstPatchedVersion: github.FirstPatchedVersion{
												Identifier: "33",
											},
											UpdatedAt: time.Time{},
										},
										CreatedAt: time.Time{},
										FixedAt:   nil,
									},
								},
							},
							PullRequests: github.RootNode[github.PullRequest]{
								Nodes: []github.PullRequest{
									{
										ID:        "333",
										Title:     "444",
										State:     "555",
										CreatedAt: time.Time{},
										MergedAt:  nil,
										Permalink: "",
										Author: github.Author{
											Login: "333",
										},
									},
								},
							},
						},
					},
				}, nil
			},
		}
		s = &Service{
			ctx:      ctx,
			ghClient: ghMock,
			orgSvc: organisations.New(_testDB, _testTxnDB, func(tx *sql.Tx) organisations.OrgStore {
				return _testDB.WithTx(tx)
			}),
			productSvc: products.New(_testDB),
		}

		// Create test organizations with unique namespaces to avoid conflicts
		timestamp := time.Now().UnixNano()
		org1, err := s.CreateOrganisation(fmt.Sprintf("test_org_1_sync_%d", timestamp), fmt.Sprintf("test_namespace_1_%d", timestamp), "token1", "Test org 1")
		odize.AssertNoError(t, err)
		org1ID = org1.ID

		org2, err := s.CreateOrganisation(fmt.Sprintf("test_org_2_sync_%d", timestamp), fmt.Sprintf("test_namespace_2_%d", timestamp), "token2", "Test org 2")
		odize.AssertNoError(t, err)
		org2ID = org2.ID

		// Create products for each organization
		tags1 := []string{"sync-tag-1"}
		_, err = s.CreateProduct(fmt.Sprintf("sync_product_1_sync_%d", timestamp), "Product for sync test 1", tags1, org1ID)
		odize.AssertNoError(t, err)

		tags2 := []string{"sync-tag-2"}
		_, err = s.CreateProduct(fmt.Sprintf("sync_product_2_sync_%d", timestamp), "Product for sync test 2", tags2, org2ID)
		odize.AssertNoError(t, err)
	})

	err := group.
		Test("should handle no organizations to sync when all are recent", func(t *testing.T) {
			// Since we just created the orgs, they should be recent and not eligible for sync
			err := s.SyncOrgs()
			odize.AssertNoError(t, err)

			// Should not have called SearchRepos since orgs are too recent
			calls := ghMock.SearchReposCalls()
			odize.AssertEqual(t, len(calls), 0)
		}).
		Test("should return error when database operations fail", func(t *testing.T) {
			// Create a service with a cancelled context to simulate database errors
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel() // Cancel immediately to simulate connection issues

			brokenService := &Service{
				ctx:        cancelCtx,
				ghClient:   ghMock,
				orgSvc:     organisations.New(_testDB, nil, nil),
				productSvc: products.New(_testDB),
			}

			err := brokenService.SyncOrgs()
			// This should return an error due to the cancelled context
			odize.AssertError(t, err)
		}).
		Test("should handle organization with no products gracefully", func(t *testing.T) {

			timestamp := time.Now().UnixNano()
			emptyOrg, err := s.CreateOrganisation(fmt.Sprintf("empty_org_sync_%d", timestamp), fmt.Sprintf("empty_namespace_%d", timestamp), "empty_token", "Empty org")
			odize.AssertNoError(t, err)

			err = s.SyncOrg(emptyOrg.ID)
			odize.AssertNoError(t, err)
		}).
		Test("should successfully sync organization with products", func(t *testing.T) {
			err := s.SyncOrg(org1ID)
			odize.AssertNoError(t, err)

			calls := ghMock.SearchReposCalls()
			odize.AssertTrue(t, len(calls) > 0)
		}).
		Test("should call SearchRepos with correct parameters", func(t *testing.T) {
			// Reset call tracking by creating new mock
			ghMock = &ghClientMock{
				SearchReposFunc: func(owner string, topic string, token string) (github.QuerySearch[github.Repository], error) {
					return github.QuerySearch[github.Repository]{}, nil
				},
				GetRepoDetailsFunc: func(owner, repo, token string) (github.QueryRepository, error) {
					return github.QueryRepository{
						QueryModel: github.QueryModel{},
						Data:       github.RepositoryData{},
					}, nil
				},
			}
			s.ghClient = ghMock

			err := s.SyncOrg(org1ID)
			odize.AssertNoError(t, err)

			calls := ghMock.SearchReposCalls()
			odize.AssertTrue(t, len(calls) > 0)

			firstCall := calls[0]
			odize.AssertEqual(t, firstCall.Topic, "sync-tag-1")
			odize.AssertEqual(t, firstCall.Token, "token1")

			odize.AssertTrue(t, firstCall.Owner != "")
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestService_CreateOrganisation(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var s *Service

	ctx := context.Background()

	group.BeforeEach(func() {
		s = NewService(ctx, _testDB, _testTxnDB)
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
		s = NewService(ctx, _testDB, _testTxnDB)
	})

	err := group.
		Test("should return error when no default organisation exists", func(t *testing.T) {
			_ = _testDB.SetOrgsDefaultFalse(ctx)
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
		s = NewService(ctx, _testDB, _testTxnDB)
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
		s = NewService(ctx, _testDB, _testTxnDB)
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
		s = NewService(ctx, _testDB, _testTxnDB)
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
		s = NewService(ctx, _testDB, _testTxnDB)
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
		s = NewService(ctx, _testDB, _testTxnDB)
	})

	err := group.
		Test("should return error when trying to update non-existent organisation", func(t *testing.T) {
			params := organisations.UpdateOrgParams{
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

			params := organisations.UpdateOrgParams{
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

			params := organisations.UpdateOrgParams{
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

			params := organisations.UpdateOrgParams{
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

			params := organisations.UpdateOrgParams{
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
