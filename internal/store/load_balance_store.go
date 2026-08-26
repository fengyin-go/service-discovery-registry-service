package store

import "serviceregistry/internal/model"

func (s *MemoryStore) CreateLoadBalanceStrategy(l *model.LoadBalanceStrategy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.loadBalances {
		if exist.ServiceID == l.ServiceID {
			return ErrConflict
		}
	}
	s.loadBalances[l.ID] = l
	return nil
}

func (s *MemoryStore) GetLoadBalanceStrategy(id string) (*model.LoadBalanceStrategy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.loadBalances[id]
	if !ok {
		return nil, ErrNotFound
	}
	return l, nil
}

func (s *MemoryStore) GetLoadBalanceByService(serviceID string) (*model.LoadBalanceStrategy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, l := range s.loadBalances {
		if l.ServiceID == serviceID {
			return l, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListLoadBalanceStrategies() []*model.LoadBalanceStrategy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.LoadBalanceStrategy, 0, len(s.loadBalances))
	for _, l := range s.loadBalances {
		list = append(list, l)
	}
	return list
}

func (s *MemoryStore) UpdateLoadBalanceStrategy(l *model.LoadBalanceStrategy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.loadBalances[l.ID]; !ok {
		return ErrNotFound
	}
	s.loadBalances[l.ID] = l
	return nil
}

func (s *MemoryStore) DeleteLoadBalanceStrategy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.loadBalances[id]; !ok {
		return ErrNotFound
	}
	delete(s.loadBalances, id)
	return nil
}
