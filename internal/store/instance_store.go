package store

import "serviceregistry/internal/model"

func (s *MemoryStore) CreateInstance(i *model.Instance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.instances {
		if exist.ServiceID == i.ServiceID && exist.Host == i.Host && exist.Port == i.Port {
			return ErrConflict
		}
	}
	s.instances[i.ID] = i
	return nil
}

func (s *MemoryStore) GetInstance(id string) (*model.Instance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.instances[id]
	if !ok {
		return nil, ErrNotFound
	}
	return i, nil
}

func (s *MemoryStore) ListInstances() []*model.Instance {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Instance, 0, len(s.instances))
	for _, i := range s.instances {
		list = append(list, i)
	}
	return list
}

func (s *MemoryStore) UpdateInstance(i *model.Instance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.instances[i.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.instances {
		if exist.ID != i.ID && exist.ServiceID == i.ServiceID &&
			exist.Host == i.Host && exist.Port == i.Port {
			return ErrConflict
		}
	}
	s.instances[i.ID] = i
	return nil
}

func (s *MemoryStore) DeleteInstance(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.instances[id]; !ok {
		return ErrNotFound
	}
	delete(s.instances, id)
	return nil
}
