package service

import (
	"strings"

	"serviceregistry/internal/model"
)

// SearchResult 全局搜索结果。
type SearchResult struct {
	Services  []*model.Service      `json:"services"`
	Instances []*model.Instance     `json:"instances"`
	Labels    []*model.ServiceLabel `json:"labels"`
}

// GlobalSearch 在服务、实例、标签中做关键字搜索。
func (s *Service) GlobalSearch(keyword string) SearchResult {
	kw := strings.ToLower(strings.TrimSpace(keyword))
	result := SearchResult{
		Services:  []*model.Service{},
		Instances: []*model.Instance{},
		Labels:    []*model.ServiceLabel{},
	}
	if kw == "" {
		return result
	}
	for _, sv := range s.store.ListServices() {
		if strings.Contains(strings.ToLower(sv.Name), kw) ||
			strings.Contains(strings.ToLower(sv.Version), kw) ||
			strings.Contains(strings.ToLower(sv.Owner), kw) {
			result.Services = append(result.Services, sv)
		}
	}
	for _, i := range s.store.ListInstances() {
		if strings.Contains(strings.ToLower(i.Host), kw) ||
			strings.Contains(strings.ToLower(i.Address()), kw) {
			result.Instances = append(result.Instances, i)
		}
	}
	for _, l := range s.store.ListServiceLabels() {
		if strings.Contains(strings.ToLower(l.Key), kw) ||
			strings.Contains(strings.ToLower(l.Value), kw) {
			result.Labels = append(result.Labels, l)
		}
	}
	return result
}
