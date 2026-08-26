package service

import (
	"sort"

	"serviceregistry/internal/model"
)

// InstanceMetrics 实例心跳延迟指标。
type InstanceMetrics struct {
	InstanceID string  `json:"instance_id"`
	Count      int     `json:"count"`
	AvgLatency float64 `json:"avg_latency_ms"`
	MaxLatency int64   `json:"max_latency_ms"`
	Timeouts   int     `json:"timeouts"`
}

// InstanceMetrics 计算某实例的心跳延迟指标。
func (s *Service) InstanceMetrics(instanceID string) InstanceMetrics {
	m := InstanceMetrics{InstanceID: instanceID}
	var total int64
	for _, h := range s.store.ListHeartbeatsByInstance(instanceID) {
		m.Count++
		total += h.LatencyMs
		if h.LatencyMs > m.MaxLatency {
			m.MaxLatency = h.LatencyMs
		}
		if h.Result == model.HeartbeatTimeout {
			m.Timeouts++
		}
	}
	if m.Count > 0 {
		m.AvgLatency = float64(total) / float64(m.Count)
	}
	return m
}

// AllInstanceMetrics 返回全部实例的延迟指标（按平均延迟降序）。
func (s *Service) AllInstanceMetrics() []InstanceMetrics {
	list := make([]InstanceMetrics, 0)
	for _, i := range s.store.ListInstances() {
		list = append(list, s.InstanceMetrics(i.ID))
	}
	sort.Slice(list, func(a, b int) bool {
		return list[a].AvgLatency > list[b].AvgLatency
	})
	return list
}
