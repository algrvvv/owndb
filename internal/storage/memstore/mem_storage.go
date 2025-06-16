package memstore

import (
	"fmt"
	"sync"

	"github.com/algrvvv/owndb/internal/storage"
)

type MemStorage struct {
	mu   sync.RWMutex
	data map[string]any
}

func NewMemStorage(initData map[string]any) storage.Storage {
	if initData == nil {
		return &MemStorage{
			data: make(map[string]any),
		}
	}

	return &MemStorage{
		data: initData,
	}
}

// Get implements storage.Storage.
func (m *MemStorage) Get(key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, ok := m.data[key]
	return data, ok
}

// GetAll implements storage.Storage.
func (m *MemStorage) GetAll() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data
}

// Set implements storage.Storage.
func (m *MemStorage) Set(key string, data any) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// TODO: делать запись в owl logs.
	println("!!!TODO: SAVE IN OWL LOGS!!!")
	m.data[key] = data
	return nil
}

// Remove implements storage.Storage.
func (m *MemStorage) Remove(key string) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	return nil
}

// Save implements storage.Storage.
func (m *MemStorage) Save() (err error) {
	fmt.Println(m.data)
	panic("unimplemented")
}

// Keys implements storage.Storage.
func (m *MemStorage) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var keys []string
	for k := range m.data {
		keys = append(keys, k)
	}

	return keys
}
