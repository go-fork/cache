# Changelog

## [Unreleased]

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
- Repository located at `github.com/Fork/cache`
- Built with Go 1.23.9
- Full test coverage and documentation included

[Unreleased]: github.com/go-fork/cache/compare/v0.1.0...HEAD
[v0.1.0]: github.com/go-fork/cache/releases/tag/v0.1.0
