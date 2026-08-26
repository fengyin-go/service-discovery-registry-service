package service

import (
	"sort"
	"time"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/idgen"
)

// CreateDependency 声明服务依赖关系。
func (s *Service) CreateDependency(input model.ServiceDependency) (*model.ServiceDependency, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "源服务不存在")
	}
	if _, err := s.store.GetService(input.DependsOnID); err != nil {
		return nil, model.NewValidationError("depends_on_id", "被依赖服务不存在")
	}
	now := time.Now()
	d := &model.ServiceDependency{
		ID:          idgen.Hex(),
		ServiceID:   input.ServiceID,
		DependsOnID: input.DependsOnID,
		Protocol:    input.Protocol,
		CallCount:   0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateDependency(d); err != nil {
		return nil, err
	}
	s.emitEvent(input.ServiceID, "", model.EventDependencyAdd, "新增依赖 "+input.DependsOnID)
	return d, nil
}

// GetDependency 按 ID 获取依赖。
func (s *Service) GetDependency(id string) (*model.ServiceDependency, error) {
	return s.store.GetDependency(id)
}

// ListDependencies 返回全部依赖关系。
func (s *Service) ListDependencies() []*model.ServiceDependency {
	list := s.store.ListDependencies()
	sort.Slice(list, func(i, j int) bool { return list[i].ServiceID < list[j].ServiceID })
	return list
}

// DependenciesOf 返回某服务依赖的所有服务。
func (s *Service) DependenciesOf(serviceID string) []*model.ServiceDependency {
	list := make([]*model.ServiceDependency, 0)
	for _, d := range s.store.ListDependencies() {
		if d.ServiceID == serviceID {
			list = append(list, d)
		}
	}
	return list
}

// DependentsOf 返回依赖某服务的所有服务。
func (s *Service) DependentsOf(serviceID string) []*model.ServiceDependency {
	list := make([]*model.ServiceDependency, 0)
	for _, d := range s.store.ListDependencies() {
		if d.DependsOnID == serviceID {
			list = append(list, d)
		}
	}
	return list
}

// IncrementCallCount 累加依赖调用次数。
func (s *Service) IncrementCallCount(id string) (*model.ServiceDependency, error) {
	d, err := s.store.GetDependency(id)
	if err != nil {
		return nil, err
	}
	d.CallCount++
	d.UpdatedAt = time.Now()
	if err := s.store.UpdateDependency(d); err != nil {
		return nil, err
	}
	return d, nil
}

// DeleteDependency 删除依赖关系。
func (s *Service) DeleteDependency(id string) error {
	return s.store.DeleteDependency(id)
}
