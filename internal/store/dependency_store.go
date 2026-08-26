package store

import "serviceregistry/internal/model"

func (s *MemoryStore) CreateDependency(d *model.ServiceDependency) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.dependencies {
		if exist.ServiceID == d.ServiceID && exist.DependsOnID == d.DependsOnID {
			return ErrConflict
		}
	}
	s.dependencies[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDependency(id string) (*model.ServiceDependency, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.dependencies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListDependencies() []*model.ServiceDependency {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ServiceDependency, 0, len(s.dependencies))
	for _, d := range s.dependencies {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) UpdateDependency(d *model.ServiceDependency) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dependencies[d.ID]; !ok {
		return ErrNotFound
	}
	s.dependencies[d.ID] = d
	return nil
}

func (s *MemoryStore) DeleteDependency(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dependencies[id]; !ok {
		return ErrNotFound
	}
	delete(s.dependencies, id)
	return nil
}
