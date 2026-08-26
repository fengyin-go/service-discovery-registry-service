package model

import (
	"strings"
	"time"
)

// 健康检查类型常量。
const (
	CheckTypeTCP = "tcp"
	CheckTypeHTTP = "http"
)

// 健康检查状态常量。
const (
	CheckStatusHealthy   = "healthy"
	CheckStatusUnhealthy = "unhealthy"
)

// HealthCheck 实例健康检查配置与结果。
type HealthCheck struct {
	ID            string    `json:"id"`
	InstanceID    string    `json:"instance_id"`
	Type          string    `json:"type"`
	Target        string    `json:"target"`
	IntervalSec   int       `json:"interval_sec"`
	TimeoutSec    int       `json:"timeout_sec"`
	Threshold     int       `json:"threshold"`
	Status        string    `json:"status"`
	LastCheckedAt time.Time `json:"last_checked_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// Validate 规范化并校验健康检查字段。
func (h *HealthCheck) Validate() error {
	h.InstanceID = strings.TrimSpace(h.InstanceID)
	h.Target = strings.TrimSpace(h.Target)
	if h.InstanceID == "" {
		return NewValidationError("instance_id", "实例 ID 不能为空")
	}
	if h.Type == "" {
		h.Type = CheckTypeTCP
	}
	if h.Type != CheckTypeTCP && h.Type != CheckTypeHTTP {
		return NewValidationError("type", "健康检查类型不合法")
	}
	if h.Target == "" {
		return NewValidationError("target", "健康检查目标不能为空")
	}
	if h.IntervalSec <= 0 {
		h.IntervalSec = 10
	}
	if h.TimeoutSec <= 0 {
		h.TimeoutSec = 3
	}
	if h.Threshold <= 0 {
		h.Threshold = 3
	}
	if h.Status == "" {
		h.Status = CheckStatusHealthy
	}
	if h.Status != CheckStatusHealthy && h.Status != CheckStatusUnhealthy {
		return NewValidationError("status", "健康检查状态不合法")
	}
	return nil
}
