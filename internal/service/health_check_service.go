package service

import (
	"sort"
	"time"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/idgen"
)

// CreateHealthCheck 为实例配置健康检查。
func (s *Service) CreateHealthCheck(input model.HealthCheck) (*model.HealthCheck, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetInstance(input.InstanceID); err != nil {
		return nil, model.NewValidationError("instance_id", "实例不存在")
	}
	h := &model.HealthCheck{
		ID:            idgen.Hex(),
		InstanceID:    input.InstanceID,
		Type:          input.Type,
		Target:        input.Target,
		IntervalSec:   input.IntervalSec,
		TimeoutSec:    input.TimeoutSec,
		Threshold:     input.Threshold,
		Status:        input.Status,
		LastCheckedAt: time.Now(),
		CreatedAt:     time.Now(),
	}
	if err := s.store.CreateHealthCheck(h); err != nil {
		return nil, err
	}
	return h, nil
}

// GetHealthCheck 按 ID 获取健康检查。
func (s *Service) GetHealthCheck(id string) (*model.HealthCheck, error) {
	return s.store.GetHealthCheck(id)
}

// GetHealthCheckByInstance 按实例获取健康检查。
func (s *Service) GetHealthCheckByInstance(instanceID string) (*model.HealthCheck, error) {
	return s.store.GetHealthCheckByInstance(instanceID)
}

// ListHealthChecks 返回全部健康检查。
func (s *Service) ListHealthChecks() []*model.HealthCheck {
	list := s.store.ListHealthChecks()
	sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.Before(list[j].CreatedAt) })
	return list
}

// UpdateHealthCheck 更新健康检查配置。
func (s *Service) UpdateHealthCheck(id string, input model.HealthCheck) (*model.HealthCheck, error) {
	existing, err := s.store.GetHealthCheck(id)
	if err != nil {
		return nil, err
	}
	existing.Type = input.Type
	existing.Target = input.Target
	existing.IntervalSec = input.IntervalSec
	existing.TimeoutSec = input.TimeoutSec
	existing.Threshold = input.Threshold
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateHealthCheck(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteHealthCheck 删除健康检查。
func (s *Service) DeleteHealthCheck(id string) error {
	return s.store.DeleteHealthCheck(id)
}

// ReportCheckResult 上报一次健康检查结果，状态变化时联动实例与事件。
func (s *Service) ReportCheckResult(instanceID, status string) (*model.HealthCheck, error) {
	if status != model.CheckStatusHealthy && status != model.CheckStatusUnhealthy {
		return nil, model.NewValidationError("status", "健康状态不合法")
	}
	h, err := s.store.GetHealthCheckByInstance(instanceID)
	if err != nil {
		return nil, err
	}
	if h.Status != status {
		h.Status = status
		inst, _ := s.store.GetInstance(instanceID)
		if inst != nil {
			if status == model.CheckStatusUnhealthy {
				_, _ = s.ChangeInstanceStatus(instanceID, model.InstanceStatusDown)
			} else if inst.Status == model.InstanceStatusDown {
				_, _ = s.ChangeInstanceStatus(instanceID, model.InstanceStatusUp)
			}
		}
		s.emitEvent(inst.ServiceID, instanceID, model.EventHealthChange, "健康检查变更为 "+status)
	}
	h.LastCheckedAt = time.Now()
	if err := s.store.UpdateHealthCheck(h); err != nil {
		return nil, err
	}
	return h, nil
}
