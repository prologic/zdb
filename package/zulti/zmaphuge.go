package zulti

import (
	"fmt"
	"hash/fnv"
	"sync"
)

var hashLen = uint32(5000)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s[0:10]))
	return h.Sum32()
}

// ConcurrentMapHuge :
type ConcurrentMapHuge struct {
	sync.RWMutex
	items map[uint32]*ConcurrentMap
}

// NewConcurrentMapHuge :
func NewConcurrentMapHuge() *ConcurrentMapHuge {
	// Make items
	times := map[uint32]*ConcurrentMap{}
	for i := uint32(0); i < hashLen; i++ {
		times[i] = NewConcurrentMap()
	}

	cm := &ConcurrentMapHuge{
		items: times,
	}

	return cm
}

// Set : Sets a key in a concurrent map
func (cm *ConcurrentMapHuge) Set(key string, value interface{}) {
	cm.Lock()
	defer cm.Unlock()
	cm.items[hash(key)%hashLen].Set(key, value)
}

// Get : Gets a key from a concurrent map
func (cm *ConcurrentMapHuge) Get(key string) (interface{}, bool) {
	cm.Lock()
	defer cm.Unlock()

	value, ok := cm.items[hash(key)%hashLen].Get(key)
	return value, ok
}

// Len : Get len
func (cm *ConcurrentMapHuge) Len() int {
	cm.Lock()
	defer cm.Unlock()

	total := 0
	for _, element := range cm.items {
		templen := element.Len()
		total += templen
	}

	return total
}

// Save : Save map
func (cm *ConcurrentMapHuge) Save(name string) {
	cm.Lock()
	defer cm.Unlock()

	for key, element := range cm.items {
		element.Save(name + fmt.Sprintf("%v", key))
	}
}

// Load :
func (cm *ConcurrentMapHuge) Load(name string) {
	cm.Lock()
	defer cm.Unlock()

	for key, element := range cm.items {
		element.Load(name + fmt.Sprintf("%v", key))
	}
}
