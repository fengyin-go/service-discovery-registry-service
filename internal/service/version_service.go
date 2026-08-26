package service

import (
	"serviceregistry/internal/model"
	"serviceregistry/pkg/semver"
)

// VersionSummary 某服务名的版本汇总。
type VersionSummary struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
	Latest   string   `json:"latest"`
}

// VersionSummary 返回某服务名的版本汇总。
func (s *Service) VersionSummary(name string) (VersionSummary, error) {
	versions := s.ListServiceVersions(name)
	if len(versions) == 0 {
		return VersionSummary{}, model.NewValidationError("name", "服务不存在")
	}
	return VersionSummary{Name: name, Versions: versions, Latest: versions[0]}, nil
}

// OutdatedServices 返回每个服务名下非最新版本的服务列表。
func (s *Service) OutdatedServices() []*model.Service {
	latestByVersion := make(map[string]string)
	for _, sv := range s.store.ListServices() {
		cur, ok := latestByVersion[sv.Name]
		if !ok || semver.Compare(sv.Version, cur) > 0 {
			latestByVersion[sv.Name] = sv.Version
		}
	}
	outdated := make([]*model.Service, 0)
	for _, sv := range s.store.ListServices() {
		if semver.Compare(sv.Version, latestByVersion[sv.Name]) < 0 {
			outdated = append(outdated, sv)
		}
	}
	return outdated
}
