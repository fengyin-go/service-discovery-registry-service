package store

import "sync"

type RegistryRefreshLeaseState struct {
	mu     sync.Mutex
	active map[string]uint64
	next   uint64
}

func NewRegistryRefreshLeaseState() *RegistryRefreshLeaseState {
	return &RegistryRefreshLeaseState{active: make(map[string]uint64)}
}

func (s *RegistryRefreshLeaseState) Acquire(key string) (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.active[key]; exists {
		return 0, false
	}
	// 单调递增的 token，用于 Release 校验归属，避免误删新租约。
	s.next++
	token := s.next
	s.active[key] = token
	return token, true
}

// Release 释放 key 对应的占用标记。
func (s *RegistryRefreshLeaseState) Release(key string, token uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.active[key]
	if !ok {
		return false
	}
	// 旧租约已被覆盖（例如重启拿到新 token）时不应误删。
	if token != 0 && current != token {
		return false
	}
	delete(s.active, key)
	return true
}

func (s *RegistryRefreshLeaseState) Active() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.active)
}
