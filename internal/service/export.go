package service

import (
	"serviceregistry/internal/model"
	"serviceregistry/pkg/idgen"
)

// Snapshot 注册中心全量快照。
type Snapshot struct {
	Services     []*model.Service             `json:"services"`
	Instances    []*model.Instance            `json:"instances"`
	HealthChecks []*model.HealthCheck         `json:"health_checks"`
	LoadBalances []*model.LoadBalanceStrategy `json:"load_balances"`
	Dependencies []*model.ServiceDependency   `json:"dependencies"`
}

// ExportSnapshot 导出注册中心全量快照。
func (s *Service) ExportSnapshot() Snapshot {
	return Snapshot{
		Services:     s.store.ListServices(),
		Instances:    s.store.ListInstances(),
		HealthChecks: s.store.ListHealthChecks(),
		LoadBalances: s.store.ListLoadBalanceStrategies(),
		Dependencies: s.store.ListDependencies(),
	}
}

// ImportSnapshot 导入快照并重建索引（跳过重复项）。
func (s *Service) ImportSnapshot(snap Snapshot) (map[string]int, error) {
	imported := map[string]int{
		"services":     0,
		"instances":    0,
		"health_checks": 0,
		"load_balances": 0,
		"dependencies": 0,
	}
	for _, sv := range snap.Services {
		if sv.ID == "" {
			sv.ID = idgen.Hex()
		}
		if err := s.store.CreateService(sv); err == nil {
			imported["services"]++
		}
	}
	for _, i := range snap.Instances {
		if i.ID == "" {
			i.ID = idgen.Hex()
		}
		if err := s.store.CreateInstance(i); err == nil {
			imported["instances"]++
		}
	}
	for _, h := range snap.HealthChecks {
		if h.ID == "" {
			h.ID = idgen.Hex()
		}
		if err := s.store.CreateHealthCheck(h); err == nil {
			imported["health_checks"]++
		}
	}
	for _, l := range snap.LoadBalances {
		if l.ID == "" {
			l.ID = idgen.Hex()
		}
		if err := s.store.CreateLoadBalanceStrategy(l); err == nil {
			imported["load_balances"]++
		}
	}
	for _, d := range snap.Dependencies {
		if d.ID == "" {
			d.ID = idgen.Hex()
		}
		if err := s.store.CreateDependency(d); err == nil {
			imported["dependencies"]++
		}
	}
	s.emitEvent("", "", model.EventRegister, "导入快照")
	return imported, nil
}
