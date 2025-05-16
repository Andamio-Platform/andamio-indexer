package cache

import (
	"container/list"
	"sync"

	input_chainsync "github.com/blinklabs-io/adder/input/chainsync"
)

// CacheItem represents an item in the cache
type CacheItem struct {
	TxHash string
	Event  input_chainsync.TransactionEvent
	Context input_chainsync.TransactionContext
}

// TransactionCache is a simple LRU cache for transactions
type TransactionCache struct {
	limit int
	lock  sync.Mutex
	cache map[string]*list.Element
	ll    *list.List
}

var globalTransactionCache *TransactionCache
var onceTransactionCache sync.Once

// Len returns the current number of items in the cache
func (c *TransactionCache) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.ll.Len()
}

// Limit returns the maximum number of items allowed in the cache
func (c *TransactionCache) Limit() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.limit
}

// NewTransactionCache creates a new TransactionCache with the given limit
func NewTransactionCache(limit int) *TransactionCache {
	return &TransactionCache{
		limit: limit,
		cache: make(map[string]*list.Element),
		ll:    list.New(),
	}
}

// InitTransactionCache initializes the global transaction cache
func InitTransactionCache(limit int) {
	onceTransactionCache.Do(func() {
		globalTransactionCache = NewTransactionCache(limit)
	})
}

// GetTransactionCache returns the singleton instance of the transaction cache
func GetTransactionCache() *TransactionCache {
	// Ensure the cache is initialized before returning
	if globalTransactionCache == nil {
		// This should ideally be initialized during startup, but as a fallback
		// we can initialize with a default limit if accessed before Init
		InitTransactionCache(1000) // Default limit
	}
	return globalTransactionCache
}


// Add adds a transaction event to the cache
func (c *TransactionCache) Add(eventTx input_chainsync.TransactionEvent, eventCtx input_chainsync.TransactionContext) {
	c.lock.Lock()
	defer c.lock.Unlock()

	txHash := string(eventTx.Transaction.Hash().Bytes())

	if element, ok := c.cache[txHash]; ok {
		// Item exists, move to front
		c.ll.MoveToFront(element)
		element.Value.(*CacheItem).Event = eventTx // Update event if needed
		element.Value.(*CacheItem).Context = eventCtx // Update context if needed
	} else {
		// Item does not exist, add to front
		item := &CacheItem{
			TxHash: txHash,
			Event:  eventTx,
			Context: eventCtx,
		}
		element := c.ll.PushFront(item)
		c.cache[txHash] = element

		// Remove oldest if limit is reached
		if c.ll.Len() > c.limit {
			oldest := c.ll.Back()
			if oldest != nil {
				c.ll.Remove(oldest)
				delete(c.cache, oldest.Value.(*CacheItem).TxHash)
			}
		}
	}
}

// Get retrieves a transaction event from the cache
func (c *TransactionCache) Get(txHash string) (*input_chainsync.TransactionEvent, *input_chainsync.TransactionContext, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if element, ok := c.cache[txHash]; ok {
		c.ll.MoveToFront(element)
		item := element.Value.(*CacheItem)
		return &item.Event, &item.Context, true
	}
	return nil, nil, false
}

// GetAll retrieves all transaction events from the cache and clears the cache
func (c *TransactionCache) GetAll() []CacheItem {
	c.lock.Lock()
	defer c.lock.Unlock()

	items := []CacheItem{}
	for element := c.ll.Back(); element != nil; element = element.Prev() {
		items = append(items, *element.Value.(*CacheItem))
	}

	// Clear the cache
	c.cache = make(map[string]*list.Element)
	c.ll = list.New()

	return items
}