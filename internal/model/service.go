package model

import (
	"strings"
	"time"
)

// 服务状态常量。
const (
	ServiceStatusUp   = "up"
	ServiceStatusDown = "down"
)

// Service 服务实体。
type Service struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Owner       string    `json:"owner"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 规范化并校验服务字段。
func (s *Service) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Version = strings.TrimSpace(s.Version)
	s.Description = strings.TrimSpace(s.Description)
	s.Owner = strings.TrimSpace(s.Owner)
	if s.Name == "" {
		return NewValidationError("name", "服务名称不能为空")
	}
	if s.Version == "" {
		return NewValidationError("version", "服务版本不能为空")
	}
	if s.Status == "" {
		s.Status = ServiceStatusUp
	}
	if s.Status != ServiceStatusUp && s.Status != ServiceStatusDown {
		return NewValidationError("status", "服务状态不合法")
	}
	return nil
}

// ServiceFilter 服务列表筛选条件。
type ServiceFilter struct {
	Keyword string
	Status  string
	Owner   string
}

// Match 判断服务是否命中筛选条件。
func (f ServiceFilter) Match(s *Service) bool {
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Owner != "" && s.Owner != f.Owner {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) &&
			!strings.Contains(strings.ToLower(s.Version), k) {
			return false
		}
	}
	return true
}
