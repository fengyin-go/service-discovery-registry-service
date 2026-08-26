package store

import "serviceregistry/internal/model"

func (s *MemoryStore) CreateRegistryConfig(c *model.RegistryConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.configs {
		if exist.Key == c.Key {
			return ErrConflict
		}
	}
	s.configs[c.ID] = c
	return nil
}

func (s *MemoryStore) GetRegistryConfig(id string) (*model.RegistryConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.configs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) GetRegistryConfigByKey(key string) (*model.RegistryConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.configs {
		if c.Key == key {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRegistryConfigs() []*model.RegistryConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RegistryConfig, 0, len(s.configs))
	for _, c := range s.configs {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateRegistryConfig(c *model.RegistryConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.configs[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.configs {
		if exist.ID != c.ID && exist.Key == c.Key {
			return ErrConflict
		}
	}
	s.configs[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteRegistryConfig(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.configs[id]; !ok {
		return ErrNotFound
	}
	delete(s.configs, id)
	return nil
}
