package store

import "sync"

type ProbeLeaseLeaseState struct {
	mu     sync.Mutex
	active map[string]uint64
	next   uint64
}

func NewProbeLeaseLeaseState() *ProbeLeaseLeaseState {
	return &ProbeLeaseLeaseState{active: make(map[string]uint64)}
}

func (s *ProbeLeaseLeaseState) Acquire(key string) (uint64, bool) {
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

// Release 释放租约，从 active 中彻底移除 key，使同一实例可被再次获取。
func (s *ProbeLeaseLeaseState) Release(key string, token uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur, ok := s.active[key]
	if !ok {
		return false
	}
	// token 不匹配时拒绝释放，避免误归还他人持有的租约。
	if token != 0 && cur != token {
		return false
	}
	delete(s.active, key)
	return true
}

func (s *ProbeLeaseLeaseState) Active() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.active)
}
