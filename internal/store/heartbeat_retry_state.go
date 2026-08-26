package store

import "errors"

// RetryFailure 表示一次可重试或不可重试的写入失败。
// Temporary 为 true 表示临时繁忙，可重试；为 false 表示永久拒绝，不可重试。
type RetryFailure struct {
	Temporary bool
	Message   string
}

func (e *RetryFailure) Error() string { return e.Message }

// NewTemporaryFailure 构造一次临时繁忙失败（可重试）。
func NewTemporaryFailure(msg string) error {
	return &RetryFailure{Temporary: true, Message: msg}
}

// NewPermanentFailure 构造一次永久拒绝失败（不可重试）。
func NewPermanentFailure(msg string) error {
	return &RetryFailure{Temporary: false, Message: msg}
}

// IsRetryable 判断错误是否为可重试的临时繁忙失败。
// 非临时失败（永久拒绝）或非 RetryFailure 错误均不可重试。
func IsRetryable(err error) bool {
	var f *RetryFailure
	if !errors.As(err, &f) {
		return false
	}
	return f.Temporary
}

type HeartbeatRetryRetryState struct {
	steps    []error
	attempts int
	commits  int
}

func NewHeartbeatRetryRetryState(steps ...error) *HeartbeatRetryRetryState {
	return &HeartbeatRetryRetryState{steps: append([]error(nil), steps...)}
}

// Next 推进一次写入步骤：返回遇到的失败，返回 nil 表示本次写入成功并计入提交。
// 返回的错误保留原始类型（如 *RetryFailure），以便调用方区分临时繁忙与永久拒绝。
func (s *HeartbeatRetryRetryState) Next() error {
	s.attempts++
	var err error
	if len(s.steps) > 0 {
		err = s.steps[0]
		s.steps = s.steps[1:]
	}
	if err != nil {
		return err
	}
	s.commits++
	return nil
}

// HasNext 判断是否仍有可推进的写入步骤。
func (s *HeartbeatRetryRetryState) HasNext() bool { return len(s.steps) > 0 }

func (s *HeartbeatRetryRetryState) Attempts() int { return s.attempts }
func (s *HeartbeatRetryRetryState) Commits() int  { return s.commits }
