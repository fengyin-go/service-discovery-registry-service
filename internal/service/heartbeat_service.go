package service

import (
	"sort"
	"time"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/idgen"
)

// SendHeartbeat 记录一次实例心跳并刷新其最后心跳时间。
func (s *Service) SendHeartbeat(instanceID string, latencyMs int64, result string) (*model.Heartbeat, error) {
	if _, err := s.store.GetInstance(instanceID); err != nil {
		return nil, err
	}
	now := time.Now()
	h := &model.Heartbeat{
		ID:         idgen.Hex(),
		InstanceID: instanceID,
		SentAt:     now,
		ReceivedAt: now,
		LatencyMs:  latencyMs,
		Result:     result,
	}
	if err := h.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateHeartbeat(h); err != nil {
		return nil, err
	}
	if err := s.TouchInstance(instanceID); err != nil {
		return nil, err
	}
	return h, nil
}

// ListHeartbeats 按实例与结果筛选心跳记录。
func (s *Service) ListHeartbeats(filter model.HeartbeatFilter) []*model.Heartbeat {
	var list []*model.Heartbeat
	if filter.InstanceID != "" {
		list = s.store.ListHeartbeatsByInstance(filter.InstanceID)
	} else {
		list = s.store.ListHeartbeats()
	}
	matched := make([]*model.Heartbeat, 0, len(list))
	for _, h := range list {
		if filter.Match(h) {
			matched = append(matched, h)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].ReceivedAt.After(matched[j].ReceivedAt) })
	return matched
}

// StaleInstances 返回最后心跳早于 timeoutSec 秒的实例。
func (s *Service) StaleInstances(timeoutSec int) []*model.Instance {
	deadline := time.Now().Add(-time.Duration(timeoutSec) * time.Second)
	stale := make([]*model.Instance, 0)
	for _, i := range s.store.ListInstances() {
		if i.LastHeartbeat.Before(deadline) {
			stale = append(stale, i)
		}
	}
	return stale
}

// ReapStaleInstances 将心跳超时的实例置为 down。
func (s *Service) ReapStaleInstances(timeoutSec int) int {
	reaped := 0
	for _, i := range s.StaleInstances(timeoutSec) {
		if i.Status == model.InstanceStatusUp {
			_, _ = s.ChangeInstanceStatus(i.ID, model.InstanceStatusDown)
			reaped++
		}
	}
	return reaped
}
