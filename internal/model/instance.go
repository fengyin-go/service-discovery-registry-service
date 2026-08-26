package model

import (
	"strconv"
	"strings"
	"time"
)

// 实例状态常量。
const (
	InstanceStatusStarting = "starting"
	InstanceStatusUp       = "up"
	InstanceStatusDown     = "down"
)

// instanceTransitions 定义实例合法状态流转。
var instanceTransitions = map[string]map[string]bool{
	InstanceStatusStarting: {InstanceStatusUp: true, InstanceStatusDown: true},
	InstanceStatusUp:       {InstanceStatusDown: true},
	InstanceStatusDown:     {InstanceStatusUp: true},
}

// CanInstanceTransition 判断实例能否从 from 流转到 to。
func CanInstanceTransition(from, to string) bool {
	if m, ok := instanceTransitions[from]; ok {
		return m[to]
	}
	return false
}

// Instance 服务实例实体。
type Instance struct {
	ID            string            `json:"id"`
	ServiceID     string            `json:"service_id"`
	Host          string            `json:"host"`
	Port          int               `json:"port"`
	Weight        int               `json:"weight"`
	Status        string            `json:"status"`
	Region        string            `json:"region"`
	Zone          string            `json:"zone"`
	Metadata      map[string]string `json:"metadata"`
	LastHeartbeat time.Time         `json:"last_heartbeat"`
	CreatedAt     time.Time         `json:"created_at"`
}

// Validate 规范化并校验实例字段。
func (i *Instance) Validate() error {
	i.Host = strings.TrimSpace(i.Host)
	i.ServiceID = strings.TrimSpace(i.ServiceID)
	i.Region = strings.TrimSpace(i.Region)
	i.Zone = strings.TrimSpace(i.Zone)
	if i.ServiceID == "" {
		return NewValidationError("service_id", "所属服务 ID 不能为空")
	}
	if i.Host == "" {
		return NewValidationError("host", "实例主机不能为空")
	}
	if i.Port <= 0 || i.Port > 65535 {
		return NewValidationError("port", "实例端口超出合法范围")
	}
	if i.Weight <= 0 {
		i.Weight = 1
	}
	if i.Status == "" {
		i.Status = InstanceStatusStarting
	}
	if i.Status != InstanceStatusStarting && i.Status != InstanceStatusUp && i.Status != InstanceStatusDown {
		return NewValidationError("status", "实例状态不合法")
	}
	return nil
}

// Address 返回实例的 host:port 地址。
func (i *Instance) Address() string {
	return i.Host + ":" + strconv.Itoa(i.Port)
}

// InstanceFilter 实例列表筛选条件。
type InstanceFilter struct {
	ServiceID string
	Status    string
	Region    string
	Zone      string
}

// Match 判断实例是否命中筛选条件。
func (f InstanceFilter) Match(i *Instance) bool {
	if f.ServiceID != "" && i.ServiceID != f.ServiceID {
		return false
	}
	if f.Status != "" && i.Status != f.Status {
		return false
	}
	if f.Region != "" && i.Region != f.Region {
		return false
	}
	if f.Zone != "" && i.Zone != f.Zone {
		return false
	}
	return true
}
