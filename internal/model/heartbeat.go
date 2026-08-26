package model

import (
	"strings"
	"time"
)

// 心跳结果常量。
const (
	HeartbeatOK      = "ok"
	HeartbeatTimeout = "timeout"
)

// Heartbeat 实例心跳记录。
type Heartbeat struct {
	ID         string    `json:"id"`
	InstanceID string    `json:"instance_id"`
	SentAt     time.Time `json:"sent_at"`
	ReceivedAt time.Time `json:"received_at"`
	LatencyMs  int64     `json:"latency_ms"`
	Result     string    `json:"result"`
}

// Validate 规范化并校验心跳记录字段。
func (h *Heartbeat) Validate() error {
	h.InstanceID = strings.TrimSpace(h.InstanceID)
	if h.InstanceID == "" {
		return NewValidationError("instance_id", "实例 ID 不能为空")
	}
	if h.Result == "" {
		h.Result = HeartbeatOK
	}
	if h.Result != HeartbeatOK && h.Result != HeartbeatTimeout {
		return NewValidationError("result", "心跳结果不合法")
	}
	if h.LatencyMs < 0 {
		return NewValidationError("latency_ms", "心跳延迟不能为负数")
	}
	return nil
}

// HeartbeatFilter 心跳记录筛选条件。
type HeartbeatFilter struct {
	InstanceID string
	Result     string
}

// Match 判断心跳记录是否命中筛选条件。
func (f HeartbeatFilter) Match(h *Heartbeat) bool {
	if f.InstanceID != "" && h.InstanceID != f.InstanceID {
		return false
	}
	if f.Result != "" && h.Result != f.Result {
		return false
	}
	return true
}
