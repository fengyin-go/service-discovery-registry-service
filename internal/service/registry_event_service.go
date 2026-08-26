package service

import (
	"sort"
	"time"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/idgen"
)

// emitEvent 写入一条注册中心事件。
func (s *Service) emitEvent(serviceID, instanceID, eventType, detail string) {
	e := &model.RegistryEvent{
		ID:         idgen.Hex(),
		ServiceID:  serviceID,
		InstanceID: instanceID,
		EventType:  eventType,
		Detail:     detail,
		CreatedAt:  time.Now(),
	}
	_ = s.store.CreateRegistryEvent(e)
}

// ListEvents 按筛选条件分页查询事件（时间倒序）。
func (s *Service) ListEvents(filter model.EventFilter, page, size int) ([]*model.RegistryEvent, int, error) {
	all := s.store.ListRegistryEvents()
	matched := make([]*model.RegistryEvent, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RegistryEvent{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
