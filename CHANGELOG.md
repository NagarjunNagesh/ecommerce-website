# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
- - Refactored `ProductsRepository` to preload the category, added `Category` type to the Products.


## [0.1.0] - 2026-04-02

### Added
- Initial project 

### Changed
- Refactored `CatalogService` to depend on a `IProductsRepository` interface instead of the concrete `*models.ProductsRepository`, following the Dependency Inversion Principle
- Introduced `CatalogService` to separate HTTP handling from business logic

### Fixed
-- Refactored `CatalogHandler` to depend on a `ICategoryService` interface instead of the class `CatalogService` , following the Dependency Inversion Principle.
- Added Relevant Logs to the `CatalogService` to better analyse the code.
- Added `CatalogServiceTest` to expand the unit testing to business logic.

[Unreleased]: https://github.com/NagarjunNagesh/ecommerce-website/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/NagarjunNagesh/ecommerce-website/releases/tag/v0.1.0