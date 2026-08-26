package service

import (
	"time"

	"serviceregistry/internal/model"
)

// MaintenanceResult 一次维护巡检的结果。
type MaintenanceResult struct {
	ReapedInstances int      `json:"reaped_instances"`
	StaleInstances  []string `json:"stale_instances"`
}

// RunMaintenance 巡检：将心跳超时的 up 实例置为 down。
func (s *Service) RunMaintenance() MaintenanceResult {
	timeout := 30
	if s.cfg != nil && s.cfg.HeartbeatTimeoutSec > 0 {
		timeout = s.cfg.HeartbeatTimeoutSec
	}
	stale := s.StaleInstances(timeout)
	result := MaintenanceResult{StaleInstances: make([]string, 0)}
	for _, i := range stale {
		result.StaleInstances = append(result.StaleInstances, i.Address())
		if i.Status == model.InstanceStatusUp {
			_, _ = s.ChangeInstanceStatus(i.ID, model.InstanceStatusDown)
			result.ReapedInstances++
		}
	}
	return result
}

// InstanceHealthView 实例综合健康视图。
type InstanceHealthView struct {
	Instance      *model.Instance    `json:"instance"`
	HealthCheck   *model.HealthCheck `json:"health_check,omitempty"`
	LastHeartbeat time.Time          `json:"last_heartbeat"`
	Healthy       bool               `json:"healthy"`
}

// InstanceHealth 返回实例的综合健康视图。
func (s *Service) InstanceHealth(instanceID string) (*InstanceHealthView, error) {
	i, err := s.store.GetInstance(instanceID)
	if err != nil {
		return nil, err
	}
	view := &InstanceHealthView{
		Instance:      i,
		LastHeartbeat: i.LastHeartbeat,
		Healthy:       i.Status == model.InstanceStatusUp,
	}
	if hc, err := s.store.GetHealthCheckByInstance(instanceID); err == nil {
		view.HealthCheck = hc
		if hc.Status == model.CheckStatusUnhealthy {
			view.Healthy = false
		}
	}
	return view, nil
}

// UptimeSeconds 返回实例自创建以来的存活秒数。
func (s *Service) UptimeSeconds(instanceID string) int64 {
	i, err := s.store.GetInstance(instanceID)
	if err != nil {
		return 0
	}
	return int64(time.Since(i.CreatedAt).Seconds())
}
