package store

import "serviceregistry/internal/model"

func (s *MemoryStore) CreateHeartbeat(h *model.Heartbeat) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.heartbeats[h.ID] = h
	return nil
}

func (s *MemoryStore) ListHeartbeats() []*model.Heartbeat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Heartbeat, 0, len(s.heartbeats))
	for _, h := range s.heartbeats {
		list = append(list, h)
	}
	return list
}

func (s *MemoryStore) ListHeartbeatsByInstance(instanceID string) []*model.Heartbeat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Heartbeat, 0)
	for _, h := range s.heartbeats {
		if h.InstanceID == instanceID {
			list = append(list, h)
		}
	}
	return list
}
