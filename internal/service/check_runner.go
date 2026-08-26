package service

import (
	"time"

	"serviceregistry/internal/model"
)

// 探测结论常量。
const (
	VerdictHealthy   = "healthy"
	VerdictUnhealthy = "unhealthy"
	VerdictUnknown   = "unknown"
)

// ProbeVerdict 实例健康探测结论。
type ProbeVerdict struct {
	InstanceID     string   `json:"instance_id"`
	Address        string   `json:"address"`
	Verdict        string   `json:"verdict"`
	HeartbeatAgeSec int64   `json:"heartbeat_age_seconds"`
	Reasons        []string `json:"reasons"`
}

// ProbeInstance 综合心跳、健康检查与实例状态给出探测结论。
func (s *Service) ProbeInstance(instanceID string) (ProbeVerdict, error) {
	i, err := s.store.GetInstance(instanceID)
	if err != nil {
		return ProbeVerdict{}, err
	}
	v := ProbeVerdict{
		InstanceID: i.ID,
		Address:    i.Address(),
		Verdict:    VerdictUnknown,
		Reasons:    []string{},
	}
	age := time.Since(i.LastHeartbeat)
	v.HeartbeatAgeSec = int64(age.Seconds())
	timeout := 30
	if s.cfg != nil && s.cfg.HeartbeatTimeoutSec > 0 {
		timeout = s.cfg.HeartbeatTimeoutSec
	}
	if age > time.Duration(timeout)*time.Second {
		v.Reasons = append(v.Reasons, "心跳超时")
		v.Verdict = VerdictUnhealthy
	} else {
		v.Reasons = append(v.Reasons, "心跳正常")
		v.Verdict = VerdictHealthy
	}
	if hc, hcErr := s.store.GetHealthCheckByInstance(instanceID); hcErr == nil {
		if hc.Status == model.CheckStatusUnhealthy {
			v.Reasons = append(v.Reasons, "健康检查不通过")
			v.Verdict = VerdictUnhealthy
		}
	}
	if i.Status == model.InstanceStatusDown {
		v.Reasons = append(v.Reasons, "实例已下线")
		v.Verdict = VerdictUnhealthy
	}
	return v, nil
}

// ProbeAll 探测全部实例。
func (s *Service) ProbeAll() []ProbeVerdict {
	out := make([]ProbeVerdict, 0)
	for _, i := range s.store.ListInstances() {
		if v, err := s.ProbeInstance(i.ID); err == nil {
			out = append(out, v)
		}
	}
	return out
}
