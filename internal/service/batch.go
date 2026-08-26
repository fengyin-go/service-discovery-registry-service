package service

import "serviceregistry/internal/model"

// BatchRegisterInstances 批量注册实例，跳过失败项并返回成功数量。
func (s *Service) BatchRegisterInstances(serviceID string, items []model.Instance) (int, error) {
	if _, err := s.store.GetService(serviceID); err != nil {
		return 0, model.NewValidationError("service_id", "服务不存在")
	}
	ok := 0
	for _, item := range items {
		item.ServiceID = serviceID
		if _, err := s.RegisterInstance(item); err == nil {
			ok++
		}
	}
	return ok, nil
}

// BatchDeregister 批量注销实例，返回成功数量。
func (s *Service) BatchDeregister(instanceIDs []string) int {
	ok := 0
	for _, id := range instanceIDs {
		if err := s.DeregisterInstance(id); err == nil {
			ok++
		}
	}
	return ok
}

// BulkHealthResult 单条健康上报。
type BulkHealthResult struct {
	InstanceID string `json:"instance_id"`
	Status     string `json:"status"`
}

// BulkHealthReport 批量上报健康检查结果，返回成功数量。
func (s *Service) BulkHealthReport(reports []BulkHealthResult) int {
	ok := 0
	for _, r := range reports {
		if _, err := s.ReportCheckResult(r.InstanceID, r.Status); err == nil {
			ok++
		}
	}
	return ok
}
