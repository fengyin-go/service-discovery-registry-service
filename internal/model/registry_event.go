package model

import (
	"strings"
	"time"
)

// 注册事件类型常量。
const (
	EventRegister      = "register"
	EventDeregister    = "deregister"
	EventStatusChange  = "status_change"
	EventHealthChange  = "health_change"
	EventDependencyAdd = "dependency_add"
)

// RegistryEvent 注册中心事件。
type RegistryEvent struct {
	ID         string    `json:"id"`
	ServiceID  string    `json:"service_id"`
	InstanceID string    `json:"instance_id"`
	EventType  string    `json:"event_type"`
	Detail     string    `json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

// Validate 规范化并校验事件字段。
func (e *RegistryEvent) Validate() error {
	e.ServiceID = strings.TrimSpace(e.ServiceID)
	e.InstanceID = strings.TrimSpace(e.InstanceID)
	e.Detail = strings.TrimSpace(e.Detail)
	if e.EventType == "" {
		return NewValidationError("event_type", "事件类型不能为空")
	}
	return nil
}

// EventFilter 事件筛选条件。
type EventFilter struct {
	ServiceID string
	EventType string
}

// Match 判断事件是否命中筛选条件。
func (f EventFilter) Match(e *RegistryEvent) bool {
	if f.ServiceID != "" && e.ServiceID != f.ServiceID {
		return false
	}
	if f.EventType != "" && e.EventType != f.EventType {
		return false
	}
	return true
}
