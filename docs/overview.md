# Tổng Quan Cache Module

## Giới Thiệu

Cache Module (`go.fork.vn/cache`) là một framework cache enterprise-grade được thiết kế để cung cấp các giải pháp lưu trữ tạm thời hiệu suất cao cho ứng dụng Go. Module này được phát triển với triết lý "performance-first" và "developer-friendly", mang lại sự cân bằng hoàn hảo giữa hiệu suất, độ tin cậy và tính dễ sử dụng.

## Triết Lý Thiết Kế

### 🎯 **Performance-First Architecture**
Cache module được thiết kế từ đầu với mục tiêu đạt hiệu suất tối ưu:

- **Zero-allocation paths**: Tối thiểu allocation trong các hot paths
- **Lock-free operations**: Sử dụng atomic operations khi có thể
- **Batch processing**: Tối ưu hóa throughput với batch operations
- **Efficient serialization**: Hỗ trợ multiple serialization formats

### 🔧 **Modular Design**
Kiến trúc modular cho phép tùy chỉnh và mở rộng dễ dàng:

```go
// Driver interface cho extensibility
type Driver interface {
    Get(ctx context.Context, key string) (interface{}, bool)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    // ... other methods
}
```

### 🔄 **Production-Ready**
Được thiết kế cho môi trường production với:

- **Thread-safe operations**: Hoàn toàn an toàn với concurrent access
- **Graceful degradation**: Xử lý lỗi elegant và fallback mechanisms
- **Comprehensive monitoring**: Built-in metrics và monitoring
- **Memory management**: Automatic cleanup và memory optimization

## Các Thành Phần Chính

### 1. **Cache Manager**
Là component trung tâm quản lý tất cả các driver cache:

```go
type Manager interface {
    // Core operations
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration) error
    
    // Batch operations
    GetMultiple(keys []string) (map[string]interface{}, []string)
    SetMultiple(values map[string]interface{}, ttl time.Duration) error
    
    // Management
    AddDriver(name string, driver driver.Driver)
    SetDefaultDriver(name string) error
}
```

**Tính năng chính:**
- Multi-driver support với seamless switching
- Default driver configuration
- Centralized statistics và monitoring
- Unified API cho tất cả driver types

### 2. **Driver System**
Hệ thống driver linh hoạt hỗ trợ multiple backends:

#### **Memory Driver**
```go
// High-performance in-memory cache
type memoryDriver struct {
    items           sync.Map
    defaultTTL     time.Duration
    cleanupInterval time.Duration
    maxItems       int
}
```

**Đặc điểm:**
- Fastest access times (nanoseconds)
- Automatic cleanup của expired items
- Memory limit protection
- Concurrent-safe operations

#### **File Driver**
```go
// Persistent filesystem-based cache
type fileDriver struct {
    basePath        string
    defaultTTL     time.Duration
    extension      string
    cleanupInterval time.Duration
}
```

**Đặc điểm:**
- Persistent storage across restarts
- File-based TTL management
- Automatic directory management
- Safe concurrent file operations

#### **Redis Driver**
```go
// Distributed cache với Redis backend
type redisDriver struct {
    client        redis.Cmdable
    prefix        string
    defaultTTL   time.Duration
    serializer   func(interface{}) ([]byte, error)
    deserializer func([]byte, interface{}) error
}
```

**Đặc điểm:**
- Distributed caching capabilities
- Multiple serialization formats (JSON, GOB, MessagePack)
- Redis cluster support
- Network-optimized operations

#### **MongoDB Driver**
```go
// Document-based cache với MongoDB
type mongodbDriver struct {
    client     *mongo.Client
    collection *mongo.Collection
    defaultTTL time.Duration
}
```

**Đặc điểm:**
- Document-based storage
- TTL indexes for automatic expiration
- Rich query capabilities
- Horizontal scaling support

### 3. **Configuration System**
Hệ thống cấu hình linh hoạt và type-safe:

```go
type Config struct {
    DefaultDriver string        `yaml:"default_driver"`
    DefaultTTL   int           `yaml:"default_ttl"`
    Prefix       string        `yaml:"prefix"`
    Drivers      DriversConfig `yaml:"drivers"`
}
```

**Features:**
- YAML/JSON configuration support
- Environment variable substitution
- Runtime configuration updates
- Validation và default values

### 4. **Service Provider**
Integration seamless với DI containers:

```go
type ServiceProvider interface {
    di.ServiceProvider
    Requires() []string
    Register(app di.Application)
    Boot(app di.Application)
}
```

**Capabilities:**
- Automatic dependency resolution
- Lifecycle management
- Multiple driver registration
- Error handling và validation

## Use Cases và Patterns

### 1. **Web Application Caching**
```go
// Session caching
manager.Set("session:"+sessionID, sessionData, 30*time.Minute)

// User profile caching
userProfile, found := manager.Get("user:profile:" + userID)
if !found {
    userProfile = fetchUserProfile(userID)
    manager.Set("user:profile:"+userID, userProfile, 1*time.Hour)
}
```

### 2. **API Response Caching**
```go
// Remember pattern cho expensive operations
response, err := manager.Remember("api:report:"+reportID, 15*time.Minute, func() (interface{}, error) {
    return generateExpensiveReport(reportID)
})
```

### 3. **Database Query Caching**
```go
// Query result caching
cacheKey := fmt.Sprintf("query:%x", hash(query))
results, found := manager.Get(cacheKey)
if !found {
    results = database.Query(query)
    manager.Set(cacheKey, results, 5*time.Minute)
}
```

### 4. **Distributed Configuration Caching**
```go
// Configuration caching across services
configs := map[string]interface{}{
    "feature:flags":    featureFlags,
    "service:config":   serviceConfig,
    "rate:limits":      rateLimits,
}
manager.SetMultiple(configs, 1*time.Hour)
```

## Performance Characteristics

### Memory Driver
- **Read Latency**: ~10-50 nanoseconds
- **Write Latency**: ~50-100 nanoseconds
- **Throughput**: >1M ops/second (single core)
- **Memory Overhead**: ~50 bytes/item

### File Driver
- **Read Latency**: ~100-500 microseconds
- **Write Latency**: ~500-2000 microseconds
- **Throughput**: ~10K ops/second
- **Storage Overhead**: ~100 bytes/item

### Redis Driver
- **Read Latency**: ~0.1-1 millisecond
- **Write Latency**: ~0.5-2 milliseconds
- **Throughput**: ~100K ops/second
- **Network Overhead**: ~50-200 bytes/item

### MongoDB Driver
- **Read Latency**: ~1-5 milliseconds
- **Write Latency**: ~2-10 milliseconds
- **Throughput**: ~10K ops/second
- **Storage Overhead**: ~200-500 bytes/item

## Monitoring và Observability

### Built-in Metrics
```go
stats := manager.Stats()
// Returns map[string]interface{} with:
// - hits: int64
// - misses: int64
// - hit_rate: float64
// - total_operations: int64
// - cache_size: int64
// - memory_usage: int64
```

### Health Checks
```go
// Driver health monitoring
health := manager.HealthCheck()
// Returns detailed health status for each driver
```

### Performance Profiling
```go
// Enable performance profiling
manager.EnableProfiling(true)

// Get performance metrics
metrics := manager.GetMetrics()
```

## Concurrency Model

### Thread Safety
Tất cả operations đều thread-safe và có thể được sử dụng concurrent từ multiple goroutines:

```go
// Safe to call from multiple goroutines
go func() {
    manager.Set("key1", value1, ttl)
}()

go func() {
    value, found := manager.Get("key1")
}()
```

### Lock Strategy
- **Memory Driver**: sync.Map cho lock-free reads
- **File Driver**: Per-file locking để avoid contention
- **Redis Driver**: Connection pooling với pipeline support
- **MongoDB Driver**: Session management với connection pooling

## Error Handling

### Graceful Degradation
```go
// Cache failures không ảnh hưởng đến application logic
value, found := manager.Get("key")
if !found {
    // Fallback to primary data source
    value = fetchFromDatabase("key")
}
```

### Error Reporting
```go
// Comprehensive error information
err := manager.Set("key", value, ttl)
if err != nil {
    // Error contains driver-specific details
    log.Printf("Cache error: %v", err)
}
```

## Best Practices

### 1. **Key Design**
```go
// Use hierarchical keys
"user:profile:123"
"session:data:abc"
"api:response:endpoint:/users"
```

### 2. **TTL Strategy**
```go
// Different TTLs for different data types
manager.Set("user:session:"+id, session, 30*time.Minute)
manager.Set("user:profile:"+id, profile, 24*time.Hour)
manager.Set("config:app", config, 7*24*time.Hour)
```

### 3. **Batch Operations**
```go
// Use batch operations cho efficiency
keys := []string{"key1", "key2", "key3"}
results, missed := manager.GetMultiple(keys)
```

### 4. **Memory Management**
```go
// Regular cleanup và monitoring
go func() {
    ticker := time.NewTicker(1*time.Hour)
    for range ticker.C {
        stats := manager.Stats()
        if stats["memory_usage"].(int64) > maxMemoryThreshold {
            manager.Flush()
        }
    }
}()
```

Cache Module cung cấp foundation mạnh mẽ cho caching requirements trong ứng dụng Go modern, với focus vào performance, reliability và developer experience.
