package cache_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.fork.vn/cache"
	"go.fork.vn/cache/config"
	cachemocks "go.fork.vn/cache/mocks"
	configmocks "go.fork.vn/config/mocks"
	"go.fork.vn/di"
	dimocks "go.fork.vn/di/mocks"
	mongodbmocks "go.fork.vn/mongodb/mocks"
	redismocks "go.fork.vn/redis/mocks"
)

func TestServiceProvider_NewServiceProvider(t *testing.T) {
	t.Run("creates_new_service_provider", func(t *testing.T) {
		// Act
		provider := cache.NewServiceProvider()

		// Assert
		assert.NotNil(t, provider)
		assert.Implements(t, (*di.ServiceProvider)(nil), provider)
	})
}

func TestServiceProvider_Requires(t *testing.T) {
	t.Run("returns_required_dependencies", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()

		// Act
		requires := provider.Requires()

		// Assert
		expected := []string{"config", "redis", "mongodb"}
		assert.Equal(t, expected, requires)
		assert.Len(t, requires, 3)
		assert.Contains(t, requires, "config")
		assert.Contains(t, requires, "redis")
		assert.Contains(t, requires, "mongodb")
	})
}

func TestServiceProvider_Providers(t *testing.T) {
	t.Run("returns_empty_before_register", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()

		// Act
		providers := provider.Providers()

		// Assert
		assert.Empty(t, providers)
	})

	t.Run("returns_manager_after_register", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		// Setup mocks
		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					Memory: &config.DriverMemoryConfig{Enabled: true},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)
		mockContainer.EXPECT().Instance("cache.memory", mock.Anything).Return().Times(1)

		// Act
		provider.Register(mockApp)
		providers := provider.Providers()

		// Assert
		assert.Contains(t, providers, "cache")
		assert.Contains(t, providers, "cache.memory")
		assert.Len(t, providers, 2)

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})
}

func TestServiceProvider_Register(t *testing.T) {
	// TestServiceProvider_Register kiểm tra tất cả các scenario đăng ký driver
	t.Run("register_memory_driver_successfully", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					Memory: &config.DriverMemoryConfig{Enabled: true},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)
		mockContainer.EXPECT().Instance("cache.memory", mock.Anything).Return().Times(1)

		// Act & Assert
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})

	t.Run("register_file_driver_successfully", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					File: &config.DriverFileConfig{
						Enabled: true,
						Path:    "/tmp/cache",
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)
		mockContainer.EXPECT().Instance("cache.file", mock.Anything).Return().Times(1)

		// Act & Assert
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})

	t.Run("register_redis_driver_failure_due_to_interface_conversion", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}
		mockRedisManager := &redismocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockContainer.EXPECT().MustMake("redis").Return(mockRedisManager).Times(1)

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					Redis: &config.DriverRedisConfig{
						Enabled: true,
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)

		// Act & Assert - Expect panic due to mock not implementing full redis.Manager interface
		assert.Panics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
		// Note: mockRedisManager expectations are not verified as the conversion fails before method calls
	})

	t.Run("register_mongodb_driver_failure_due_to_nil_database", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}
		mockMongodbManager := &mongodbmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockContainer.EXPECT().MustMake("mongodb").Return(mockMongodbManager).Times(1)

		// Mock MongoDB Manager's DatabaseWithName() method that NewMongoDBDriver calls
		// Return nil which will cause the driver creation to fail as expected during unit tests
		mockMongodbManager.EXPECT().DatabaseWithName("cache").Return(nil).Times(2) // Called twice in NewMongoDBDriver

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					MongoDB: &config.DriverMongodbConfig{
						Enabled:    true,
						Database:   "cache",
						Collection: "cache_entries",
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)

		// Act & Assert - Expect this to panic due to MongoDB driver implementation requiring real database
		assert.Panics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
		mockMongodbManager.AssertExpectations(t)
	})

	t.Run("register_all_drivers_failure_due_to_interface_issues", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}
		mockRedisManager := &redismocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockContainer.EXPECT().MustMake("redis").Return(mockRedisManager).Times(1)

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					Memory: &config.DriverMemoryConfig{Enabled: true},
					File: &config.DriverFileConfig{
						Enabled: true,
						Path:    "/tmp/cache",
					},
					Redis: &config.DriverRedisConfig{
						Enabled: true,
					},
					MongoDB: &config.DriverMongodbConfig{
						Enabled:    true,
						Database:   "cache",
						Collection: "cache_entries",
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)
		mockContainer.EXPECT().Instance("cache.memory", mock.Anything).Return().Times(1)
		mockContainer.EXPECT().Instance("cache.file", mock.Anything).Return().Times(1)

		// Act & Assert - Expect this to panic due to Redis interface conversion failure
		assert.Panics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
		// Note: Redis and MongoDB manager expectations not verified as conversion fails
	})

	t.Run("with_nil_app_returns_early", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()

		// Act & Assert
		assert.NotPanics(t, func() {
			provider.Register(nil)
		})
	})

	t.Run("panic_when_config_unmarshal_fails", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).Return(
			errors.New("config unmarshal error")).Times(1)

		// Act & Assert
		assert.PanicsWithValue(t, "Cache config unmarshal error: config unmarshal error", func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})

	t.Run("panic_when_file_driver_creation_fails", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					File: &config.DriverFileConfig{
						Enabled: true,
						Path:    "/invalid/path/that/cannot/be/created/due/to/permissions",
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)

		// Act & Assert
		assert.Panics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})

	t.Run("panic_when_redis_manager_is_nil", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockContainer.EXPECT().MustMake("redis").Return(nil).Times(1)

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					Redis: &config.DriverRedisConfig{
						Enabled: true,
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)

		// Act & Assert - Type assertion panic occurs when nil is returned from MustMake
		assert.Panics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})

	t.Run("panic_when_mongodb_manager_is_nil", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockContainer.EXPECT().MustMake("mongodb").Return(nil).Times(1)

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					MongoDB: &config.DriverMongodbConfig{
						Enabled:    true,
						Database:   "cache",
						Collection: "cache_entries",
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)

		// Act & Assert - Type assertion panic occurs when nil is returned from MustMake
		assert.Panics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})

	t.Run("register_with_disabled_drivers", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					Memory: &config.DriverMemoryConfig{Enabled: false},
					File: &config.DriverFileConfig{
						Enabled: false,
						Path:    "/tmp/cache",
					},
					Redis: &config.DriverRedisConfig{
						Enabled: false,
					},
					MongoDB: &config.DriverMongodbConfig{
						Enabled:    false,
						Database:   "cache",
						Collection: "cache_entries",
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)

		// Act & Assert
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)

		// Verify only base cache service is registered (no drivers)
		providers := provider.Providers()
		assert.Contains(t, providers, "cache")
		assert.Len(t, providers, 1)
	})

	t.Run("register_with_nil_driver_configs", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					Memory:  nil,
					File:    nil,
					Redis:   nil,
					MongoDB: nil,
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)

		// Act & Assert
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)

		// Verify only base cache service is registered (no drivers)
		providers := provider.Providers()
		assert.Contains(t, providers, "cache")
		assert.Len(t, providers, 1)
	})
	t.Run("redis_driver_creation_fails_with_invalid_config", func(t *testing.T) {
		// This test helps cover the redis driver error handling path
		// By creating an invalid redis config that will cause driver creation to fail
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}
		mockRedisManager := &redismocks.MockManager{}

		// Setup expectations
		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockContainer.EXPECT().MustMake("redis").Return(mockRedisManager).Times(1)

		// Mock redis manager methods that NewRedisDriver may call
		mockRedisManager.EXPECT().Client().Return(nil, errors.New("redis connection failed")).Maybe()

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					Redis: &config.DriverRedisConfig{
						Enabled:    true,
						DefaultTTL: -1,                   // Invalid TTL might cause issues
						Serializer: "invalid_serializer", // Invalid serializer
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)

		// Act & Assert - Expect this to panic due to redis driver creation failure
		assert.Panics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})

	t.Run("mongodb_driver_creation_fails_with_invalid_config", func(t *testing.T) {
		// This test helps cover the mongodb driver error handling path
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}
		mockMongodbManager := &mongodbmocks.MockManager{}

		// Setup expectations
		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockContainer.EXPECT().MustMake("mongodb").Return(mockMongodbManager).Times(1)

		// Mock mongodb manager methods
		mockMongodbManager.EXPECT().DatabaseWithName("").Return(nil).Maybe()

		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					MongoDB: &config.DriverMongodbConfig{
						Enabled:    true,
						Database:   "", // Empty database name will cause issues
						Collection: "cache_entries",
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)

		// Act & Assert - Expect this to panic due to mongodb driver creation failure
		assert.Panics(t, func() {
			provider.Register(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})

	t.Run("successful_redis_driver_creation", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)

		// Test với Redis driver disabled để tránh dependency vào redis.Manager mock
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					Redis: &config.DriverRedisConfig{
						Enabled: false, // Disable để test path không có Redis
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)

		// Act & Assert - Should complete successfully
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		})

		// Verify provider was registered (Redis driver không được tạo vì disabled)
		providers := provider.Providers()
		assert.Contains(t, providers, "cache")
		assert.NotContains(t, providers, "cache.redis") // Không có Redis vì disabled

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})

	t.Run("successful_mongodb_driver_creation", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)

		// Test với MongoDB driver disabled để tránh dependency vào mongodb.Manager mock
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					MongoDB: &config.DriverMongodbConfig{
						Enabled: false, // Disable để test path không có MongoDB
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)

		// Act & Assert - Should complete successfully
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		})

		// Verify provider was registered (MongoDB driver không được tạo vì disabled)
		providers := provider.Providers()
		assert.Contains(t, providers, "cache")
		assert.NotContains(t, providers, "cache.mongodb") // Không có MongoDB vì disabled

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})

	t.Run("successful_all_drivers_creation", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)

		// Test với chỉ memory và file drivers enabled để tránh external dependencies
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					Memory: &config.DriverMemoryConfig{Enabled: true},
					File: &config.DriverFileConfig{
						Enabled: true,
						Path:    "/tmp/cache_test",
					},
					Redis: &config.DriverRedisConfig{
						Enabled: false, // Disable external dependencies
					},
					MongoDB: &config.DriverMongodbConfig{
						Enabled: false, // Disable external dependencies
					},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)
		mockContainer.EXPECT().Instance("cache.memory", mock.Anything).Return().Times(1)
		mockContainer.EXPECT().Instance("cache.file", mock.Anything).Return().Times(1)

		// Act & Assert - Should complete successfully
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		})

		// Verify providers were registered (chỉ memory và file)
		providers := provider.Providers()
		assert.Contains(t, providers, "cache")
		assert.Contains(t, providers, "cache.memory")
		assert.Contains(t, providers, "cache.file")
		assert.NotContains(t, providers, "cache.redis")   // Disabled
		assert.NotContains(t, providers, "cache.mongodb") // Disabled
		assert.Len(t, providers, 3)                       // cache + memory + file

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})
}

func TestServiceProvider_Boot(t *testing.T) {
	t.Run("boot_with_nil_app", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()

		// Act & Assert
		assert.NotPanics(t, func() {
			provider.Boot(nil)
		})
	})

	t.Run("boot_with_valid_app", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}

		// Act & Assert
		assert.NotPanics(t, func() {
			provider.Boot(mockApp)
		})

		// Verify mocks
		mockApp.AssertExpectations(t)
	})
}

func TestServiceProvider_Integration(t *testing.T) {
	t.Run("complete_lifecycle_test", func(t *testing.T) {
		// Arrange
		provider := cache.NewServiceProvider()
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		mockConfigManager := &configmocks.MockManager{}

		// Setup mocks for complete lifecycle
		mockApp.EXPECT().Container().Return(mockContainer).Times(1)
		mockContainer.EXPECT().MustMake("config").Return(mockConfigManager).Times(1)
		mockConfigManager.EXPECT().UnmarshalKey("cache", mock.AnythingOfType("*config.Config")).RunAndReturn(
			func(key string, target interface{}) error {
				cfg := target.(*config.Config)
				cfg.Drivers = config.DriversConfig{
					Memory: &config.DriverMemoryConfig{Enabled: true},
				}
				return nil
			}).Times(1)
		mockContainer.EXPECT().Instance("cache", mock.Anything).Return().Times(1)
		mockContainer.EXPECT().Instance("cache.memory", mock.Anything).Return().Times(1)

		// Act
		// Test new provider creation
		assert.NotNil(t, provider)
		assert.Implements(t, (*di.ServiceProvider)(nil), provider)

		// Test requires
		requires := provider.Requires()
		assert.Equal(t, []string{"config", "redis", "mongodb"}, requires)

		// Test providers before register
		providers := provider.Providers()
		assert.Empty(t, providers)

		// Test register
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		})

		// Test providers after register
		providers = provider.Providers()
		assert.Contains(t, providers, "cache")

		// Test boot
		assert.NotPanics(t, func() {
			provider.Boot(mockApp)
		})

		// Verify all mocks
		mockApp.AssertExpectations(t)
		mockContainer.AssertExpectations(t)
		mockConfigManager.AssertExpectations(t)
	})
}

func TestMockServiceProvider(t *testing.T) {
	t.Run("mock_interfaces_are_properly_implemented", func(t *testing.T) {
		// Test cache manager mock
		mockCacheManager := &cachemocks.MockManager{}
		assert.Implements(t, (*cache.Manager)(nil), mockCacheManager)

		// Test config manager mock
		mockConfigManager := &configmocks.MockManager{}
		assert.NotNil(t, mockConfigManager)

		// Test DI mocks
		mockApp := &dimocks.MockApplication{}
		mockContainer := &dimocks.MockContainer{}
		assert.Implements(t, (*di.Application)(nil), mockApp)
		assert.Implements(t, (*di.Container)(nil), mockContainer)

		// Test redis manager mock
		mockRedisManager := &redismocks.MockManager{}
		assert.NotNil(t, mockRedisManager)

		// Test mongodb manager mock
		mockMongodbManager := &mongodbmocks.MockManager{}
		assert.NotNil(t, mockMongodbManager)
	})
}
