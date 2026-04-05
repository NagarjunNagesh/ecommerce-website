# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

### Unimplemented
- `POST /categories` json field `code` and `name` sanitation where these two json objects create multiple entries because of space. 
  - `{"code":"baby  clothes","name": "Baby Clothes"}`  and `{"code":"baby clothes","name": "Baby Clothes"}` create two entries in DB.

## [0.1.0] - 2026-04-05

### Added
- Update the catalog handler and relevant repositories to support offset pagination.
   - The endpoint should accept query parameters `offset` and `limit`.
   - If `offset` is not provided, default to 0.
   - If `limit` is not provided, default to 10. Maximum limit should be 100. Minimum limit should be 1.
   - The response should include the total number of products available.
   - **GORM Order**: https://gorm.io/docs/query.html#Order
   - **GORM LIMIT and OFFSET**: https://gorm.io/docs/query.html#Limit-Offset
   - *Total Number of products available in the DB?* It is ambiguous? - Fetching the total number of products available in the DB rather than the total products returned after offset and limit

## [0.1.0] - 2026-04-04

### Added
- Implement an endpoint to create new categories at `POST /categories`.
  - This endpoint should accept a JSON body with the category details from the category model and create a new entry in the DB.
  - Provide unit tests for this endpoint.


- Implement the product details endpoint at `GET /catalog/{code}`.
  - This endpoint should return the product details including its variants. Do note that variants without specific price should inherit the price from the product.
  - The product details should include the product's category.
  - Provide unit tests for this endpoint.

### Fixed
- GORM: Translate Errors - https://gorm.io/docs/error_handling.html#Dialect-Translated-Errors
- `POST /categories` - Transfer Duplicate Key identification to Gorm layer and send business logic specific error.
- Extract Constants: Maintainable Unit Tests


## [0.1.0] - 2026-04-03

### Added
- Update the catalog handler and relevant repositories to include the product category in the response.

- Implement the categories endpoint at `/categories`. Created `app/categories` - handler and service classes.
  - This endpoint should return a list of all categories.
  - Provide unit tests for this endpoint.
  - Added the Service Layer and Handler Test Layer.

- Create a new model for Product Categories
  - Added ncessary code to link categories in DB to Repository. 
  - Added the `sql/006-category-data.sql` where the mentioned categories were added.
  - Added the `sql/007-product-category-mapping.sql` modification of products to include relevant category id.

### Changed
- Refactored `products` to add `category_id` as an extra column with Foreign Key Constraint and Indexing for category id based product retrieval from Postgres.
- Refactored `ProductsRepository` to preload the category, added `Category` type to the Products.


## [0.1.0] - 2026-04-02

### Added
- Initial project 

### Changed
- Refactored `CatalogService` to depend on a `IProductsRepository` interface instead of the concrete `*models.ProductsRepository`, following the Dependency Inversion Principle
- Introduced `CatalogService` to separate HTTP handling from business logic

### Fixed
- Refactored `CatalogHandler` to depend on a `ICategoryService` interface instead of the class `CatalogService` , following the Dependency Inversion Principle.
- Added Relevant Logs to the `CatalogService` to better analyse the code.
- Added `CatalogServiceTest` to expand the unit testing to business logic.

[Unreleased]: https://github.com/NagarjunNagesh/ecommerce-website/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/NagarjunNagesh/ecommerce-website/releases/tag/v0.1.0