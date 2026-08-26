package service

import (
	"sort"

	"serviceregistry/internal/model"
)

// RecentEvents 返回某服务最近的 N 条事件。
func (s *Service) RecentEvents(serviceID string, n int) []*model.RegistryEvent {
	if n <= 0 {
		n = 20
	}
	list := make([]*model.RegistryEvent, 0)
	for _, e := range s.store.ListRegistryEvents() {
		if serviceID == "" || e.ServiceID == serviceID {
			list = append(list, e)
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
	if len(list) > n {
		list = list[:n]
	}
	return list
}

// EventsForInstance 返回某实例的全部事件。
func (s *Service) EventsForInstance(instanceID string) []*model.RegistryEvent {
	list := make([]*model.RegistryEvent, 0)
	for _, e := range s.store.ListRegistryEvents() {
		if e.InstanceID == instanceID {
			list = append(list, e)
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
	return list
}

// EventTypeCounts 返回各事件类型计数（按数量降序）。
type EventTypeCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// EventTypeCounts 返回各事件类型计数。
func (s *Service) EventTypeCounts() []EventTypeCount {
	counts := make(map[string]int)
	for _, e := range s.store.ListRegistryEvents() {
		counts[e.EventType]++
	}
	list := make([]EventTypeCount, 0, len(counts))
	for typ, cnt := range counts {
		list = append(list, EventTypeCount{Type: typ, Count: cnt})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Count > list[j].Count })
	return list
}
