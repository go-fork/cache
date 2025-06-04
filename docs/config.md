# Hướng Dẫn Cấu Hình Cache Module

## Tổng Quan

Cache Module sử dụng hệ thống cấu hình dựa trên YAML/JSON với support cho environment variables và validation tự động. Cấu hình được tổ chức theo cấu trúc phân cấp với global settings và driver-specific configurations.

## Cấu Trúc Cấu Hình

### Global Configuration

```yaml
cache:
  # Driver mặc định được sử dụng khi không chỉ định cụ thể
  default_driver: "memory"
  
  # TTL mặc định cho tất cả cache entries (seconds)
  # 0 = không hết hạn, -1 = sử dụng default của driver
  default_ttl: 3600  # 1 hour
  
  # Prefix cho cache keys để tránh conflicts
  prefix: "cache:"
  
  # Cấu hình cho từng driver
  drivers:
    # ... driver configurations
```

### Config Structure Reference

```go
type Config struct {
    DefaultDriver string        `mapstructure:"default_driver" yaml:"default_driver"`
    DefaultTTL    int          `mapstructure:"default_ttl" yaml:"default_ttl"`
    Prefix        string       `mapstructure:"prefix" yaml:"prefix"`
    Drivers       DriversConfig `mapstructure:"drivers" yaml:"drivers"`
}

type DriversConfig struct {
    Memory  *DriverMemoryConfig  `mapstructure:"memory" yaml:"memory"`
    File    *DriverFileConfig    `mapstructure:"file" yaml:"file"`
    Redis   *DriverRedisConfig   `mapstructure:"redis" yaml:"redis"`
    MongoDB *DriverMongodbConfig `mapstructure:"mongodb" yaml:"mongodb"`
}
```

## Driver Configurations

### 1. Memory Driver Configuration

Memory driver cung cấp high-performance in-memory caching với automatic cleanup.

```yaml
cache:
  drivers:
    memory:
      # Bật/tắt Memory driver
      enabled: true
      
      # TTL mặc định cho memory cache (seconds)
      default_ttl: 3600  # 1 hour
      
      # Interval dọn dẹp expired items (seconds)
      cleanup_interval: 600  # 10 minutes
      
      # Số lượng items tối đa (0 = unlimited)
      max_items: 10000
```

**Configuration Fields:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Kích hoạt memory driver |
| `default_ttl` | int | `3600` | TTL mặc định (seconds) |
| `cleanup_interval` | int | `600` | Interval cleanup (seconds) |
| `max_items` | int | `10000` | Giới hạn số items (0=unlimited) |

**Best Practices:**
- Set `cleanup_interval` từ 1/10 đến 1/6 của `default_ttl`
- Monitor memory usage với `max_items` appropriate
- Sử dụng cho session data, user profiles, temporary computations

### 2. File Driver Configuration

File driver cung cấp persistent cache storage sử dụng filesystem.

```yaml
cache:
  drivers:
    file:
      # Bật/tắt File driver
      enabled: true
      
      # Đường dẫn thư mục lưu cache files
      path: "./storage/cache"
      
      # TTL mặc định cho file cache (seconds)
      default_ttl: 3600  # 1 hour
      
      # Extension cho cache files
      extension: ".cache"
      
      # Interval dọn dẹp expired files (seconds)
      cleanup_interval: 600  # 10 minutes
```

**Configuration Fields:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Kích hoạt file driver |
| `path` | string | `"./storage/cache"` | Thư mục lưu cache files |
| `default_ttl` | int | `3600` | TTL mặc định (seconds) |
| `extension` | string | `".cache"` | File extension |
| `cleanup_interval` | int | `600` | Interval cleanup (seconds) |

**Path Configuration:**
- Sử dụng absolute paths trong production
- Ensure write permissions cho cache directory
- Consider disk space và I/O performance

**Best Practices:**
- Separate cache directory từ application code
- Regular backup strategy cho important cached data
- Monitor disk usage với appropriate cleanup intervals

### 3. Redis Driver Configuration

Redis driver cung cấp distributed caching với Redis backend.

```yaml
cache:
  drivers:
    redis:
      # Bật/tắt Redis driver
      enabled: true
      
      # TTL mặc định cho Redis cache (seconds)
      default_ttl: 3600  # 1 hour
      
      # Serialization format: json, gob, msgpack
      serializer: "json"
```

**Configuration Fields:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Kích hoạt Redis driver |
| `default_ttl` | int | `3600` | TTL mặc định (seconds) |
| `serializer` | string | `"json"` | Serialization format |

**Serialization Options:**

1. **JSON Serializer** (Recommended)
   ```yaml
   serializer: "json"
   ```
   - Human-readable format
   - Cross-language compatibility
   - Good performance cho most use cases

2. **GOB Serializer**
   ```yaml
   serializer: "gob"
   ```
   - Go-native binary format
   - Excellent performance
   - Type safety với Go structs

3. **MessagePack Serializer**
   ```yaml
   serializer: "msgpack"
   ```
   - Compact binary format
   - Cross-language support
   - Balanced performance/size

**Redis Connection:**
Redis driver relies on `go.fork.vn/redis` module configuration:

```yaml
redis:
  default: "main"
  connections:
    main:
      host: "localhost"
      port: 6379
      database: 0
      password: ""
      pool_size: 10
```

### 4. MongoDB Driver Configuration

MongoDB driver cung cấp document-based cache storage với MongoDB backend.

```yaml
cache:
  drivers:
    mongodb:
      # Bật/tắt MongoDB driver
      enabled: true
      
      # Database name cho cache storage
      database: "cache_db"
      
      # Collection name cho cache storage
      collection: "cache_items"
      
      # TTL mặc định cho MongoDB cache (seconds)
      default_ttl: 3600  # 1 hour
      
      # Statistics (readonly, managed by driver)
      hits: 0
      misses: 0
```

**Configuration Fields:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Kích hoạt MongoDB driver |
| `database` | string | `"cache_db"` | Database name |
| `collection` | string | `"cache_items"` | Collection name |
| `default_ttl` | int | `3600` | TTL mặc định (seconds) |
| `hits` | int64 | `0` | Cache hits (readonly) |
| `misses` | int64 | `0` | Cache misses (readonly) |

**MongoDB Connection:**
MongoDB driver relies on `go.fork.vn/mongodb` module configuration:

```yaml
mongodb:
  default: "main"
  connections:
    main:
      uri: "mongodb://localhost:27017"
      database: "app_db"
      timeout: 10s
```

## Environment-Specific Configurations

### Development Environment

```yaml
development:
  cache:
    default_driver: "memory"
    default_ttl: 300  # 5 minutes for faster testing
    prefix: "dev:cache:"
    
    drivers:
      memory:
        enabled: true
        default_ttl: 300
        cleanup_interval: 60  # 1 minute
        max_items: 1000
        
      file:
        enabled: true
        path: "./tmp/cache"
        default_ttl: 300
        cleanup_interval: 60
```

### Testing Environment

```yaml
testing:
  cache:
    default_driver: "memory"
    default_ttl: 60  # 1 minute
    prefix: "test:cache:"
    
    drivers:
      memory:
        enabled: true
        default_ttl: 60
        cleanup_interval: 10  # 10 seconds
        max_items: 100
```

### Production Environment

```yaml
production:
  cache:
    default_driver: "redis"
    default_ttl: 7200  # 2 hours
    prefix: "${APP_NAME:app}:cache:"
    
    drivers:
      redis:
        enabled: true
        default_ttl: 7200
        serializer: "json"
        
      mongodb:
        enabled: true
        database: "${MONGODB_CACHE_DB:cache_db}"
        collection: "${MONGODB_CACHE_COLLECTION:cache_items}"
        default_ttl: 7200
        
      memory:
        enabled: true
        default_ttl: 3600
        cleanup_interval: 300
        max_items: 50000
```

## Environment Variables

Cache module hỗ trợ environment variable substitution trong configuration:

```yaml
cache:
  default_driver: "${CACHE_DEFAULT_DRIVER:memory}"
  default_ttl: "${CACHE_DEFAULT_TTL:3600}"
  prefix: "${CACHE_PREFIX:cache:}"
  
  drivers:
    redis:
      enabled: "${REDIS_CACHE_ENABLED:true}"
      serializer: "${REDIS_SERIALIZER:json}"
      
    mongodb:
      database: "${MONGODB_CACHE_DB:cache_db}"
      collection: "${MONGODB_CACHE_COLLECTION:cache_items}"
```

**Environment Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `CACHE_DEFAULT_DRIVER` | `memory` | Default cache driver |
| `CACHE_DEFAULT_TTL` | `3600` | Default TTL (seconds) |
| `CACHE_PREFIX` | `cache:` | Cache key prefix |
| `REDIS_CACHE_ENABLED` | `true` | Enable Redis driver |
| `REDIS_SERIALIZER` | `json` | Redis serialization format |
| `MONGODB_CACHE_DB` | `cache_db` | MongoDB database |
| `MONGODB_CACHE_COLLECTION` | `cache_items` | MongoDB collection |

## Configuration Loading

### Using Service Provider

```go
import (
    "go.fork.vn/cache"
    "go.fork.vn/di"
)

// Service provider tự động load configuration
provider := cache.NewServiceProvider()
app.Register(provider)
```

### Manual Configuration

```go
import (
    "go.fork.vn/cache/config"
    "go.fork.vn/cache/driver"
    "go.fork.vn/cache"
)

// Load configuration manually
cfg := config.DefaultConfig()
cfg.DefaultDriver = "redis"
cfg.DefaultTTL = 7200

// Create drivers
manager := cache.NewManager()

// Memory driver
if cfg.Drivers.Memory != nil && cfg.Drivers.Memory.Enabled {
    memoryDriver := driver.NewMemoryDriver(*cfg.Drivers.Memory)
    manager.AddDriver("memory", memoryDriver)
}

// Set default driver
manager.SetDefaultDriver(cfg.DefaultDriver)
```

### Configuration Validation

```go
// Validation được thực hiện automatically
func (c *Config) Validate() error {
    if c.DefaultDriver == "" {
        return errors.New("default_driver is required")
    }
    
    if c.DefaultTTL < 0 {
        return errors.New("default_ttl must be >= 0")
    }
    
    // Validate driver configurations
    return c.Drivers.Validate()
}
```

## Helper Methods

Cache configuration cung cấp helper methods để convert time values:

```go
// Global TTL as time.Duration
duration := config.GetDefaultExpiration()

// Driver-specific TTL
memoryTTL := config.Drivers.Memory.GetDefaultExpiration()
fileTTL := config.Drivers.File.GetDefaultExpiration()
redisTTL := config.Drivers.Redis.GetDefaultExpiration()
mongoTTL := config.Drivers.MongoDB.GetDefaultExpiration()

// Cleanup intervals
memoryCleanup := config.Drivers.Memory.GetCleanupInterval()
fileCleanup := config.Drivers.File.GetFileCleanupInterval()
```

## Configuration Examples

### High-Performance Memory-First Setup

```yaml
cache:
  default_driver: "memory"
  default_ttl: 1800  # 30 minutes
  prefix: "hp:cache:"
  
  drivers:
    memory:
      enabled: true
      default_ttl: 1800
      cleanup_interval: 300  # 5 minutes
      max_items: 100000  # 100K items
      
    redis:
      enabled: true
      default_ttl: 3600  # 1 hour fallback
      serializer: "msgpack"  # Compact format
```

### Distributed Cache Setup

```yaml
cache:
  default_driver: "redis"
  default_ttl: 3600
  prefix: "dist:cache:"
  
  drivers:
    redis:
      enabled: true
      default_ttl: 3600
      serializer: "json"
      
    mongodb:
      enabled: true
      database: "distributed_cache"
      collection: "cache_data"
      default_ttl: 7200  # Longer TTL for persistence
```

### Multi-Tier Cache Setup

```yaml
cache:
  default_driver: "memory"
  default_ttl: 1800
  prefix: "multi:cache:"
  
  drivers:
    memory:
      enabled: true
      default_ttl: 900   # 15 minutes (L1)
      cleanup_interval: 180
      max_items: 50000
      
    redis:
      enabled: true
      default_ttl: 3600  # 1 hour (L2)
      serializer: "json"
      
    file:
      enabled: true
      path: "/var/cache/app"
      default_ttl: 7200  # 2 hours (L3)
      cleanup_interval: 1800
```

Configuration system của Cache Module được thiết kế để linh hoạt và dễ maintain, với support đầy đủ cho các deployment scenarios khác nhau từ development đến large-scale production environments.
