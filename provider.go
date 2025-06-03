package cache

import (
	"go.fork.vn/cache/config"
	"go.fork.vn/cache/driver"
	configService "go.fork.vn/config"
	"go.fork.vn/di"
	"go.fork.vn/mongodb"
	"go.fork.vn/redis"
)

// ServiceProvider là interface cho cache service provider.
//
// ServiceProvider định nghĩa các phương thức cần thiết cho một cache service provider
// và kế thừa từ di.ServiceProvider.
type ServiceProvider interface {
	di.ServiceProvider
}

// serviceProvider là service provider cho module cache.
//
// serviceProvider đảm nhận việc đăng ký và khởi tạo các dịch vụ cache vào container DI.
// Nó cung cấp cơ chế để đăng ký cache manager và các driver mặc định vào ứng dụng.
type serviceProvider struct {
	providers []string
}

// NewServiceProvider tạo một cache service provider mới.
//
// Phương thức này khởi tạo một service provider mới cho module cache.
//
// Returns:
//   - ServiceProvider: Đối tượng service provider đã sẵn sàng đăng ký vào ứng dụng
func NewServiceProvider() ServiceProvider {
	return &serviceProvider{}
}

// Requires trả về danh sách các service provider mà provider này phụ thuộc.
//
// Cache provider phụ thuộc vào config provider để đọc cấu hình.
//
// Trả về:
//   - []string: danh sách các service provider khác mà provider này yêu cầu
func (p *serviceProvider) Requires() []string {
	return []string{"config", "redis", "mongodb"}
}

// Register đăng ký các binding vào container.
//
// Phương thức này đăng ký cache manager vào container DI của ứng dụng.
// Nó khởi tạo một cache manager mới và đăng ký nó với khóa "cache".
// Cấu hình sẽ được load từ config manager và các driver được khởi tạo theo cấu hình.
//
// Params:
//   - app: Application instance với DI container và lifecycle management
func (p *serviceProvider) Register(app di.Application) {
	if app == nil {
		return
	}

	c := app.Container()

	// Load cache manager và config manager
	configManager, ok := c.MustMake("config").(configService.Manager)
	if !ok {
		panic("Config manager is not available, please ensure config provider is registered")
	}
	var cfg config.Config
	if err := configManager.UnmarshalKey("cache", &cfg); err != nil {
		panic("Cache config unmarshal error: " + err.Error())
	}

	manager := NewManager()
	// Đăng ký cache service
	c.Instance("cache", manager)
	p.providers = append(p.providers, "cache")

	if cfg.Drivers.Memory != nil && cfg.Drivers.Memory.Enabled {
		// Đăng ký Memory Driver vào cache manager
		memoryDriver := driver.NewMemoryDriver(*cfg.Drivers.Memory)
		manager.AddDriver("memory", memoryDriver)
		c.Instance("cache.memory", memoryDriver)
		p.providers = append(p.providers, "cache.memory")
	}

	if cfg.Drivers.File != nil && cfg.Drivers.File.Enabled {
		// Đăng ký File Driver vào cache manager
		fileDriver, err := driver.NewFileDriver(*cfg.Drivers.File)
		if err != nil {
			panic("Failed to create File driver: " + err.Error())
		}
		manager.AddDriver("file", fileDriver)
		c.Instance("cache.file", fileDriver)
		p.providers = append(p.providers, "cache.file")
	}

	if cfg.Drivers.Redis != nil && cfg.Drivers.Redis.Enabled {
		redisManager := c.MustMake("redis").(redis.Manager)
		if redisManager == nil {
			panic("Redis manager is nil, please ensure Redis provider is registered")
		}
		// Đăng ký Redis Driver vào cache manager
		redisDriver, err := driver.NewRedisDriver(*cfg.Drivers.Redis, redisManager)
		if err != nil {
			panic("Failed to create Redis driver: " + err.Error())
		}
		manager.AddDriver("redis", redisDriver)
		c.Instance("cache.redis", redisDriver)
		p.providers = append(p.providers, "cache.redis")
	}

	if cfg.Drivers.MongoDB != nil && cfg.Drivers.MongoDB.Enabled {
		mongodbManager := c.MustMake("mongodb").(mongodb.Manager)
		if mongodbManager == nil {
			panic("MongoDB manager is nil, please ensure MongoDB provider is registered")
		}

		// Đăng ký MongoDB Driver vào cache manager
		mongodbDriver, err := driver.NewMongoDBDriver(*cfg.Drivers.MongoDB, mongodbManager)
		if err != nil {
			panic("Failed to create MongoDB driver: " + err.Error())
		}
		manager.AddDriver("mongodb", mongodbDriver)
		c.Instance("cache.mongodb", mongodbDriver)
		p.providers = append(p.providers, "cache.mongodb")
	}
}

// Boot được gọi sau khi tất cả các service provider đã được đăng ký.
//
// Phương thức này thực hiện các thao tác khởi tạo cuối cùng sau khi tất cả provider đã đăng ký.
// Hiện tại không cần thực hiện gì đặc biệt vì tất cả cấu hình đã được thực hiện trong Register.
//
// Params:
//   - app: Application instance với DI container và lifecycle management
func (p *serviceProvider) Boot(app di.Application) {
	if app == nil {
		return
	}
}

// Providers trả về danh sách các dịch vụ được cung cấp bởi provider.
//
// Trả về:
//   - []string: danh sách các khóa dịch vụ mà provider này cung cấp
func (p *serviceProvider) Providers() []string {
	return p.providers
}
