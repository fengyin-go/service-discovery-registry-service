package store

import "sync"

type DependencyRebuildLeaseState struct {
	mu     sync.Mutex
	active map[string]uint64
	next   uint64
}

func NewDependencyRebuildLeaseState() *DependencyRebuildLeaseState {
	return &DependencyRebuildLeaseState{active: make(map[string]uint64)}
}

func (s *DependencyRebuildLeaseState) Acquire(key string) (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.active[key]; exists {
		return 0, false
	}
	s.next++
	token := s.next
	s.active[key] = token
	return token, true
}

func (s *DependencyRebuildLeaseState) Release(key string, token uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur, ok := s.active[key]
	if !ok {
		return false
	}
	// 仅在 token 匹配时删除租约，避免误清后续新一轮重建持有的租约。
	if cur != token {
		return false
	}
	delete(s.active, key)
	return true
}

func (s *DependencyRebuildLeaseState) Active() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.active)
}
