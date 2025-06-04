# Cache Service Provider

## Tổng Quan

Cache Service Provider là component chính để tích hợp Cache Module vào ứng dụng Go thông qua Dependency Injection (DI) container. Provider này đảm nhận việc đăng ký, cấu hình và khởi tạo tất cả các cache drivers dựa trên configuration được định nghĩa.

## Interface Definition

```go
type ServiceProvider interface {
    di.ServiceProvider
    Requires() []string
    Register(app di.Application)
    Boot(app di.Application)
    Providers() []string
}
```

## Service Provider Lifecycle

### 1. **Dependencies Resolution**

Service Provider khai báo các dependencies cần thiết:

```go
func (p *serviceProvider) Requires() []string {
    return []string{"config", "redis", "mongodb"}
}
```

**Dependencies:**
- **config**: Configuration management service
- **redis**: Redis connection manager (optional)
- **mongodb**: MongoDB connection manager (optional)

### 2. **Registration Phase**

Trong phase này, provider đăng ký tất cả cache services vào DI container:

```go
func (p *serviceProvider) Register(app di.Application) {
    if app == nil {
        panic("di application can not nill.")
    }

    c := app.Container()
    
    // Load configuration
    configManager, ok := c.MustMake("config").(configService.Manager)
    if !ok {
        panic("Config manager is not available, please ensure config provider is registered")
    }
    
    var cfg config.Config
    if err := configManager.UnmarshalKey("cache", &cfg); err != nil {
        panic("Cache config unmarshal error: " + err.Error())
    }

    // Create and register cache manager
    manager := NewManager()
    c.Instance("cache", manager)
    
    // Register individual drivers based on configuration
    // ...
}
```

### 3. **Boot Phase**

Boot phase thực hiện final initialization sau khi tất cả dependencies đã được resolved:

```go
func (p *serviceProvider) Boot(app di.Application) {
    if app == nil {
        panic("di application can not nill.")
    }
    
    // Perform any final initialization tasks
    // Currently no additional boot operations needed
}
```

## Service Registration Details

### Cache Manager Registration

```go
// Tạo và đăng ký cache manager chính
manager := NewManager()
c.Instance("cache", manager)
p.providers = append(p.providers, "cache")
```

**Service Key**: `"cache"`
**Type**: `cache.Manager`
**Description**: Main cache manager interface cho application

### Memory Driver Registration

```go
if cfg.Drivers.Memory != nil && cfg.Drivers.Memory.Enabled {
    memoryDriver := driver.NewMemoryDriver(*cfg.Drivers.Memory)
    manager.AddDriver("memory", memoryDriver)
    c.Instance("cache.memory", memoryDriver)
    p.providers = append(p.providers, "cache.memory")
}
```

**Service Key**: `"cache.memory"`
**Type**: `driver.Driver`
**Condition**: Memory driver enabled trong configuration

### File Driver Registration

```go
if cfg.Drivers.File != nil && cfg.Drivers.File.Enabled {
    fileDriver, err := driver.NewFileDriver(*cfg.Drivers.File)
    if err != nil {
        panic("Failed to create File driver: " + err.Error())
    }
    manager.AddDriver("file", fileDriver)
    c.Instance("cache.file", fileDriver)
    p.providers = append(p.providers, "cache.file")
}
```

**Service Key**: `"cache.file"`
**Type**: `driver.Driver`
**Condition**: File driver enabled trong configuration
**Error Handling**: Panics nếu không thể create file driver

### Redis Driver Registration

```go
if cfg.Drivers.Redis != nil && cfg.Drivers.Redis.Enabled {
    redisManager := c.MustMake("redis").(redis.Manager)
    if redisManager == nil {
        panic("Redis manager is nil, please ensure Redis provider is registered")
    }
    
    redisDriver, err := driver.NewRedisDriver(*cfg.Drivers.Redis, redisManager)
    if err != nil {
        panic("Failed to create Redis driver: " + err.Error())
    }
    manager.AddDriver("redis", redisDriver)
    c.Instance("cache.redis", redisDriver)
    p.providers = append(p.providers, "cache.redis")
}
```

**Service Key**: `"cache.redis"`
**Type**: `driver.Driver`
**Dependencies**: Requires `redis.Manager`
**Condition**: Redis driver enabled trong configuration

### MongoDB Driver Registration

```go
if cfg.Drivers.MongoDB != nil && cfg.Drivers.MongoDB.Enabled {
    mongodbManager := c.MustMake("mongodb").(mongodb.Manager)
    if mongodbManager == nil {
        panic("MongoDB manager is nil, please ensure MongoDB provider is registered")
    }
    
    mongodbDriver, err := driver.NewMongoDBDriver(*cfg.Drivers.MongoDB, mongodbManager)
    if err != nil {
        panic("Failed to create MongoDB driver: " + err.Error())
    }
    manager.AddDriver("mongodb", mongodbDriver)
    c.Instance("cache.mongodb", mongodbDriver)
    p.providers = append(p.providers, "cache.mongodb")
}
```

**Service Key**: `"cache.mongodb"`
**Type**: `driver.Driver`
**Dependencies**: Requires `mongodb.Manager`
**Condition**: MongoDB driver enabled trong configuration

## Usage Examples

### Basic Integration

```go
package main

import (
    "go.fork.vn/cache"
    "go.fork.vn/di"
)

func main() {
    // Create DI application
    app := di.NewApplication()
    
    // Register dependencies first
    app.Register(config.NewServiceProvider())
    app.Register(redis.NewServiceProvider())
    app.Register(mongodb.NewServiceProvider())
    
    // Register cache service provider
    cacheProvider := cache.NewServiceProvider()
    app.Register(cacheProvider)
    
    // Boot application
    app.Boot()
    
    // Use cache manager
    cacheManager := app.Container().MustMake("cache").(cache.Manager)
    
    // Cache operations
    cacheManager.Set("key", "value", 5*time.Minute)
    value, found := cacheManager.Get("key")
}
```

### Direct Driver Access

```go
// Access specific drivers directly
memoryDriver := app.Container().MustMake("cache.memory").(driver.Driver)
redisDriver := app.Container().MustMake("cache.redis").(driver.Driver)

// Use drivers directly với context
ctx := context.Background()
memoryDriver.Set(ctx, "memory_key", "value", 1*time.Hour)
redisDriver.Set(ctx, "redis_key", "value", 1*time.Hour)
```

### Multiple Driver Configuration

```yaml
# config.yaml
cache:
  default_driver: "memory"
  drivers:
    memory:
      enabled: true
      default_ttl: 3600
      max_items: 10000
    
    redis:
      enabled: true
      default_ttl: 7200
      serializer: "json"
    
    file:
      enabled: true
      path: "./storage/cache"
      default_ttl: 86400
```

```go
// Multiple drivers sẽ được registered tự động
providers := cacheProvider.Providers()
// Returns: ["cache", "cache.memory", "cache.redis", "cache.file"]
```

## Error Handling

### Configuration Errors

Service Provider thực hiện extensive validation và error handling:

```go
// Nil application check
if app == nil {
    panic("di application can not nill.")
}

// Config manager availability
configManager, ok := c.MustMake("config").(configService.Manager)
if !ok {
    panic("Config manager is not available, please ensure config provider is registered")
}

// Configuration parsing
if err := configManager.UnmarshalKey("cache", &cfg); err != nil {
    panic("Cache config unmarshal error: " + err.Error())
}
```

### Driver Creation Errors

```go
// File driver creation error
fileDriver, err := driver.NewFileDriver(*cfg.Drivers.File)
if err != nil {
    panic("Failed to create File driver: " + err.Error())
}

// Redis manager dependency check
redisManager := c.MustMake("redis").(redis.Manager)
if redisManager == nil {
    panic("Redis manager is nil, please ensure Redis provider is registered")
}

// Redis driver creation error
redisDriver, err := driver.NewRedisDriver(*cfg.Drivers.Redis, redisManager)
if err != nil {
    panic("Failed to create Redis driver: " + err.Error())
}
```

### Error Recovery Strategies

```go
// Graceful degradation pattern
func setupCacheWithFallback(app di.Application) cache.Manager {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Cache setup failed: %v, using in-memory fallback", r)
            // Setup minimal memory-only cache
        }
    }()
    
    return app.Container().MustMake("cache").(cache.Manager)
}
```

## Advanced Configuration

### Conditional Driver Registration

```go
// Custom provider với conditional logic
type CustomCacheProvider struct {
    *serviceProvider
    environment string
}

func (p *CustomCacheProvider) Register(app di.Application) {
    // Base registration
    p.serviceProvider.Register(app)
    
    // Environment-specific drivers
    switch p.environment {
    case "production":
        // Only enable Redis trong production
        // Disable memory driver để save resources
    case "development":
        // Enable all drivers cho testing
    case "testing":
        // Only memory driver cho fast tests
    }
}
```

### Dynamic Driver Registration

```go
// Runtime driver addition
func AddDynamicDriver(app di.Application, name string, driver driver.Driver) {
    cacheManager := app.Container().MustMake("cache").(cache.Manager)
    cacheManager.AddDriver(name, driver)
    
    // Register trong DI container
    app.Container().Instance("cache."+name, driver)
}
```

### Health Check Integration

```go
// Health check provider extension
func (p *serviceProvider) RegisterHealthChecks(app di.Application) {
    healthManager := app.Container().MustMake("health").(health.Manager)
    cacheManager := app.Container().MustMake("cache").(cache.Manager)
    
    healthManager.Register("cache", func() health.Status {
        // Check cache health
        stats := cacheManager.Stats()
        if stats["status"] == "healthy" {
            return health.StatusUp
        }
        return health.StatusDown
    })
}
```

## Provider Information

### Service List

```go
func (p *serviceProvider) Providers() []string {
    return p.providers
}
```

**Returns:** Danh sách các services đã được registered bởi provider

**Example Output:**
```go
[]string{
    "cache",         // Main cache manager
    "cache.memory",  // Memory driver (if enabled)
    "cache.file",    // File driver (if enabled)
    "cache.redis",   // Redis driver (if enabled)
    "cache.mongodb"  // MongoDB driver (if enabled)
}
```

### Dependency Tree

```
Cache Service Provider
├── config (required)
├── redis (optional, for Redis driver)
├── mongodb (optional, for MongoDB driver)
└── Provides:
    ├── cache (Manager)
    ├── cache.memory (Driver)
    ├── cache.file (Driver)
    ├── cache.redis (Driver)
    └── cache.mongodb (Driver)
```

## Best Practices

### 1. **Dependency Order**
```go
// Correct dependency registration order
app.Register(config.NewServiceProvider())     // First
app.Register(redis.NewServiceProvider())      // Before cache
app.Register(mongodb.NewServiceProvider())    // Before cache
app.Register(cache.NewServiceProvider())      // Last
```

### 2. **Configuration Validation**
```go
// Validate configuration trước khi register
func validateCacheConfig(cfg config.Config) error {
    if cfg.DefaultDriver == "" {
        return errors.New("default_driver must be specified")
    }
    
    // Check if default driver is enabled
    switch cfg.DefaultDriver {
    case "memory":
        if cfg.Drivers.Memory == nil || !cfg.Drivers.Memory.Enabled {
            return errors.New("default driver 'memory' is not enabled")
        }
    case "redis":
        if cfg.Drivers.Redis == nil || !cfg.Drivers.Redis.Enabled {
            return errors.New("default driver 'redis' is not enabled")
        }
    }
    
    return nil
}
```

### 3. **Graceful Degradation**
```go
// Setup cache với fallback strategies
func setupProductionCache(app di.Application) {
    // Try Redis first
    if redisAvailable() {
        setupRedisCache(app)
        return
    }
    
    // Fallback to Memory cache
    log.Warn("Redis unavailable, falling back to memory cache")
    setupMemoryCache(app)
}
```

### 4. **Testing Support**
```go
// Test-specific provider configuration
func setupTestCache() cache.Manager {
    cfg := config.Config{
        DefaultDriver: "memory",
        Drivers: config.DriversConfig{
            Memory: &config.DriverMemoryConfig{
                Enabled:         true,
                DefaultTTL:      60,  // Short TTL cho tests
                CleanupInterval: 10,
                MaxItems:        100,
            },
        },
    }
    
    manager := cache.NewManager()
    memoryDriver := driver.NewMemoryDriver(*cfg.Drivers.Memory)
    manager.AddDriver("memory", memoryDriver)
    manager.SetDefaultDriver("memory")
    
    return manager
}
```

Cache Service Provider cung cấp foundation mạnh mẽ cho cache integration trong Go applications, với comprehensive error handling, flexible configuration, và seamless DI container integration.
