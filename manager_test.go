// Package cache_test cung cấp các test cho package cache
package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.fork.vn/cache"
	cache_mocks "go.fork.vn/cache/mocks"
)

// TestManager_New kiểm tra constructor NewManager
func TestManager_New(t *testing.T) {
	t.Run("creates_new_manager_instance", func(t *testing.T) {
		// Act
		manager := cache.NewManager()

		// Assert
		assert.NotNil(t, manager)

		// Test that newly created manager has no default driver
		_, found := manager.Get("any-key")
		assert.False(t, found)
	})
}

// TestManager_Get kiểm tra phương thức Get với các kịch bản khác nhau
func TestManager_Get(t *testing.T) {
	t.Run("returns_value_when_default_driver_is_set_and_key_exists", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Get(context.Background(), "test-key").Return("test-value", true)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, found := manager.Get("test-key")

		// Assert
		assert.True(t, found)
		assert.Equal(t, "test-value", value)
	})

	t.Run("returns_not_found_when_default_driver_is_set_but_key_doesnt_exist", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Get(context.Background(), "nonexistent-key").Return(nil, false)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, found := manager.Get("nonexistent-key")

		// Assert
		assert.False(t, found)
		assert.Nil(t, value)
	})

	t.Run("returns_not_found_when_no_default_driver_is_set", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()

		// Act
		value, found := manager.Get("any-key")

		// Assert
		assert.False(t, found)
		assert.Nil(t, value)
	})

	t.Run("returns_not_found_when_default_driver_doesnt_exist", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()
		manager.SetDefaultDriver("nonexistent")

		// Act
		value, found := manager.Get("any-key")

		// Assert
		assert.False(t, found)
		assert.Nil(t, value)
	})
}

// TestManager_Set kiểm tra phương thức Set với các kịch bản khác nhau
func TestManager_Set(t *testing.T) {
	t.Run("sets_value_successfully_when_default_driver_is_configured", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Set(context.Background(), "test-key", "test-value", 5*time.Minute).Return(nil)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Set("test-key", "test-value", 5*time.Minute)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns_error_when_default_driver_set_operation_fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("driver set error")
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Set(context.Background(), "test-key", "test-value", 5*time.Minute).Return(expectedError)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Set("test-key", "test-value", 5*time.Minute)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("returns_error_when_no_default_driver_is_set", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()

		// Act
		err := manager.Set("test-key", "test-value", 5*time.Minute)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})

	t.Run("returns_error_when_default_driver_doesnt_exist", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()
		manager.SetDefaultDriver("nonexistent")

		// Act
		err := manager.Set("test-key", "test-value", 5*time.Minute)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManager_Has kiểm tra phương thức Has với các kịch bản khác nhau
func TestManager_Has(t *testing.T) {
	t.Run("returns_true_when_key_exists", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Has(context.Background(), "existing-key").Return(true)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		exists := manager.Has("existing-key")

		// Assert
		assert.True(t, exists)
	})

	t.Run("returns_false_when_key_doesnt_exist", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Has(context.Background(), "nonexistent-key").Return(false)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		exists := manager.Has("nonexistent-key")

		// Assert
		assert.False(t, exists)
	})

	t.Run("returns_false_when_no_default_driver_is_set", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()

		// Act
		exists := manager.Has("any-key")

		// Assert
		assert.False(t, exists)
	})

	t.Run("returns_false_when_default_driver_doesnt_exist", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()
		manager.SetDefaultDriver("nonexistent")

		// Act
		exists := manager.Has("any-key")

		// Assert
		assert.False(t, exists)
	})
}

// TestManager_Delete kiểm tra phương thức Delete với các kịch bản khác nhau
func TestManager_Delete(t *testing.T) {
	t.Run("deletes_key_successfully_when_default_driver_is_configured", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Delete(context.Background(), "test-key").Return(nil)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Delete("test-key")

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns_error_when_driver_delete_operation_fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("driver delete error")
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Delete(context.Background(), "test-key").Return(expectedError)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Delete("test-key")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("returns_error_when_no_default_driver_is_set", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()

		// Act
		err := manager.Delete("test-key")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManager_Flush kiểm tra phương thức Flush với các kịch bản khác nhau
func TestManager_Flush(t *testing.T) {
	t.Run("flushes_cache_successfully_when_default_driver_is_configured", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Flush(context.Background()).Return(nil)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Flush()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns_error_when_driver_flush_operation_fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("driver flush error")
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Flush(context.Background()).Return(expectedError)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.Flush()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("returns_error_when_no_default_driver_is_set", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()

		// Act
		err := manager.Flush()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManager_GetMultiple kiểm tra phương thức GetMultiple với các kịch bản khác nhau
func TestManager_GetMultiple(t *testing.T) {
	t.Run("gets_multiple_values_successfully_when_default_driver_is_configured", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2", "key3"}
		expectedFound := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		expectedMissing := []string{"key3"}

		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().GetMultiple(context.Background(), keys).Return(expectedFound, expectedMissing)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		found, missing := manager.GetMultiple(keys)

		// Assert
		assert.Equal(t, expectedFound, found)
		assert.Equal(t, expectedMissing, missing)
	})

	t.Run("returns_empty_map_and_all_keys_as_missing_when_no_default_driver_is_set", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2", "key3"}
		manager := cache.NewManager()

		// Act
		found, missing := manager.GetMultiple(keys)

		// Assert
		assert.Empty(t, found)
		assert.Equal(t, keys, missing)
	})

	t.Run("returns_empty_map_and_all_keys_as_missing_when_default_driver_doesnt_exist", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2", "key3"}
		manager := cache.NewManager()
		manager.SetDefaultDriver("nonexistent")

		// Act
		found, missing := manager.GetMultiple(keys)

		// Assert
		assert.Empty(t, found)
		assert.Equal(t, keys, missing)
	})
}

// TestManager_SetMultiple kiểm tra phương thức SetMultiple với các kịch bản khác nhau
func TestManager_SetMultiple(t *testing.T) {
	t.Run("sets_multiple_values_successfully_when_default_driver_is_configured", func(t *testing.T) {
		// Arrange
		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		ttl := 10 * time.Minute

		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().SetMultiple(context.Background(), values, ttl).Return(nil)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.SetMultiple(values, ttl)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns_error_when_driver_setMultiple_operation_fails", func(t *testing.T) {
		// Arrange
		values := map[string]interface{}{
			"key1": "value1",
		}
		ttl := 10 * time.Minute
		expectedError := errors.New("driver setMultiple error")

		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().SetMultiple(context.Background(), values, ttl).Return(expectedError)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.SetMultiple(values, ttl)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("returns_error_when_no_default_driver_is_set", func(t *testing.T) {
		// Arrange
		values := map[string]interface{}{
			"key1": "value1",
		}
		manager := cache.NewManager()

		// Act
		err := manager.SetMultiple(values, 10*time.Minute)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManager_DeleteMultiple kiểm tra phương thức DeleteMultiple với các kịch bản khác nhau
func TestManager_DeleteMultiple(t *testing.T) {
	t.Run("deletes_multiple_keys_successfully_when_default_driver_is_configured", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2", "key3"}

		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().DeleteMultiple(context.Background(), keys).Return(nil)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.DeleteMultiple(keys)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns_error_when_driver_deleteMultiple_operation_fails", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2"}
		expectedError := errors.New("driver deleteMultiple error")

		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().DeleteMultiple(context.Background(), keys).Return(expectedError)

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		err := manager.DeleteMultiple(keys)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("returns_error_when_no_default_driver_is_set", func(t *testing.T) {
		// Arrange
		keys := []string{"key1", "key2"}
		manager := cache.NewManager()

		// Act
		err := manager.DeleteMultiple(keys)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManager_Remember kiểm tra phương thức Remember với các kịch bản khác nhau
func TestManager_Remember(t *testing.T) {
	t.Run("returns_cached_value_when_key_exists", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Remember(context.Background(), "cached-key", 10*time.Minute, mock.AnythingOfType("func() (interface {}, error)")).Return("cached-value", nil)

		callback := func() (interface{}, error) {
			return "new-value", nil
		}

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, err := manager.Remember("cached-key", 10*time.Minute, callback)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "cached-value", value)
	})

	t.Run("calls_callback_and_caches_result_when_key_doesnt_exist", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Remember(context.Background(), "new-key", 10*time.Minute, mock.AnythingOfType("func() (interface {}, error)")).Return("callback-value", nil)

		callback := func() (interface{}, error) {
			return "callback-value", nil
		}

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, err := manager.Remember("new-key", 10*time.Minute, callback)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "callback-value", value)
	})

	t.Run("returns_callback_error_when_callback_fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("callback error")
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Remember(context.Background(), "new-key", 10*time.Minute, mock.AnythingOfType("func() (interface {}, error)")).Return(nil, expectedError)

		callback := func() (interface{}, error) {
			return nil, expectedError
		}

		manager := cache.NewManager()
		manager.AddDriver("mock", mockDriver)
		manager.SetDefaultDriver("mock")

		// Act
		value, err := manager.Remember("new-key", 10*time.Minute, callback)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, value)
	})

	t.Run("returns_error_when_no_default_driver_is_set", func(t *testing.T) {
		// Arrange
		callback := func() (interface{}, error) {
			return "value", nil
		}
		manager := cache.NewManager()

		// Act
		value, err := manager.Remember("key", 10*time.Minute, callback)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, value)
		assert.Contains(t, err.Error(), "no default cache driver set")
	})
}

// TestManager_AddDriver kiểm tra phương thức AddDriver với các kịch bản khác nhau
func TestManager_AddDriver(t *testing.T) {
	t.Run("adds_driver_successfully", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		manager := cache.NewManager()

		// Act
		manager.AddDriver("test-driver", mockDriver)

		// Assert
		driver, err := manager.Driver("test-driver")
		assert.NoError(t, err)
		assert.Equal(t, mockDriver, driver)
	})

	t.Run("sets_first_added_driver_as_default", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Get(context.Background(), "test-key").Return("test-value", true)

		manager := cache.NewManager()

		// Act
		manager.AddDriver("first-driver", mockDriver)

		// Assert - should be able to use operations without explicitly setting default
		value, found := manager.Get("test-key")
		assert.True(t, found)
		assert.Equal(t, "test-value", value)
	})

	t.Run("replaces_existing_driver_with_same_name", func(t *testing.T) {
		// Arrange
		oldDriver := cache_mocks.NewMockDriver(t)
		newDriver := cache_mocks.NewMockDriver(t)
		manager := cache.NewManager()

		// Act
		manager.AddDriver("same-name", oldDriver)
		manager.AddDriver("same-name", newDriver)

		// Assert
		driver, err := manager.Driver("same-name")
		assert.NoError(t, err)
		assert.Equal(t, newDriver, driver)
	})
}

// TestManager_SetDefaultDriver kiểm tra phương thức SetDefaultDriver
func TestManager_SetDefaultDriver(t *testing.T) {
	t.Run("sets_default_driver_successfully_when_driver_exists", func(t *testing.T) {
		// Arrange
		driver1 := cache_mocks.NewMockDriver(t)
		driver2 := cache_mocks.NewMockDriver(t)
		driver2.EXPECT().Get(context.Background(), "test-key").Return("value-from-driver2", true)

		manager := cache.NewManager()
		manager.AddDriver("driver1", driver1)
		manager.AddDriver("driver2", driver2)

		// Act
		manager.SetDefaultDriver("driver2")

		// Assert - operations should use driver2
		value, found := manager.Get("test-key")
		assert.True(t, found)
		assert.Equal(t, "value-from-driver2", value)
	})

	t.Run("setting_non-existent_driver_as_default_doesnt_cause_immediate_error", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()

		// Act - this should not panic or error immediately
		manager.SetDefaultDriver("non-existent")

		// Assert - but operations should fail
		_, found := manager.Get("test-key")
		assert.False(t, found)
	})
}

// TestManager_Driver kiểm tra phương thức Driver với các kịch bản khác nhau
func TestManager_Driver(t *testing.T) {
	t.Run("returns_driver_when_it_exists", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		manager := cache.NewManager()
		manager.AddDriver("test-driver", mockDriver)

		// Act
		driver, err := manager.Driver("test-driver")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, mockDriver, driver)
	})

	t.Run("returns_error_when_driver_doesnt_exist", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()

		// Act
		driver, err := manager.Driver("non-existent")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, driver)
		assert.Contains(t, err.Error(), "driver 'non-existent' not found")
	})
}

// TestManager_Stats kiểm tra phương thức Stats với các kịch bản khác nhau
func TestManager_Stats(t *testing.T) {
	t.Run("returns_stats_for_all_drivers", func(t *testing.T) {
		// Arrange
		driver1Stats := map[string]interface{}{
			"hits":   100,
			"misses": 20,
		}
		driver2Stats := map[string]interface{}{
			"hits":   50,
			"misses": 10,
		}

		mockDriver1 := cache_mocks.NewMockDriver(t)
		mockDriver1.EXPECT().Stats(context.Background()).Return(driver1Stats)

		mockDriver2 := cache_mocks.NewMockDriver(t)
		mockDriver2.EXPECT().Stats(context.Background()).Return(driver2Stats)

		manager := cache.NewManager()
		manager.AddDriver("driver1", mockDriver1)
		manager.AddDriver("driver2", mockDriver2)

		// Act
		stats := manager.Stats()

		// Assert
		assert.Len(t, stats, 2)
		assert.Equal(t, driver1Stats, stats["driver1"])
		assert.Equal(t, driver2Stats, stats["driver2"])
	})

	t.Run("returns_empty_map_when_no_drivers_are_registered", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()

		// Act
		stats := manager.Stats()

		// Assert
		assert.Empty(t, stats)
	})
}

// TestManager_Close kiểm tra phương thức Close với các kịch bản khác nhau
func TestManager_Close(t *testing.T) {
	t.Run("closes_all_drivers_successfully", func(t *testing.T) {
		// Arrange
		mockDriver1 := cache_mocks.NewMockDriver(t)
		mockDriver1.EXPECT().Close().Return(nil)

		mockDriver2 := cache_mocks.NewMockDriver(t)
		mockDriver2.EXPECT().Close().Return(nil)

		manager := cache.NewManager()
		manager.AddDriver("driver1", mockDriver1)
		manager.AddDriver("driver2", mockDriver2)

		// Act
		err := manager.Close()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("returns_error_when_any_driver_close_fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("close error")

		mockDriver1 := cache_mocks.NewMockDriver(t)
		mockDriver1.EXPECT().Close().Return(nil)

		mockDriver2 := cache_mocks.NewMockDriver(t)
		mockDriver2.EXPECT().Close().Return(expectedError)

		manager := cache.NewManager()
		manager.AddDriver("driver1", mockDriver1)
		manager.AddDriver("driver2", mockDriver2)

		// Act
		err := manager.Close()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "close error")
	})

	t.Run("succeeds_when_no_drivers_are_registered", func(t *testing.T) {
		// Arrange
		manager := cache.NewManager()

		// Act
		err := manager.Close()

		// Assert
		assert.NoError(t, err)
	})
}

// TestManager_Concurrency kiểm tra việc truy cập đồng thời vào manager
func TestManager_Concurrency(t *testing.T) {
	t.Run("concurrent_operations_dont_cause_race_conditions", func(t *testing.T) {
		// Arrange
		mockDriver := cache_mocks.NewMockDriver(t)
		mockDriver.EXPECT().Get(mock.Anything, mock.Anything).Return("value", true).Maybe()
		mockDriver.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
		mockDriver.EXPECT().Has(mock.Anything, mock.Anything).Return(true).Maybe()

		manager := cache.NewManager()
		manager.AddDriver("concurrent-driver", mockDriver)
		manager.SetDefaultDriver("concurrent-driver")

		// Act - perform concurrent operations
		done := make(chan bool, 3)

		go func() {
			for i := 0; i < 100; i++ {
				manager.Get("key")
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 100; i++ {
				_ = manager.Set("key", "value", time.Minute)
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 100; i++ {
				manager.Has("key")
			}
			done <- true
		}()

		// Wait for all goroutines to complete
		for i := 0; i < 3; i++ {
			<-done
		}

		// Assert - if we reach here without panic, the test passes
		assert.True(t, true)
	})
}
