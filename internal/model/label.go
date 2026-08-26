package model

import (
	"strings"
	"time"
)

// ServiceLabel 服务标签，用于分组与检索。
type ServiceLabel struct {
	ID        string    `json:"id"`
	ServiceID string    `json:"service_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 规范化并校验标签字段。
func (l *ServiceLabel) Validate() error {
	l.ServiceID = strings.TrimSpace(l.ServiceID)
	l.Key = strings.TrimSpace(l.Key)
	l.Value = strings.TrimSpace(l.Value)
	if l.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if l.Key == "" {
		return NewValidationError("key", "标签键不能为空")
	}
	return nil
}

// LabelFilter 标签筛选条件。
type LabelFilter struct {
	ServiceID string
	Key       string
	Value     string
}

// Match 判断标签是否命中筛选条件。
func (f LabelFilter) Match(l *ServiceLabel) bool {
	if f.ServiceID != "" && l.ServiceID != f.ServiceID {
		return false
	}
	if f.Key != "" && l.Key != f.Key {
		return false
	}
	if f.Value != "" && l.Value != f.Value {
		return false
	}
	return true
}
