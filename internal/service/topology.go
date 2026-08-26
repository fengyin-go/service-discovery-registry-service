package service

import (
	"sort"

	"serviceregistry/internal/model"
)

// TopologyNode 服务依赖拓扑节点。
type TopologyNode struct {
	Service    *model.Service   `json:"service"`
	Upstream   []*model.Service `json:"upstream"`
	Downstream []*model.Service `json:"downstream"`
}

// Topology 返回全部服务的依赖拓扑（含上下游）。
func (s *Service) Topology() []TopologyNode {
	services := s.store.ListServices()
	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })
	nodes := make([]TopologyNode, 0, len(services))
	for _, sv := range services {
		node := TopologyNode{
			Service:    sv,
			Upstream:   []*model.Service{},
			Downstream: []*model.Service{},
		}
		for _, d := range s.store.ListDependencies() {
			if d.ServiceID == sv.ID {
				if dep, err := s.store.GetService(d.DependsOnID); err == nil {
					node.Upstream = append(node.Upstream, dep)
				}
			}
			if d.DependsOnID == sv.ID {
				if dep, err := s.store.GetService(d.ServiceID); err == nil {
					node.Downstream = append(node.Downstream, dep)
				}
			}
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// DetectCycles 检测依赖图中的环，返回成环的服务名链。
func (s *Service) DetectCycles() [][]string {
	adj := make(map[string][]string)
	for _, d := range s.store.ListDependencies() {
		adj[d.ServiceID] = append(adj[d.ServiceID], d.DependsOnID)
	}
	const (
		white = 0 // 未访问
		gray  = 1 // 访问中
		black = 2 // 已完成
	)
	color := make(map[string]int)
	var cycles [][]string
	var stack []string

	var dfs func(node string)
	dfs = func(node string) {
		color[node] = gray
		stack = append(stack, node)
		for _, next := range adj[node] {
			if color[next] == white {
				dfs(next)
			} else if color[next] == gray {
				// 找到环：从 stack 中 next 的位置截取
				idx := -1
				for i, n := range stack {
					if n == next {
						idx = i
						break
					}
				}
				if idx >= 0 {
					cycle := append([]string{}, stack[idx:]...)
					cycle = append(cycle, next)
					cycles = append(cycles, s.serviceNames(cycle))
				}
			}
		}
		stack = stack[:len(stack)-1]
		color[node] = black
	}

	// 保证遍历顺序稳定
	nodes := make([]string, 0, len(adj))
	for k := range adj {
		nodes = append(nodes, k)
	}
	sort.Strings(nodes)
	for _, n := range nodes {
		if color[n] == white {
			dfs(n)
		}
	}
	return cycles
}

// serviceNames 将服务 ID 链转换为服务名链。
func (s *Service) serviceNames(ids []string) []string {
	names := make([]string, 0, len(ids))
	for _, id := range ids {
		if sv, err := s.store.GetService(id); err == nil {
			names = append(names, sv.Name)
		} else {
			names = append(names, id)
		}
	}
	return names
}
