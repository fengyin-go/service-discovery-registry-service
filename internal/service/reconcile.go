package service

import (
	"serviceregistry/internal/model"
	"serviceregistry/pkg/idgen"
)

// ReconcileResult 一致性巡检结果。
type ReconcileResult struct {
	MissingHealthChecks []string `json:"missing_health_checks"`
	OrphanInstances     []string `json:"orphan_instances"`
	CreatedChecks       int      `json:"created_checks"`
	RemovedInstances    int      `json:"removed_instances"`
}

// Reconcile 巡检注册表一致性：
//  1. 为缺少健康检查的实例补建默认 TCP 检查；
//  2. 清理所属服务不存在的孤儿实例。
func (s *Service) Reconcile() ReconcileResult {
	res := ReconcileResult{
		MissingHealthChecks: []string{},
		OrphanInstances:     []string{},
	}
	for _, i := range s.store.ListInstances() {
		if _, err := s.store.GetService(i.ServiceID); err != nil {
			res.OrphanInstances = append(res.OrphanInstances, i.Address())
			_ = s.store.DeleteInstance(i.ID)
			res.RemovedInstances++
			continue
		}
		if _, err := s.store.GetHealthCheckByInstance(i.ID); err != nil {
			res.MissingHealthChecks = append(res.MissingHealthChecks, i.Address())
			hc := &model.HealthCheck{
				ID:          idgen.Hex(),
				InstanceID:  i.ID,
				Type:        model.CheckTypeTCP,
				Target:      i.Address(),
				IntervalSec: 10,
				TimeoutSec:  3,
				Threshold:   3,
				Status:      model.CheckStatusHealthy,
			}
			if hcErr := s.store.CreateHealthCheck(hc); hcErr == nil {
				res.CreatedChecks++
			}
		}
	}
	return res
}
