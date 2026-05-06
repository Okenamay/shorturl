package memstorage

import "sync"

type URLMap struct {
	mu sync.RWMutex
	m  map[string]string
}

var Store = NewURLMap()

func NewURLMap() *URLMap {
	return &URLMap{
		m: make(map[string]string),
	}
}

func (u *URLMap) Set(shortID, fullURL string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.m[shortID] = fullURL
}

func (u *URLMap) Get(shortID string) (string, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	val, ok := u.m[shortID]

	return val, ok
}

func (u *URLMap) GetAll() map[string]string {
	u.mu.RLock()
	defer u.mu.RUnlock()

	newMap := make(map[string]string, len(u.m))
	for k, v := range u.m {
		newMap[k] = v
	}

	return newMap
}
