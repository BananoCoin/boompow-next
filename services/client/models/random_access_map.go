package models

import (
	"math/rand"
	"sync"

	serializableModels "github.com/bbedward/boompow-ng/libs/models"
)

// RandomAccessMap provides a data struture that can be access randomly
// This helps workers clears backlogs of work more evenly
// If there is 20 items on 3 workers, each worker will access the next unit of work randomly
type RandomAccessMap struct {
	mu     sync.Mutex
	Hashes []serializableModels.ClientWorkRequest
}

func NewRandomAccessMap() *RandomAccessMap {
	return &RandomAccessMap{
		Hashes: []serializableModels.ClientWorkRequest{},
	}
}

// See if element exists
func (r *RandomAccessMap) Exists(hash string) bool {
	for _, v := range r.Hashes {
		if v.Hash == hash {
			return true
		}
	}
	return false
}

// Put value into map - synchronized
func (r *RandomAccessMap) Put(value serializableModels.ClientWorkRequest) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.Exists(value.Hash) {
		r.Hashes = append(r.Hashes, value)
	}
}

// Removes and returns a random value from the map - synchronized
func (r *RandomAccessMap) PopRandom() *serializableModels.ClientWorkRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.Hashes) == 0 {
		return nil
	}
	index := rand.Intn(len(r.Hashes))
	ret := r.Hashes[index]
	r.Hashes = remove(r.Hashes, index)

	return &ret
}

// Gets a value from the map - synchronized
func (r *RandomAccessMap) Get(hash string) *serializableModels.ClientWorkRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Exists(hash) {
		return &r.Hashes[r.IndexOf(hash)]
	}

	return nil
}

// Removes specified hash - synchronized
func (r *RandomAccessMap) Delete(hash string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	index := r.IndexOf(hash)
	if index > -1 {
		r.Hashes = remove(r.Hashes, r.IndexOf(hash))
	}
}

func (r *RandomAccessMap) IndexOf(hash string) int {
	for i, v := range r.Hashes {
		if v.Hash == hash {
			return i
		}
	}
	return -1
}

func remove(s []serializableModels.ClientWorkRequest, i int) []serializableModels.ClientWorkRequest {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
