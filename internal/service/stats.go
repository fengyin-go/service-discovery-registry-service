package service

import (
	"sort"

	"serviceregistry/internal/model"
)

// RegistryStats 注册中心整体统计。
type RegistryStats struct {
	TotalServices   int            `json:"total_services"`
	TotalInstances  int            `json:"total_instances"`
	ByServiceStatus map[string]int `json:"by_service_status"`
	ByInstanceStatus map[string]int `json:"by_instance_status"`
	HealthyRatio    float64        `json:"healthy_ratio"`
	TotalEvents     int            `json:"total_events"`
	TotalDependencies int          `json:"total_dependencies"`
}

// StatsRegistry 返回注册中心整体统计。
func (s *Service) StatsRegistry() RegistryStats {
	stats := RegistryStats{
		ByServiceStatus:  make(map[string]int),
		ByInstanceStatus: make(map[string]int),
	}
	healthy := 0
	for _, sv := range s.store.ListServices() {
		stats.TotalServices++
		stats.ByServiceStatus[sv.Status]++
	}
	for _, i := range s.store.ListInstances() {
		stats.TotalInstances++
		stats.ByInstanceStatus[i.Status]++
		if i.Status == model.InstanceStatusUp {
			healthy++
		}
	}
	if stats.TotalInstances > 0 {
		stats.HealthyRatio = float64(healthy) / float64(stats.TotalInstances)
	}
	stats.TotalEvents = len(s.store.ListRegistryEvents())
	stats.TotalDependencies = len(s.store.ListDependencies())
	return stats
}

// ServiceStats 单服务统计。
type ServiceStats struct {
	Service         *model.Service `json:"service"`
	TotalInstances  int            `json:"total_instances"`
	UpInstances     int            `json:"up_instances"`
	DownInstances   int            `json:"down_instances"`
	HealthyChecks   int            `json:"healthy_checks"`
	UnhealthyChecks int            `json:"unhealthy_checks"`
}

// StatsService 返回某服务的实例与健康统计。
func (s *Service) StatsService(serviceID string) (*ServiceStats, error) {
	sv, err := s.store.GetService(serviceID)
	if err != nil {
		return nil, err
	}
	stats := &ServiceStats{Service: sv}
	for _, i := range s.store.ListInstances() {
		if i.ServiceID != serviceID {
			continue
		}
		stats.TotalInstances++
		switch i.Status {
		case model.InstanceStatusUp:
			stats.UpInstances++
		case model.InstanceStatusDown:
			stats.DownInstances++
		}
	}
	for _, h := range s.store.ListHealthChecks() {
		inst, err := s.store.GetInstance(h.InstanceID)
		if err != nil || inst.ServiceID != serviceID {
			continue
		}
		if h.Status == model.CheckStatusHealthy {
			stats.HealthyChecks++
		} else {
			stats.UnhealthyChecks++
		}
	}
	return stats, nil
}

// DependencyRank 依赖调用排行项。
type DependencyRank struct {
	Dependency *model.ServiceDependency `json:"dependency"`
	CallCount  int64                    `json:"call_count"`
}

// TopDependencies 返回调用次数最高的 N 个依赖。
func (s *Service) TopDependencies(n int) []DependencyRank {
	if n <= 0 {
		n = 10
	}
	list := make([]DependencyRank, 0)
	for _, d := range s.store.ListDependencies() {
		list = append(list, DependencyRank{Dependency: d, CallCount: d.CallCount})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].CallCount > list[j].CallCount })
	if len(list) > n {
		list = list[:n]
	}
	return list
}
