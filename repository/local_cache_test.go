package repository

import (
	"sync"
	"testing"

	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/stretchr/testify/assert"
)

func TestDronePlanCache_SetAndGet(t *testing.T) {
	cache := NewDronePlanCache()

	estateId := "1"
	distance := 10
	coord := model.Coordinate{X: 5, Y: 8}

	cache.Set(estateId, distance, coord)

	retrievedCoordinate, exists := cache.Get(estateId, distance)
	assert.True(t, exists)
	assert.Equal(t, coord, retrievedCoordinate)
}

func TestDronePlanCache_GetNonExistentKey(t *testing.T) {
	cache := NewDronePlanCache()

	estateId := "1"
	distance := 999

	_, exists := cache.Get(estateId, distance)
	assert.False(t, exists)
}

func TestDronePlanCache_Clear(t *testing.T) {
	cache := NewDronePlanCache()

	estateId := "1"
	distance := 10
	coordinate := model.Coordinate{X: 5, Y: 8}

	// Set cache value
	cache.Set(estateId, distance, coordinate)

	// Verify that the cache has the value
	retrievedCoord, exists := cache.Get(estateId, distance)
	assert.True(t, exists)
	assert.Equal(t, coordinate, retrievedCoord)

	// Clear the cache for this estateId
	cache.ClearByEstateId(estateId)

	// Verify that the cache for estateId is cleared
	_, exists = cache.Get(estateId, distance)
	assert.False(t, exists)
}

func TestDronePlanCache_ConcurrentAccess(t *testing.T) {
	cache := NewDronePlanCache()
	var wg sync.WaitGroup

	estateId := "1"
	distance := 10
	coordinate := model.Coordinate{X: 5, Y: 8}

	// Run multiple goroutines to test concurrent access to cache
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Set(estateId, distance, coordinate)
		}()
	}

	wg.Wait()

	// After all goroutines finish, ensure the cache is set correctly
	retrievedCoordinate, exists := cache.Get(estateId, distance)
	assert.True(t, exists)
	assert.Equal(t, coordinate, retrievedCoordinate)
}
