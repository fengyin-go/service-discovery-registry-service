package service

import "serviceregistry/internal/store"

type InstanceRegisterRetryFlow struct {
	state *store.InstanceRegisterRetryRetryState
}

func NewInstanceRegisterRetryFlow(state *store.InstanceRegisterRetryRetryState) *InstanceRegisterRetryFlow {
	return &InstanceRegisterRetryFlow{state: state}
}

func (f *InstanceRegisterRetryFlow) Execute() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		last = f.state.Next()
		if last == nil {
			return nil
		}
		// 仅临时故障允许重试；永久拒绝或未分类错误立即停止，
		// 避免永久拒绝被当成临时故障重复执行，以及后续的错误提交。
		if rf, ok := last.(*store.RetryFailure); !ok || !rf.Temporary {
			return last
		}
	}
	return last
}
