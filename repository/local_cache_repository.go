package repository

import (
	"sync"

	"github.com/ahmadirfaan/plantation-engine/model"
)

type DronePlanCache struct {
	data map[string]map[int]model.Coordinate
	mu   sync.RWMutex
}

func NewDronePlanCache() *DronePlanCache {
	return &DronePlanCache{
		data: make(map[string]map[int]model.Coordinate),
	}
}

func (c *DronePlanCache) Set(estateId string, distance int, coordinate model.Coordinate) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.data[estateId]; !exists {
		c.data[estateId] = make(map[int]model.Coordinate)
	}
	c.data[estateId][distance] = coordinate
}

func (c *DronePlanCache) Get(estateId string, distance int) (model.Coordinate, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if distMap, ok := c.data[estateId]; ok {
		coordinate, found := distMap[distance]
		return coordinate, found
	}
	return model.Coordinate{}, false
}

func (c *DronePlanCache) ClearByEstateId(estateId string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, estateId)
}
