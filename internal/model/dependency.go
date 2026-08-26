package model

import (
	"strings"
	"time"
)

// ServiceDependency 服务依赖关系。
type ServiceDependency struct {
	ID          string    `json:"id"`
	ServiceID   string    `json:"service_id"`
	DependsOnID string    `json:"depends_on_id"`
	Protocol    string    `json:"protocol"`
	CallCount   int64     `json:"call_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 规范化并校验依赖关系字段。
func (d *ServiceDependency) Validate() error {
	d.ServiceID = strings.TrimSpace(d.ServiceID)
	d.DependsOnID = strings.TrimSpace(d.DependsOnID)
	d.Protocol = strings.TrimSpace(d.Protocol)
	if d.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if d.DependsOnID == "" {
		return NewValidationError("depends_on_id", "被依赖服务 ID 不能为空")
	}
	if d.ServiceID == d.DependsOnID {
		return NewValidationError("depends_on_id", "服务不能依赖自身")
	}
	if d.Protocol == "" {
		d.Protocol = "http"
	}
	if d.CallCount < 0 {
		return NewValidationError("call_count", "调用次数不能为负数")
	}
	return nil
}
