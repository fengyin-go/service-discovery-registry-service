package service

import (
	"math/rand"
	"sort"
	"sync"

	"serviceregistry/internal/model"
	"serviceregistry/internal/store"
)

// selector 维护负载均衡所需的运行时状态。
type selector struct {
	mu    sync.Mutex
	rr    map[string]int
	conns map[string]int
}

func newSelector() *selector {
	return &selector{
		rr:    make(map[string]int),
		conns: make(map[string]int),
	}
}

// upInstances 返回某服务的在线（up）实例，仅 up 状态可被选择。
func (s *Service) upInstances(serviceID string) []*model.Instance {
	all := s.store.ListInstances()
	up := make([]*model.Instance, 0, len(all))
	for _, i := range all {
		if i.ServiceID == serviceID && i.Status == model.InstanceStatusUp {
			up = append(up, i)
		}
	}
	sort.Slice(up, func(a, b int) bool { return up[a].ID < up[b].ID })
	return up
}

// SelectInstance 根据服务配置的负载均衡算法选取一个实例。
func (s *Service) SelectInstance(serviceID string) (*model.Instance, error) {
	if _, err := s.store.GetService(serviceID); err != nil {
		return nil, err
	}
	candidates := s.upInstances(serviceID)
	if len(candidates) == 0 {
		return nil, store.ErrNotFound
	}
	algorithm := model.AlgorithmRoundRobin
	if lb, err := s.store.GetLoadBalanceByService(serviceID); err == nil {
		algorithm = lb.Algorithm
	}
	switch algorithm {
	case model.AlgorithmRandom:
		return candidates[rand.Intn(len(candidates))], nil
	case model.AlgorithmWeighted:
		return s.weightedSelect(candidates), nil
	case model.AlgorithmLeastConn:
		return s.leastConnSelect(candidates), nil
	default:
		return s.roundRobinSelect(serviceID, candidates), nil
	}
}

func (s *Service) roundRobinSelect(serviceID string, candidates []*model.Instance) *model.Instance {
	s.sel.mu.Lock()
	defer s.sel.mu.Unlock()
	idx := s.sel.rr[serviceID] % len(candidates)
	s.sel.rr[serviceID] = idx + 1
	return candidates[idx]
}

func (s *Service) weightedSelect(candidates []*model.Instance) *model.Instance {
	total := 0
	for _, c := range candidates {
		total += c.Weight
	}
	if total <= 0 {
		return candidates[0]
	}
	roll := rand.Intn(total)
	for _, c := range candidates {
		roll -= c.Weight
		if roll < 0 {
			return c
		}
	}
	return candidates[len(candidates)-1]
}

func (s *Service) leastConnSelect(candidates []*model.Instance) *model.Instance {
	s.sel.mu.Lock()
	defer s.sel.mu.Unlock()
	best := candidates[0]
	for _, c := range candidates {
		if s.sel.conns[c.ID] < s.sel.conns[best.ID] {
			best = c
		}
	}
	s.sel.conns[best.ID]++
	return best
}

// ReleaseInstance 释放实例的一个连接计数。
func (s *Service) ReleaseInstance(instanceID string) {
	s.sel.mu.Lock()
	defer s.sel.mu.Unlock()
	if s.sel.conns[instanceID] > 0 {
		s.sel.conns[instanceID]--
	}
}
