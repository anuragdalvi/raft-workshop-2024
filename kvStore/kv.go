package main

import (
	"encoding/json"
	"os"
	"sync"
)

type Mapping struct {
	data map[string]string
	mu   sync.RWMutex
}

// Create New Instance of Key Value store
func NewKeyValeStore() *Mapping {
	return &Mapping{
		data: make(map[string]string),
	}
}

// Sets key value pair
func (m *Mapping) Set(key, value string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return "ok"
}

// Gets value from a Key
func (m *Mapping) Get(key string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, _ := range m.data {
		if key == k {
			value := m.data[key]
			return value

		}
	}
	return "Not Found"
}

func (m *Mapping) Delete(key string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, _ := range m.data {
		if key == k {
			delete(m.data, key)
			return "ok"
		}
	}
	return "Not Found"
}

func (m *Mapping) Snapshot(filename string) string {
	// Locking using memory synchronization technique
	// Limitation: Will lead to a performance hit.
	// No new clients can interact with the system during the execution of snapshot.
	m.mu.Lock()
	defer m.mu.Unlock()
	jsonData, err := json.Marshal(m.data)
	if err != nil {
		return "Error: Loading JSON"
	}

	snapshot, err := os.Create(filename)
	defer snapshot.Close()

	_, err = snapshot.Write(jsonData)
	if err != nil {
		return "Error: Creating Snapshot"
	}

	return "ok"
}

func (m *Mapping) RetrieveSnapshot(filename string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	snapshotData, err := os.ReadFile(filename)
	if err != nil {
		return "Error: Reading the snapshot"
	}

	err = json.Unmarshal(snapshotData, &m.data)

	if err != nil {
		return "Error: Loading snapshot"
	}

	return "ok"

}
