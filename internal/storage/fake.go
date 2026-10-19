package storage

import (
	"context"
	"sort"
	"sync"
)

// Fake — хранилище в памяти для тестов. Имитирует те же команды Redis.
type Fake struct {
	mu     sync.Mutex
	hits   map[string]int64
	scores map[string]int64
}

func NewFake() *Fake {
	return &Fake{
		hits:   make(map[string]int64),
		scores: make(map[string]int64),
	}
}

func (f *Fake) Hit(_ context.Context, name string) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.hits[name]++
	f.scores[name]++
	return f.hits[name], nil
}

func (f *Fake) Stats(_ context.Context, name string) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.hits[name], nil
}

func (f *Fake) Top(_ context.Context, n int64) ([]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	names := make([]string, 0, len(f.scores))
	for name := range f.scores {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return f.scores[names[i]] > f.scores[names[j]]
	})
	if int64(len(names)) > n {
		names = names[:n]
	}
	return names, nil
}
