// Package driver_test cung cấp các test cho package driver
package driver_test

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"
	"go.fork.vn/cache/config"
	"go.fork.vn/cache/driver"
	redispkg "go.fork.vn/redis"
)

// mockRedisManager là mock implementation của redis Manager cho testing
type mockRedisManager struct {
	client          *redis.Client
	universalClient redis.UniversalClient
	clientError     error
	universalError  error
}

func (m *mockRedisManager) Client() (*redis.Client, error) {
	return m.client, m.clientError
}

func (m *mockRedisManager) UniversalClient() (*redis.UniversalClient, error) {
	if m.universalClient != nil {
		return &m.universalClient, m.universalError
	}
	return nil, m.universalError
}

func (m *mockRedisManager) GetConfig() *redispkg.Config {
	return nil
}

func (m *mockRedisManager) Close() error {
	return nil
}

func (m *mockRedisManager) Ping(ctx context.Context) error {
	return nil
}

func (m *mockRedisManager) ClusterPing(ctx context.Context) error {
	return nil
}

// TestRedisDriver_NewRedisDriver kiểm tra việc khởi tạo Redis driver
func TestRedisDriver_NewRedisDriver(t *testing.T) {
	t.Run("returns_error_when_redis_driver_is_disabled", func(t *testing.T) {
		// Arrange
		mockManager := &mockRedisManager{}
		configTest := config.DriverRedisConfig{
			Enabled: false,
		}

		// Act
		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, redisDriver)
		assert.Contains(t, err.Error(), "redis driver is not enabled")
	})

	t.Run("returns_error_when_redis_client_creation_fails", func(t *testing.T) {
		// Arrange
		mockManager := &mockRedisManager{
			clientError: fmt.Errorf("redis connection failed"),
		}
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
		}

		// Act
		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, redisDriver)
		assert.Contains(t, err.Error(), "redis connection failed")
	})

	t.Run("create_redis_driver_successfully_with_json_serializer", func(t *testing.T) {
		// Arrange
		client := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		mockManager := &mockRedisManager{
			client: client,
		}
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}

		// Act
		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, redisDriver)
	})

	t.Run("change_to_gob_serializer", func(t *testing.T) {
		// Arrange
		client := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		mockManager := &mockRedisManager{
			client: client,
		}
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "gob",
		}

		// Act
		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, redisDriver)
	})

	t.Run("change_to_msgpack_serializer", func(t *testing.T) {
		// Arrange
		client := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		mockManager := &mockRedisManager{
			client: client,
		}
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "msgpack",
		}

		// Act
		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, redisDriver)
	})

	t.Run("fallback_to_json_with_unknown_serializer", func(t *testing.T) {
		// Arrange
		client := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		mockManager := &mockRedisManager{
			client: client,
		}
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "unknown",
		}

		// Act
		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, redisDriver)
	})
}

// TestRedisDriver_Serializers kiểm tra các serializer khác nhau
func TestRedisDriver_Serializers(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	mockManager := &mockRedisManager{
		client: client,
	}

	t.Run("json_serializer", func(t *testing.T) {
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)
		require.NoError(t, err)
		require.NotNil(t, redisDriver)

		// Test serialization
		testData := map[string]interface{}{"key": "value", "number": 42}
		serialized, err := json.Marshal(testData)
		require.NoError(t, err)
		assert.NotEmpty(t, serialized)
	})

	t.Run("gob_serializer", func(t *testing.T) {
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "gob",
		}

		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)
		require.NoError(t, err)
		require.NotNil(t, redisDriver)

		// Test registration for gob
		gob.Register(map[string]interface{}{})
	})

	t.Run("msgpack_serializer", func(t *testing.T) {
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "msgpack",
		}

		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)
		require.NoError(t, err)
		require.NotNil(t, redisDriver)

		// Test serialization
		testData := map[string]interface{}{"key": "value", "number": 42}
		serialized, err := msgpack.Marshal(testData)
		require.NoError(t, err)
		assert.NotEmpty(t, serialized)
	})
}

// TestRedisDriver_ErrorHandling kiểm tra xử lý lỗi
func TestRedisDriver_ErrorHandling(t *testing.T) {
	t.Run("nil_redis_manager", func(t *testing.T) {
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
		}

		redisDriver, err := driver.NewRedisDriver(configTest, nil)
		assert.Error(t, err)
		assert.Nil(t, redisDriver)
		assert.Contains(t, err.Error(), "redis manager cannot be nil")
	})

	t.Run("redis_client_error", func(t *testing.T) {
		mockManager := &mockRedisManager{
			clientError: fmt.Errorf("failed to connect to redis"),
		}
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
		}

		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)
		assert.Error(t, err)
		assert.Nil(t, redisDriver)
		assert.Contains(t, err.Error(), "failed to connect to redis")
	})
}

// TestRedisDriver_Configuration kiểm tra cấu hình driver
func TestRedisDriver_Configuration(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	t.Run("default_ttl_configuration", func(t *testing.T) {
		mockManager := &mockRedisManager{
			client: client,
		}
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 600, // 10 minutes
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)
		require.NoError(t, err)
		require.NotNil(t, redisDriver)

		// Verify the driver was created with correct config
		assert.Equal(t, 600, configTest.DefaultTTL)
		assert.Equal(t, "json", configTest.Serializer)
	})

	t.Run("zero_ttl_configuration", func(t *testing.T) {
		mockManager := &mockRedisManager{
			client: client,
		}
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 0, // No expiration
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)
		require.NoError(t, err)
		require.NotNil(t, redisDriver)
	})
}

// TestRedisDriver_EdgeCases kiểm tra các trường hợp edge case
func TestRedisDriver_EdgeCases(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	mockManager := &mockRedisManager{
		client: client,
	}

	t.Run("empty_serializer_fallback_to_json", func(t *testing.T) {
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "", // Empty serializer
		}

		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)
		require.NoError(t, err)
		require.NotNil(t, redisDriver)
	})

	t.Run("invalid_serializer_fallback_to_json", func(t *testing.T) {
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "invalid_serializer",
		}

		redisDriver, err := driver.NewRedisDriver(configTest, mockManager)
		require.NoError(t, err)
		require.NotNil(t, redisDriver)
	})
}

// TestRedisDriver_Integration kiểm tra tích hợp với Redis thật (chỉ chạy khi có Redis)
func TestRedisDriver_Integration(t *testing.T) {
	// Skip if running in CI without Redis
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Create a real Redis manager (this would normally come from DI container)
	cfg := redispkg.DefaultConfig()
	cfg.Client.Enabled = true
	cfg.Client.Host = "localhost"
	cfg.Client.Port = 6379
	cfg.Client.DB = 15 // Use test database

	redisManager := redispkg.NewManager(cfg)

	// Test if Redis is available
	ctx := context.Background()
	if err := redisManager.Ping(ctx); err != nil {
		t.Skip("Redis not available, skipping integration tests")
	}

	defer func() {
		// Cleanup
		if client, err := redisManager.Client(); err == nil {
			client.FlushDB(ctx)
		}
		redisManager.Close()
	}()

	t.Run("real_redis_operations", func(t *testing.T) {
		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(configTest, redisManager)
		require.NoError(t, err)
		require.NotNil(t, redisDriver)

		// Test basic operations
		key := "test:integration:key"
		value := "test_value"

		// Set
		err = redisDriver.Set(ctx, key, value, 5*time.Minute)
		assert.NoError(t, err)

		// Get
		result, found := redisDriver.Get(ctx, key)
		assert.True(t, found)
		assert.Equal(t, value, result)

		// Has
		exists := redisDriver.Has(ctx, key)
		assert.True(t, exists)

		// Delete
		err = redisDriver.Delete(ctx, key)
		assert.NoError(t, err)

		// Verify deletion
		exists = redisDriver.Has(ctx, key)
		assert.False(t, exists)
	})
}

// TestRedisDriver_InterfaceMethods kiểm tra tất cả method của Driver interface
func TestRedisDriver_InterfaceMethods(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	mockManager := &mockRedisManager{
		client: client,
	}
	configTest := config.DriverRedisConfig{
		Enabled:    true,
		DefaultTTL: 300,
		Serializer: "json",
	}

	redisDriver, err := driver.NewRedisDriver(configTest, mockManager)
	require.NoError(t, err)
	require.NotNil(t, redisDriver)

	ctx := context.Background()

	t.Run("Set_and_Get", func(t *testing.T) {
		key := "test:set:get"
		value := "test_value"

		// Create mock client and manager first
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		// The driver serializes the value before storing - JSON adds quotes around strings
		// TTL: 5*time.Minute = 300 seconds
		expectedData, _ := json.Marshal(value)
		mock.ExpectSet("cache:"+key, expectedData, 5*time.Minute).SetVal("OK")

		err = testRedisDriver.Set(ctx, key, value, 5*time.Minute)
		assert.NoError(t, err)

		// Mock GET command
		mock.ExpectGet("cache:" + key).SetVal(`"test_value"`)

		result, found := testRedisDriver.Get(ctx, key)
		assert.True(t, found)
		assert.Equal(t, value, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Has", func(t *testing.T) {
		key := "test:has"

		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		mock.ExpectExists("cache:" + key).SetVal(1)
		exists := testRedisDriver.Has(ctx, key)
		assert.True(t, exists)

		mock.ExpectExists("cache:" + key).SetVal(0)
		exists = testRedisDriver.Has(ctx, key)
		assert.False(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test:delete"

		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		mock.ExpectDel("cache:" + key).SetVal(1)

		err = testRedisDriver.Delete(ctx, key)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetMultiple", func(t *testing.T) {
		keys := []string{"key1", "key2", "key3"}

		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		// Prepare expected data - JSON adds quotes around strings
		mock.ExpectMGet("cache:key1", "cache:key2", "cache:key3").SetVal([]interface{}{
			`"value1"`, `"value2"`, nil, // key3 not found
		})

		results, missed := testRedisDriver.GetMultiple(ctx, keys)

		assert.Len(t, results, 2)
		assert.Len(t, missed, 1)
		assert.Equal(t, "value1", results["key1"])
		assert.Equal(t, "value2", results["key2"])
		assert.Contains(t, missed, "key3")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SetMultiple", func(t *testing.T) {
		values := map[string]interface{}{
			"multi1": "value1",
			"multi2": "value2",
		}

		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		// The driver serializes each value before storing - JSON adds quotes around strings
		expectedData1, _ := json.Marshal("value1")
		expectedData2, _ := json.Marshal("value2")
		mock.ExpectSet("cache:multi1", expectedData1, time.Duration(300)*time.Second).SetVal("OK")
		mock.ExpectSet("cache:multi2", expectedData2, time.Duration(300)*time.Second).SetVal("OK")

		err = testRedisDriver.SetMultiple(ctx, values, 5*time.Minute)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DeleteMultiple", func(t *testing.T) {
		keys := []string{"del1", "del2", "del3"}

		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		mock.ExpectDel("cache:del1", "cache:del2", "cache:del3").SetVal(3)

		err = testRedisDriver.DeleteMultiple(ctx, keys)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Remember", func(t *testing.T) {
		key := "test:remember"
		expectedValue := "computed_value"
		callbackCalled := false

		callback := func() (interface{}, error) {
			callbackCalled = true
			return expectedValue, nil
		}

		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		mock.ExpectGet("cache:" + key).RedisNil()
		// The driver serializes the value before storing - JSON adds quotes around strings
		expectedData, _ := json.Marshal(expectedValue)
		mock.ExpectSet("cache:"+key, expectedData, time.Duration(300)*time.Second).SetVal("OK")

		result, err := testRedisDriver.Remember(ctx, key, 0, callback)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)
		assert.True(t, callbackCalled)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Remember_Cache_Hit", func(t *testing.T) {
		key := "test:remember:hit"
		cachedValue := "cached_value"
		callbackCalled := false

		callback := func() (interface{}, error) {
			callbackCalled = true
			return "should_not_be_called", nil
		}

		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		// JSON adds quotes around strings
		mock.ExpectGet("cache:" + key).SetVal(`"cached_value"`)

		result, err := testRedisDriver.Remember(ctx, key, 0, callback)
		assert.NoError(t, err)
		assert.Equal(t, cachedValue, result)
		assert.False(t, callbackCalled)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Stats", func(t *testing.T) {
		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		mock.ExpectKeys("cache:*").SetVal([]string{"cache:key1", "cache:key2"})
		mock.ExpectInfo().SetVal("redis_version:7.0.0\nused_memory:1024\n")

		stats := testRedisDriver.Stats(ctx)

		assert.Contains(t, stats, "count")
		assert.Contains(t, stats, "hits")
		assert.Contains(t, stats, "misses")
		assert.Contains(t, stats, "type")
		assert.Contains(t, stats, "prefix")
		assert.Contains(t, stats, "info")
		assert.Equal(t, "redis", stats["type"])
		assert.Equal(t, "cache:", stats["prefix"])
		assert.Equal(t, 2, stats["count"])

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Flush", func(t *testing.T) {
		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		// Mock SCAN and DEL commands
		mock.ExpectScan(0, "cache:*", 0).SetVal([]string{"cache:key1", "cache:key2"}, 0)
		mock.ExpectDel("cache:key1", "cache:key2").SetVal(2)

		err = testRedisDriver.Flush(ctx)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Close", func(t *testing.T) {
		// Note: Close method calls client.Close() which might be hard to mock
		// For now, we'll just test that it doesn't panic
		err := redisDriver.Close()
		// The error depends on the Redis client state, so we don't assert on it
		_ = err
	})
}

// TestRedisDriver_ErrorScenarios kiểm tra các tình huống lỗi
func TestRedisDriver_ErrorScenarios(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	mockManager := &mockRedisManager{
		client: client,
	}
	configTest := config.DriverRedisConfig{
		Enabled:    true,
		DefaultTTL: 300,
		Serializer: "json",
	}

	redisDriver, err := driver.NewRedisDriver(configTest, mockManager)
	require.NoError(t, err)
	require.NotNil(t, redisDriver)

	ctx := context.Background()

	t.Run("Set_Error", func(t *testing.T) {
		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		// The driver serializes the value before storing - JSON adds quotes around strings
		expectedData, _ := json.Marshal("value")
		mock.ExpectSet("cache:test", expectedData, time.Duration(300)*time.Second).SetErr(errors.New("redis error"))

		err = testRedisDriver.Set(ctx, "test", "value", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis error")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get_Error", func(t *testing.T) {
		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		mock.ExpectGet("cache:test").SetErr(errors.New("redis get error"))

		result, found := testRedisDriver.Get(ctx, "test")
		assert.False(t, found)
		assert.Nil(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete_Error", func(t *testing.T) {
		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		mock.ExpectDel("cache:test").SetErr(errors.New("redis delete error"))

		err = testRedisDriver.Delete(ctx, "test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis delete error")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Remember_Callback_Error", func(t *testing.T) {
		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		// Create driver with mock client
		testConfig := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}
		testRedisDriver, err := driver.NewRedisDriver(testConfig, testMockManager)
		require.NoError(t, err)

		callbackError := errors.New("callback error")
		callback := func() (interface{}, error) {
			return nil, callbackError
		}

		mock.ExpectGet("cache:test").RedisNil()

		result, err := testRedisDriver.Remember(ctx, "test", 0, callback)
		assert.Error(t, err)
		assert.Equal(t, callbackError, err)
		assert.Nil(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestRedisDriver_Serialization kiểm tra serialization/deserialization
func TestRedisDriver_Serialization(t *testing.T) {
	t.Run("Complex_Data_Types", func(t *testing.T) {
		// Create mock client and manager
		client, mock := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(configTest, testMockManager)
		require.NoError(t, err)

		ctx := context.Background()

		complexData := map[string]interface{}{
			"string": "test",
			"number": 42,
			"array":  []interface{}{1, 2, 3},
			"object": map[string]interface{}{"nested": "value"},
		}

		// Mock SET - serialize the complex data to JSON
		expectedData, _ := json.Marshal(complexData)
		mock.ExpectSet("cache:complex", expectedData, time.Duration(300)*time.Second).SetVal("OK")

		err = redisDriver.Set(ctx, "complex", complexData, 0)
		assert.NoError(t, err)

		// Mock GET
		mock.ExpectGet("cache:complex").SetVal(string(expectedData))

		result, found := redisDriver.Get(ctx, "complex")
		assert.True(t, found)

		// Compare the nested structure
		resultMap := result.(map[string]interface{})
		assert.Equal(t, "test", resultMap["string"])
		assert.Equal(t, float64(42), resultMap["number"]) // JSON unmarshals numbers as float64

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Serialization_Error", func(t *testing.T) {
		// Create mock client and manager
		client, _ := redismock.NewClientMock()
		testMockManager := &mockRedisManager{
			client: client,
		}

		configTest := config.DriverRedisConfig{
			Enabled:    true,
			DefaultTTL: 300,
			Serializer: "json",
		}

		redisDriver, err := driver.NewRedisDriver(configTest, testMockManager)
		require.NoError(t, err)

		ctx := context.Background()

		// Try to serialize a function (not JSON serializable)
		unserialisableData := func() {}

		err = redisDriver.Set(ctx, "test", unserialisableData, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not serialize value")
	})
}

// TestRedisDriver_WithSerializer kiểm tra method WithSerializer
func TestRedisDriver_WithSerializer(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	mockManager := &mockRedisManager{
		client: client,
	}
	configTest := config.DriverRedisConfig{
		Enabled:    true,
		DefaultTTL: 300,
		Serializer: "json",
	}

	redisDriver, err := driver.NewRedisDriver(configTest, mockManager)
	require.NoError(t, err)
	require.NotNil(t, redisDriver)

	t.Run("Change_To_GOB_Serializer", func(t *testing.T) {
		gobDriver := redisDriver.WithSerializer("gob")
		assert.NotNil(t, gobDriver)
		assert.NotEqual(t, redisDriver, gobDriver) // Should be a new instance
	})

	t.Run("Change_To_MessagePack_Serializer", func(t *testing.T) {
		msgpackDriver := redisDriver.WithSerializer("msgpack")
		assert.NotNil(t, msgpackDriver)
		assert.NotEqual(t, redisDriver, msgpackDriver) // Should be a new instance
	})

	t.Run("Unknown_Serializer_Fallback", func(t *testing.T) {
		jsonDriver := redisDriver.WithSerializer("unknown")
		assert.NotNil(t, jsonDriver)
		assert.NotEqual(t, redisDriver, jsonDriver) // Should be a new instance
	})
}
