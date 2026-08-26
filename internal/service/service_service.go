package service

import (
	"sort"
	"time"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/idgen"
	"serviceregistry/pkg/semver"
)

// CreateService 注册服务。
func (s *Service) CreateService(input model.Service) (*model.Service, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	sv := &model.Service{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Version:     input.Version,
		Description: input.Description,
		Status:      input.Status,
		Owner:       input.Owner,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateService(sv); err != nil {
		return nil, err
	}
	s.emitEvent(sv.ID, "", model.EventRegister, "服务注册: "+sv.Name+"@"+sv.Version)
	return sv, nil
}

// GetService 按 ID 获取服务。
func (s *Service) GetService(id string) (*model.Service, error) {
	return s.store.GetService(id)
}

// ListServices 按筛选条件分页查询服务。
func (s *Service) ListServices(filter model.ServiceFilter, page, size int) ([]*model.Service, int, error) {
	all := s.store.ListServices()
	matched := make([]*model.Service, 0, len(all))
	for _, sv := range all {
		if filter.Match(sv) {
			matched = append(matched, sv)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].Name != matched[j].Name {
			return matched[i].Name < matched[j].Name
		}
		return semver.Compare(matched[i].Version, matched[j].Version) > 0
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Service{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateService 更新服务可编辑字段。
func (s *Service) UpdateService(id string, input model.Service) (*model.Service, error) {
	existing, err := s.store.GetService(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Version = input.Version
	existing.Description = input.Description
	existing.Owner = input.Owner
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateService(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// ChangeServiceStatus 切换服务上下线状态。
func (s *Service) ChangeServiceStatus(id, to string) (*model.Service, error) {
	existing, err := s.store.GetService(id)
	if err != nil {
		return nil, err
	}
	if to != model.ServiceStatusUp && to != model.ServiceStatusDown {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	if existing.Status != to {
		existing.Status = to
		existing.UpdatedAt = time.Now()
		if err := s.store.UpdateService(existing); err != nil {
			return nil, err
		}
		s.emitEvent(id, "", model.EventStatusChange, "服务状态变更为 "+to)
	}
	return existing, nil
}

// DeleteService 删除服务及其全部实例与健康检查。
func (s *Service) DeleteService(id string) error {
	if _, err := s.store.GetService(id); err != nil {
		return err
	}
	for _, inst := range s.store.ListInstances() {
		if inst.ServiceID == id {
			_ = s.DeregisterInstance(inst.ID)
		}
	}
	for _, lb := range s.store.ListLoadBalanceStrategies() {
		if lb.ServiceID == id {
			_ = s.store.DeleteLoadBalanceStrategy(lb.ID)
		}
	}
	if err := s.store.DeleteService(id); err != nil {
		return err
	}
	s.emitEvent(id, "", model.EventDeregister, "服务注销")
	return nil
}

// ListServiceVersions 返回某服务名下的全部版本（按语义化版本降序）。
func (s *Service) ListServiceVersions(name string) []string {
	versions := make([]string, 0)
	for _, sv := range s.store.ListServices() {
		if sv.Name == name {
			versions = append(versions, sv.Version)
		}
	}
	sort.Slice(versions, func(i, j int) bool {
		return semver.Compare(versions[i], versions[j]) > 0
	})
	return versions
}

// LatestService 返回某服务名下的最新版本服务。
func (s *Service) LatestService(name string) (*model.Service, error) {
	var latest *model.Service
	for _, sv := range s.store.ListServices() {
		if sv.Name != name {
			continue
		}
		if latest == nil || semver.Compare(sv.Version, latest.Version) > 0 {
			latest = sv
		}
	}
	if latest == nil {
		return nil, model.NewValidationError("name", "服务不存在")
	}
	return latest, nil
}
