# Changelog

## [Unreleased]

### Added

### Changed

### Fixed

### Updated

### Removed

## v0.1.1 - 2025-06-04

### Added
- **Comprehensive Documentation**: Thêm bộ tài liệu hoàn chỉnh bằng tiếng Việt với Mermaid diagrams
  - Tài liệu hướng dẫn cấu hình chi tiết (config.md)
  - Tài liệu driver với biểu đồ so sánh hiệu suất (driver.md)
  - Tài liệu tổng quan hệ thống (overview.md, index.md)
  - Tài liệu quản lý cache và service provider (manager.md, provider.md)
  - Biểu đồ kiến trúc hệ thống với Mermaid
  - Biểu đồ so sánh hiệu suất các driver
  - Flowchart cho Remember Pattern và cache operations
- **CI/CD Workflows**: Thêm GitHub Actions workflows cho testing, release, và dependency updates
- **Release Automation**: Scripts tự động cho việc archive releases và tạo release templates
- **Documentation**: Thêm README cho automation scripts trong release management
- **MockServiceProvider**: Triển khai MockServiceProvider sử dụng mockery cho mục đích testing
- **Unit Tests**: Tạo unit tests cho NewServiceProvider, Requires, Providers, Register, và Boot methods
- **Integration Tests**: Thêm integration tests để xác minh lifecycle hoàn chỉnh của cache service provider

### Changed
- **Mock Packages**: Refactor mock packages và cập nhật generated files với formatting và structure nhất quán
- **Provider Error Handling**: Cải thiện error handling trong service provider - thay đổi từ early return thành panic khi DI application là nil

### Fixed
- **Dependencies**: Regenerate go.sum với các dependencies đã cập nhật
- **Redis Driver Tests**: Sửa lỗi trong Redis driver tests liên quan đến mock expectations và serialization
  - Khắc phục vấn đề mock expectations không khớp với byte arrays
  - Sửa lỗi serialization/deserialization trong test scenarios
  - Hoàn thiện tất cả test cases cho Redis driver operations

### Updated
- **Dependencies**: Cập nhật dependencies lên phiên bản mới nhất:
  - `go.fork.vn/config`: v0.1.2 => v0.1.3
  - `go.fork.vn/di`: v0.1.2 => v0.1.3  
  - `go.fork.vn/mongodb`: v0.1.1 => v0.1.2
  - `go.fork.vn/redis`: v0.1.1 => v0.1.2

### Removed
- **Test Files**: Loại bỏ backup test files cho cache service provider tests để dọn dẹp repository

## v0.1.0 - 2025-05-31

### Added
- **Đa dạng driver**: Hỗ trợ bộ nhớ (Memory), tệp tin (File), Redis, MongoDB và khả năng mở rộng driver tùy chỉnh
- **TTL (Time To Live)**: Tự động quản lý thời gian sống cho các mục trong cache
- **Remember pattern**: Hỗ trợ tính toán lười biếng và lưu trữ kết quả trong cache
- **Batch operations**: Thao tác hàng loạt để tối ưu hiệu suất
- **Serialization**: Tự động chuyển đổi giữa cấu trúc dữ liệu Go và định dạng lưu trữ
- **Thread-safe**: An toàn khi truy xuất và cập nhật đồng thời
- **Tích hợp DI**: Dễ dàng tích hợp với Dependency Injection container
- **Extensible**: Dễ dàng mở rộng với driver tùy chỉnh thông qua interface Driver
- Support for graceful error handling during cache invalidation
- Integration with telemetry and monitoring systems
- Multiple cache drivers supported:
  - Memory driver for in-RAM cache storage
  - File driver for filesystem-based cache
  - Redis driver for distributed caching
  - MongoDB driver for document-based cache storage
- Thread-safe cache manager with support for multiple concurrent drivers
- Direct serialization and deserialization for Go structs
- Comprehensive error handling
- Dependency Injection integration through ServiceProvider
- Extensible API for custom cache drivers through Driver interface

### Technical Details
- Initial release as standalone module `go.fork.vn/cache`
- Repository located at `github.com/go-fork/cache`
- Built with Go 1.23.9
- Full test coverage and documentation included

[Unreleased]: github.com/go-fork/cache/compare/v0.1.0...HEAD
[v0.1.0]: github.com/go-fork/cache/releases/tag/v0.1.0
