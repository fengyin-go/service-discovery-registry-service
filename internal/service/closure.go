package service

import (
	"sort"

	"serviceregistry/internal/model"
)

// DependencyClosure 依赖传递闭包。
type DependencyClosure struct {
	Service    *model.Service   `json:"service"`
	Upstream   []*model.Service `json:"upstream_closure"`
	Downstream []*model.Service `json:"downstream_closure"`
}

// DependencyClosure 计算某服务的全部传递上游与下游依赖。
func (s *Service) DependencyClosure(serviceID string) (*DependencyClosure, error) {
	sv, err := s.store.GetService(serviceID)
	if err != nil {
		return nil, err
	}
	upAdj := make(map[string][]string)
	downAdj := make(map[string][]string)
	for _, d := range s.store.ListDependencies() {
		upAdj[d.ServiceID] = append(upAdj[d.ServiceID], d.DependsOnID)
		downAdj[d.DependsOnID] = append(downAdj[d.DependsOnID], d.ServiceID)
	}
	closure := &DependencyClosure{
		Service:    sv,
		Upstream:   s.closureServices(serviceID, upAdj),
		Downstream: s.closureServices(serviceID, downAdj),
	}
	return closure, nil
}

// closureServices 从 start 出发沿邻接表做传递闭包，返回命中的服务。
func (s *Service) closureServices(start string, adj map[string][]string) []*model.Service {
	visited := map[string]bool{start: true}
	queue := []string{start}
	result := make([]*model.Service, 0)
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range adj[cur] {
			if visited[next] {
				continue
			}
			visited[next] = true
			if sv, err := s.store.GetService(next); err == nil {
				result = append(result, sv)
			}
			queue = append(queue, next)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}
