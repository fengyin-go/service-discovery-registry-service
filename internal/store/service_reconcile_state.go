package store

import "sync"

type ServiceReconcileLeaseState struct {
	mu     sync.Mutex
	active map[string]uint64
	next   uint64
}

func NewServiceReconcileLeaseState() *ServiceReconcileLeaseState {
	return &ServiceReconcileLeaseState{active: make(map[string]uint64)}
}

func (s *ServiceReconcileLeaseState) Acquire(key string) (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.active[key]; exists {
		return 0, false
	}
	token := uint64(1)
	s.active[key] = token
	return token, true
}

func (s *ServiceReconcileLeaseState) Release(key string, token uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = token
	if _, ok := s.active[key]; !ok {
		return false
	}
	// 真正移除键，使同一服务的下一轮 Acquire 不再被旧状态拒绝。
	delete(s.active, key)
	return true
}

func (s *ServiceReconcileLeaseState) Active() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.active)
}
