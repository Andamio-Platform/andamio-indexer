package cache

import (
	"sync"

	"github.com/Andamio-Platform/andamio-indexer/config"
	"github.com/Andamio-Platform/andamio-indexer/database"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

// RelevantDataCache holds the addresses and policies to filter transactions
type RelevantDataCache struct {
	Addresses []string
	Policies  []string
	mu        sync.RWMutex
}

var globalCache *RelevantDataCache
var once sync.Once

// GetRelevantDataCache returns the singleton instance of the cache
func GetRelevantDataCache() *RelevantDataCache {
	once.Do(func() {
		globalCache = &RelevantDataCache{}
		globalCache.LoadCache() // Load data on initialization
	})
	return globalCache
}

// LoadCache loads relevant addresses and policies into the cache
func (c *RelevantDataCache) LoadCache() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Load config
	cfg := config.GetGlobalConfig()

	// Load addresses from database
	globalDB := database.GetGlobalDB()
	dbAddresses, err := globalDB.GetAllAddresses()
	if err != nil {
		fiberLogger.Error("failed to retrieve addresses from database for cache: %v", err)
		// Continue with config addresses even if database read fails
	}

	// Combine database and config addresses
	c.Addresses = append(dbAddresses, cfg.Andamio.GetAllAndamioAddresses()...)

	// Load policies from config
	c.Policies = cfg.Andamio.GetAllAndamioPolicies()

	fiberLogger.Info("Relevant data cache loaded successfully")
}

// GetAddresses returns the cached addresses
func (c *RelevantDataCache) GetAddresses() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Addresses
}

// GetPolicies returns the cached policies
func (c *RelevantDataCache) GetPolicies() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Policies
}