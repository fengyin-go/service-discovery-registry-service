package service

import "serviceregistry/internal/store"

type ConfigUpdateRetryFlow struct {
	state *store.ConfigUpdateRetryRetryState
}

func NewConfigUpdateRetryFlow(state *store.ConfigUpdateRetryRetryState) *ConfigUpdateRetryFlow {
	return &ConfigUpdateRetryFlow{state: state}
}

// Execute 推进配置更新：仅当失败为临时占用（Temporary）时才重试，
// 永久拒绝（不可恢复）原样返回，不重试、不提交。
func (f *ConfigUpdateRetryFlow) Execute() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		last = f.state.Next()
		if last == nil {
			return nil
		}
		// 永久拒绝不可恢复，原样返回，不再重试。
		if failure := store.AsRetryFailure(last); failure != nil && !failure.IsTemporary() {
			return last
		}
		// 临时占用：继续下一轮重试（最后一轮后直接返回 last）。
	}
	return last
}
