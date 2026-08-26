package store

import "errors"

// RetryFailure 表示一次配置更新尝试的失败，临时失败可重试，永久失败不可恢复。
type RetryFailure struct {
	Temporary bool
	Message   string
}

func (e *RetryFailure) Error() string { return e.Message }

// IsTemporary 返回失败是否为可重试的临时占用。
func (e *RetryFailure) IsTemporary() bool { return e.Temporary }

// AsRetryFailure 从 err 中提取 *RetryFailure，未匹配时返回 nil。
func AsRetryFailure(err error) *RetryFailure {
	var f *RetryFailure
	if errors.As(err, &f) {
		return f
	}
	return nil
}

type ConfigUpdateRetryRetryState struct {
	steps    []error
	attempts int
	commits  int
}

func NewConfigUpdateRetryRetryState(steps ...error) *ConfigUpdateRetryRetryState {
	return &ConfigUpdateRetryRetryState{steps: append([]error(nil), steps...)}
}

// Next 取出并执行下一步：若有失败则原样返回该错误（保留类型，便于判定是否可重试），
// 否则记一次成功提交并返回 nil。
func (s *ConfigUpdateRetryRetryState) Next() error {
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

func (s *ConfigUpdateRetryRetryState) Attempts() int { return s.attempts }
func (s *ConfigUpdateRetryRetryState) Commits() int  { return s.commits }
