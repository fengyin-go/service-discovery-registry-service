package service

import (
	"sort"
	"time"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/idgen"
)

// CreateLoadBalanceStrategy 为服务配置负载均衡策略。
func (s *Service) CreateLoadBalanceStrategy(input model.LoadBalanceStrategy) (*model.LoadBalanceStrategy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "服务不存在")
	}
	now := time.Now()
	l := &model.LoadBalanceStrategy{
		ID:        idgen.Hex(),
		ServiceID: input.ServiceID,
		Algorithm: input.Algorithm,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateLoadBalanceStrategy(l); err != nil {
		return nil, err
	}
	return l, nil
}

// GetLoadBalanceByService 获取服务的负载均衡策略。
func (s *Service) GetLoadBalanceByService(serviceID string) (*model.LoadBalanceStrategy, error) {
	return s.store.GetLoadBalanceByService(serviceID)
}

// GetLoadBalanceStrategy 按 ID 获取负载均衡策略。
func (s *Service) GetLoadBalanceStrategy(id string) (*model.LoadBalanceStrategy, error) {
	return s.store.GetLoadBalanceStrategy(id)
}

// ListLoadBalanceStrategies 返回全部负载均衡策略。
func (s *Service) ListLoadBalanceStrategies() []*model.LoadBalanceStrategy {
	list := s.store.ListLoadBalanceStrategies()
	sort.Slice(list, func(i, j int) bool { return list[i].ServiceID < list[j].ServiceID })
	return list
}

// UpdateLoadBalanceStrategy 更新负载均衡算法。
func (s *Service) UpdateLoadBalanceStrategy(id string, algorithm string) (*model.LoadBalanceStrategy, error) {
	existing, err := s.store.GetLoadBalanceStrategy(id)
	if err != nil {
		return nil, err
	}
	existing.Algorithm = algorithm
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateLoadBalanceStrategy(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteLoadBalanceStrategy 删除负载均衡策略。
func (s *Service) DeleteLoadBalanceStrategy(id string) error {
	return s.store.DeleteLoadBalanceStrategy(id)
}
