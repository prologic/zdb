package zulti

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

//////////////////////
// S L I C E  M A P //
//////////////////////

// ConcurrentSliceMM : Slice type that can be safely shared between goroutines
type ConcurrentSliceMM struct {
	sync.RWMutex
	items map[string][]map[string]interface{}
}

// Append : Appends an item to the concurrent slice
func (cs *ConcurrentSliceMM) Append(itemtype string, item map[string]interface{}) {
	cs.Lock()
	defer cs.Unlock()

	if cs.items == nil {
		cs.items = map[string][]map[string]interface{}{}
	}

	if cs.items[itemtype] == nil {
		cs.items[itemtype] = []map[string]interface{}{}
	}
	cs.items[itemtype] = append(cs.items[itemtype], item)
}

// Len :
func (cs *ConcurrentSliceMM) Len(itemtype string) int {
	cs.Lock()
	defer cs.Unlock()
	value := len(cs.items[itemtype])
	return value
}

// Clear :
func (cs *ConcurrentSliceMM) Clear(itemtype string) []map[string]interface{} {
	cs.Lock()
	defer cs.Unlock()
	items := cs.items[itemtype]

	if cs.items == nil {
		cs.items = map[string][]map[string]interface{}{}
	}
	cs.items[itemtype] = []map[string]interface{}{}
	return items
}

// Get :
func (cs *ConcurrentSliceMM) Get(itemtype string, index int) map[string]interface{} {
	if len(cs.items[itemtype]) >= index {
		return cs.items[itemtype][index]
	}
	return map[string]interface{}{}
}

//////////////////////
// S L I C E  M A P //
//////////////////////

// ConcurrentSliceM : Slice type that can be safely shared between goroutines
type ConcurrentSliceM struct {
	sync.RWMutex
	items []map[string]interface{}
}

// Append : Appends an item to the concurrent slice
func (cs *ConcurrentSliceM) Append(item map[string]interface{}) {
	cs.Lock()
	defer cs.Unlock()
	cs.items = append(cs.items, item)
}

// Len :
func (cs *ConcurrentSliceM) Len() int {
	cs.Lock()
	defer cs.Unlock()
	value := len(cs.items)
	return value
}

// Clear :
func (cs *ConcurrentSliceM) Clear() []map[string]interface{} {
	cs.Lock()
	defer cs.Unlock()
	items := cs.items
	cs.items = []map[string]interface{}{}
	return items
}

// Get :
func (cs *ConcurrentSliceM) Get(index int) map[string]interface{} {
	if len(cs.items) >= index {
		return cs.items[index]
	}
	return map[string]interface{}{}
}

// ConcurrentSliceMItem : Concurrent slice item
type ConcurrentSliceMItem struct {
	Index int
	Value map[string]interface{}
}

// Iter :
func (cs *ConcurrentSliceM) Iter() <-chan ConcurrentSliceMItem {
	c := make(chan ConcurrentSliceMItem)
	f := func() {
		cs.Lock()
		defer cs.Lock()
		for index, value := range cs.items {
			c <- ConcurrentSliceMItem{index, value}
		}
		close(c)
	}
	go f()
	return c
}

///////////////
// S L I C E //
///////////////

// ConcurrentSlice : Slice type that can be safely shared between goroutines
type ConcurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

// ConcurrentSliceItem : Concurrent slice item
type ConcurrentSliceItem struct {
	Index int
	Value interface{}
}

// Append : Appends an item to the concurrent slice
func (cs *ConcurrentSlice) Append(item interface{}) {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, item)
}

// Len :
func (cs *ConcurrentSlice) Len() int {
	cs.Lock()
	defer cs.Unlock()

	value := len(cs.items)
	return value
}

// Iter :
// Iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builin range keyword
func (cs *ConcurrentSlice) Iter() <-chan ConcurrentSliceItem {
	c := make(chan ConcurrentSliceItem)

	f := func() {
		cs.Lock()
		defer cs.Lock()
		for index, value := range cs.items {
			c <- ConcurrentSliceItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}

// ConcurrentMap :
type ConcurrentMap struct {
	sync.RWMutex
	items map[string]interface{}
}

// ConcurrentMapItem : Concurrent map item
type ConcurrentMapItem struct {
	Key   string
	Value interface{}
}

// NewConcurrentMap :
func NewConcurrentMap() *ConcurrentMap {
	cm := &ConcurrentMap{
		items: make(map[string]interface{}),
	}

	return cm
}

// Set : Sets a key in a concurrent map
func (cm *ConcurrentMap) Set(key string, value interface{}) {
	cm.Lock()
	defer cm.Unlock()

	cm.items[key] = value
}

// Get : Gets a key from a concurrent map
func (cm *ConcurrentMap) Get(key string) (interface{}, bool) {
	cm.Lock()
	defer cm.Unlock()

	value, ok := cm.items[key]

	return value, ok
}

// Len : Get len
func (cm *ConcurrentMap) Len() int {
	cm.Lock()
	defer cm.Unlock()
	return len(cm.items)
}

// Save : Save map
func (cm *ConcurrentMap) Save(name string) {
	cm.Lock()
	defer cm.Unlock()

	// Get path
	pathFull, _ := os.Executable()
	path := filepath.Dir(pathFull)

	// Create a file
	os.Remove(path + "/gob/" + name)
	dataFile, _ := os.Create(path + "/gob/" + name)
	defer dataFile.Close()

	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(cm.items)
	for k := range cm.items {
		cm.items[k] = nil
	}
}

// Load :
func (cm *ConcurrentMap) Load(name string) {
	cm.Lock()
	defer cm.Unlock()

	// Get path
	pathFull, _ := os.Executable()
	path := filepath.Dir(pathFull)
	dataFile, err := os.Open(path + "/gob/" + name)
	if err == nil {
		dataDecoder := gob.NewDecoder(dataFile)
		dataDecoder.Decode(&cm.items)
	}
	dataFile.Close()
}

// CopyStruct :
func (cm *ConcurrentMap) CopyStruct(oldmap map[string]interface{}) error {
	cm.Lock()
	defer cm.Unlock()

	for k := range cm.items {
		delete(cm.items, k)
	}

	cm.items = make(map[string]interface{})
	for k, v := range oldmap {
		cm.items[k] = v
	}

	return nil
}

// GetStruct :
func (cm *ConcurrentMap) GetStruct(key string, v interface{}) (err error) {
	cm.Lock()
	defer cm.Unlock()

	if cm.items[key] == nil {
		v = nil
	} else {
		tempByte := []byte(cm.items[key].(string))
		err = json.Unmarshal(tempByte, v)
	}

	return err
}

// GetDump :
func (cm *ConcurrentMap) GetDump() (result string) {
	cm.Lock()
	defer cm.Unlock()

	for k, v := range cm.items {
		result += fmt.Sprintf("%v : %v, \n", k, v)
	}

	return result
}

// Iter :
// Iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (cm *ConcurrentMap) Iter() <-chan ConcurrentMapItem {
	c := make(chan ConcurrentMapItem)

	f := func() {
		cm.Lock()
		defer cm.Unlock()

		for k, v := range cm.items {
			c <- ConcurrentMapItem{k, v}
		}
		close(c)
	}
	go f()

	return c
}
