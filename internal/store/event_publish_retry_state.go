package store

type RetryFailure struct {
	Temporary bool
	Message   string
}

func (e *RetryFailure) Error() string { return e.Message }

type EventPublishRetryRetryState struct {
	steps    []error
	attempts int
	commits  int
}

func NewEventPublishRetryRetryState(steps ...error) *EventPublishRetryRetryState {
	return &EventPublishRetryRetryState{steps: append([]error(nil), steps...)}
}

func (s *EventPublishRetryRetryState) Next() error {
	s.attempts++
	var err error
	if len(s.steps) > 0 {
		err = s.steps[0]
		s.steps = s.steps[1:]
	}
	if err != nil {
		// 保留原始错误类型，使调用方能据 Temporary 区分永久错误与临时错误。
		// 此前用 errors.New(err.Error()) 重写会丢失 *RetryFailure 类型，
		// 导致永久拒绝被误当作可重试错误。
		return err
	}
	s.commits++
	return nil
}

func (s *EventPublishRetryRetryState) Attempts() int { return s.attempts }
func (s *EventPublishRetryRetryState) Commits() int  { return s.commits }

// IsTemporaryError 报告 err 是否为可重试的临时错误。
// 永久拒绝（Temporary=false）与未知错误均不可重试。
func IsTemporaryError(err error) bool {
	f, ok := err.(*RetryFailure)
	return ok && f.Temporary
}
