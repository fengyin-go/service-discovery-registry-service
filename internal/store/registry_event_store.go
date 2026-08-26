package store

import "serviceregistry/internal/model"

func (s *MemoryStore) CreateRegistryEvent(e *model.RegistryEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.ID] = e
	return nil
}

func (s *MemoryStore) ListRegistryEvents() []*model.RegistryEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RegistryEvent, 0, len(s.events))
	for _, e := range s.events {
		list = append(list, e)
	}
	return list
}
