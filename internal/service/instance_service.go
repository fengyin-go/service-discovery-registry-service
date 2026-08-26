package service

import (
	"sort"
	"time"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/idgen"
)

// RegisterInstance 注册实例。
func (s *Service) RegisterInstance(input model.Instance) (*model.Instance, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "所属服务不存在")
	}
	now := time.Now()
	inst := &model.Instance{
		ID:            idgen.Hex(),
		ServiceID:     input.ServiceID,
		Host:          input.Host,
		Port:          input.Port,
		Weight:        input.Weight,
		Status:        input.Status,
		Region:        input.Region,
		Zone:          input.Zone,
		Metadata:      input.Metadata,
		LastHeartbeat: now,
		CreatedAt:     now,
	}
	if err := s.store.CreateInstance(inst); err != nil {
		return nil, err
	}
	s.emitEvent(input.ServiceID, inst.ID, model.EventRegister, "实例注册: "+inst.Address())
	return inst, nil
}

// GetInstance 按 ID 获取实例。
func (s *Service) GetInstance(id string) (*model.Instance, error) {
	return s.store.GetInstance(id)
}

// ListInstances 按筛选条件分页查询实例。
func (s *Service) ListInstances(filter model.InstanceFilter, page, size int) ([]*model.Instance, int, error) {
	all := s.store.ListInstances()
	matched := make([]*model.Instance, 0, len(all))
	for _, i := range all {
		if filter.Match(i) {
			matched = append(matched, i)
		}
	}
	sort.Slice(matched, func(a, b int) bool {
		return matched[a].CreatedAt.Before(matched[b].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Instance{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// ListInstancesByService 返回某服务的全部实例。
func (s *Service) ListInstancesByService(serviceID string) []*model.Instance {
	return s.ListInstancesByFilter(model.InstanceFilter{ServiceID: serviceID})
}

// ListInstancesByFilter 返回命中筛选条件的实例（用于内部选择）。
func (s *Service) ListInstancesByFilter(filter model.InstanceFilter) []*model.Instance {
	all := s.store.ListInstances()
	matched := make([]*model.Instance, 0, len(all))
	for _, i := range all {
		if filter.Match(i) {
			matched = append(matched, i)
		}
	}
	return matched
}

// ChangeInstanceStatus 按状态机流转实例状态。
func (s *Service) ChangeInstanceStatus(id, to string) (*model.Instance, error) {
	existing, err := s.store.GetInstance(id)
	if err != nil {
		return nil, err
	}
	if to != model.InstanceStatusStarting && to != model.InstanceStatusUp && to != model.InstanceStatusDown {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	if !model.CanInstanceTransition(existing.Status, to) {
		return nil, model.NewValidationError("status", "不允许从 "+existing.Status+" 流转到 "+to)
	}
	existing.Status = to
	if err := s.store.UpdateInstance(existing); err != nil {
		return nil, err
	}
	s.emitEvent(existing.ServiceID, id, model.EventStatusChange, "实例状态变更为 "+to)
	return existing, nil
}

// DeregisterInstance 注销实例并清理其健康检查。
func (s *Service) DeregisterInstance(id string) error {
	inst, err := s.store.GetInstance(id)
	if err != nil {
		return err
	}
	if hc, err := s.store.GetHealthCheckByInstance(id); err == nil {
		_ = s.store.DeleteHealthCheck(hc.ID)
	}
	if err := s.store.DeleteInstance(id); err != nil {
		return err
	}
	s.emitEvent(inst.ServiceID, id, model.EventDeregister, "实例注销: "+inst.Address())
	return nil
}

// TouchInstance 更新实例心跳时间，若实例非 up 则提升为 up。
func (s *Service) TouchInstance(id string) error {
	inst, err := s.store.GetInstance(id)
	if err != nil {
		return err
	}
	inst.LastHeartbeat = time.Now()
	if inst.Status == model.InstanceStatusStarting || inst.Status == model.InstanceStatusDown {
		inst.Status = model.InstanceStatusUp
	}
	return s.store.UpdateInstance(inst)
}
