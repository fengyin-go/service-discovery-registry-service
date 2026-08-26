package service

import "serviceregistry/internal/store"

type HeartbeatRetryFlow struct {
	state *store.HeartbeatRetryRetryState
}

func NewHeartbeatRetryFlow(state *store.HeartbeatRetryRetryState) *HeartbeatRetryFlow {
	return &HeartbeatRetryFlow{state: state}
}

// Execute 推进心跳写入直至成功、耗尽步骤或遇到不可重试的永久拒绝。
// 仅重试临时繁忙（store.IsRetryable 为 true）；永久拒绝立即返回且不再提交。
func (f *HeartbeatRetryFlow) Execute() error {
	for {
		err := f.state.Next()
		if err == nil {
			return nil
		}
		// 永久拒绝（或非临时失败）不可重试，立即停止。
		if !store.IsRetryable(err) {
			return err
		}
		// 临时繁忙：若已无后续步骤可推进则停止，否则继续重试。
		if !f.state.HasNext() {
			return err
		}
	}
}
