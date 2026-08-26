package service

import "serviceregistry/internal/store"

type EventPublishRetryFlow struct {
	state *store.EventPublishRetryRetryState
}

func NewEventPublishRetryFlow(state *store.EventPublishRetryRetryState) *EventPublishRetryFlow {
	return &EventPublishRetryFlow{state: state}
}

func (f *EventPublishRetryFlow) Execute() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		last = f.state.Next()
		if last == nil {
			return nil
		}
		// 永久拒绝不可重试：立即返回，绝不进入下一轮发布。
		// 只有临时错误才允许重试制造短暂繁忙。
		if !store.IsTemporaryError(last) {
			return last
		}
	}
	return last
}
