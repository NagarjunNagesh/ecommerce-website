# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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