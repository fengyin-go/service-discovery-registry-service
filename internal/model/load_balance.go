package model

import (
	"strings"
	"time"
)

// 负载均衡算法常量。
const (
	AlgorithmRoundRobin = "round_robin"
	AlgorithmRandom     = "random"
	AlgorithmWeighted   = "weighted"
	AlgorithmLeastConn  = "least_conn"
)

// LoadBalanceStrategy 服务级负载均衡策略。
type LoadBalanceStrategy struct {
	ID        string    `json:"id"`
	ServiceID string    `json:"service_id"`
	Algorithm string    `json:"algorithm"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 规范化并校验负载均衡策略字段。
func (l *LoadBalanceStrategy) Validate() error {
	l.ServiceID = strings.TrimSpace(l.ServiceID)
	if l.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if l.Algorithm == "" {
		l.Algorithm = AlgorithmRoundRobin
	}
	switch l.Algorithm {
	case AlgorithmRoundRobin, AlgorithmRandom, AlgorithmWeighted, AlgorithmLeastConn:
	default:
		return NewValidationError("algorithm", "负载均衡算法不合法")
	}
	return nil
}
