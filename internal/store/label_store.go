package store

import "serviceregistry/internal/model"

func (s *MemoryStore) CreateServiceLabel(l *model.ServiceLabel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.labels {
		if exist.ServiceID == l.ServiceID && exist.Key == l.Key {
			return ErrConflict
		}
	}
	s.labels[l.ID] = l
	return nil
}

func (s *MemoryStore) GetServiceLabel(id string) (*model.ServiceLabel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.labels[id]
	if !ok {
		return nil, ErrNotFound
	}
	return l, nil
}

func (s *MemoryStore) ListServiceLabels() []*model.ServiceLabel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ServiceLabel, 0, len(s.labels))
	for _, l := range s.labels {
		list = append(list, l)
	}
	return list
}

func (s *MemoryStore) UpdateServiceLabel(l *model.ServiceLabel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.labels[l.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.labels {
		if exist.ID != l.ID && exist.ServiceID == l.ServiceID && exist.Key == l.Key {
			return ErrConflict
		}
	}
	s.labels[l.ID] = l
	return nil
}

func (s *MemoryStore) DeleteServiceLabel(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.labels[id]; !ok {
		return ErrNotFound
	}
	delete(s.labels, id)
	return nil
}
