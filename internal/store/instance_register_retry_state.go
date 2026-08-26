package store

import "errors"

type RetryFailure struct {
	Temporary bool
	Message   string
}

func (e *RetryFailure) Error() string { return e.Message }

type InstanceRegisterRetryRetryState struct {
	steps    []error
	attempts int
	commits  int
}

func NewInstanceRegisterRetryRetryState(steps ...error) *InstanceRegisterRetryRetryState {
	return &InstanceRegisterRetryRetryState{steps: append([]error(nil), steps...)}
}

// Next 模拟实例注册的下一次尝试结果。
// nil 表示注册成功并计入一次提交；否则返回错误。
// 返回 *RetryFailure 时保留其 Temporary 分类，便于调用方区分永久拒绝与临时故障，
// 避免把永久拒绝误当成可重试的临时故障。
func (s *InstanceRegisterRetryRetryState) Next() error {
	s.attempts++
	var err error
	if len(s.steps) > 0 {
		err = s.steps[0]
		s.steps = s.steps[1:]
	}
	if err == nil {
		s.commits++
		return nil
	}
	if rf, ok := err.(*RetryFailure); ok {
		return &RetryFailure{Temporary: rf.Temporary, Message: rf.Message}
	}
	return errors.New(err.Error())
}

func (s *InstanceRegisterRetryRetryState) Attempts() int { return s.attempts }
func (s *InstanceRegisterRetryRetryState) Commits() int  { return s.commits }
