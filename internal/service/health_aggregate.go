package service

import (
	"sort"

	"serviceregistry/internal/model"
)

// ServiceHealth 服务健康汇总。
type ServiceHealth struct {
	ServiceID   string  `json:"service_id"`
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	Total       int     `json:"total"`
	Up          int     `json:"up"`
	Down        int     `json:"down"`
	Healthy     int     `json:"healthy"`
	Unhealthy   int     `json:"unhealthy"`
	HealthRatio float64 `json:"health_ratio"`
}

// ServiceHealthMatrix 返回全部服务的健康汇总。
func (s *Service) ServiceHealthMatrix() []ServiceHealth {
	list := make([]ServiceHealth, 0)
	for _, sv := range s.store.ListServices() {
		sh := ServiceHealth{
			ServiceID: sv.ID,
			Name:      sv.Name,
			Version:   sv.Version,
		}
		instanceHealthy := make(map[string]bool)
		for _, i := range s.store.ListInstances() {
			if i.ServiceID != sv.ID {
				continue
			}
			sh.Total++
			if i.Status == model.InstanceStatusUp {
				sh.Up++
			} else {
				sh.Down++
			}
			instanceHealthy[i.ID] = true
		}
		for _, h := range s.store.ListHealthChecks() {
			if !instanceHealthy[h.InstanceID] {
				continue
			}
			if h.Status == model.CheckStatusHealthy {
				sh.Healthy++
			} else {
				sh.Unhealthy++
			}
		}
		if sh.Total > 0 {
			sh.HealthRatio = float64(sh.Up) / float64(sh.Total)
		}
		list = append(list, sh)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
	return list
}
