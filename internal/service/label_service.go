package service

import (
	"sort"
	"time"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/idgen"
)

// CreateLabel 为服务打标签。
func (s *Service) CreateLabel(serviceID, key, value string) (*model.ServiceLabel, error) {
	if _, err := s.store.GetService(serviceID); err != nil {
		return nil, model.NewValidationError("service_id", "服务不存在")
	}
	l := &model.ServiceLabel{
		ID:        idgen.Hex(),
		ServiceID: serviceID,
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
	}
	if err := l.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateServiceLabel(l); err != nil {
		return nil, err
	}
	return l, nil
}

// GetLabel 按 ID 获取标签。
func (s *Service) GetLabel(id string) (*model.ServiceLabel, error) {
	return s.store.GetServiceLabel(id)
}

// ListLabels 按筛选条件返回标签。
func (s *Service) ListLabels(filter model.LabelFilter) []*model.ServiceLabel {
	all := s.store.ListServiceLabels()
	matched := make([]*model.ServiceLabel, 0, len(all))
	for _, l := range all {
		if filter.Match(l) {
			matched = append(matched, l)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].ServiceID != matched[j].ServiceID {
			return matched[i].ServiceID < matched[j].ServiceID
		}
		return matched[i].Key < matched[j].Key
	})
	return matched
}

// ListLabelsByService 返回某服务的全部标签。
func (s *Service) ListLabelsByService(serviceID string) []*model.ServiceLabel {
	return s.ListLabels(model.LabelFilter{ServiceID: serviceID})
}

// ListServicesByLabel 返回命中某标签键值对的服务。
func (s *Service) ListServicesByLabel(key, value string) []*model.Service {
	matched := make([]*model.Service, 0)
	seen := make(map[string]bool)
	for _, l := range s.store.ListServiceLabels() {
		if l.Key != key || (value != "" && l.Value != value) {
			continue
		}
		if seen[l.ServiceID] {
			continue
		}
		if sv, err := s.store.GetService(l.ServiceID); err == nil {
			seen[l.ServiceID] = true
			matched = append(matched, sv)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].Name < matched[j].Name })
	return matched
}

// UpdateLabel 更新标签值。
func (s *Service) UpdateLabel(id, value string) (*model.ServiceLabel, error) {
	l, err := s.store.GetServiceLabel(id)
	if err != nil {
		return nil, err
	}
	l.Value = value
	if err := l.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateServiceLabel(l); err != nil {
		return nil, err
	}
	return l, nil
}

// DeleteLabel 删除标签。
func (s *Service) DeleteLabel(id string) error {
	return s.store.DeleteServiceLabel(id)
}
