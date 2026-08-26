package store

type RetryFailure struct {
	Temporary bool
	Message   string
}

func (e *RetryFailure) Error() string { return e.Message }

type HealthReportRetryRetryState struct {
	steps    []error
	attempts int
	commits  int
}

func NewHealthReportRetryRetryState(steps ...error) *HealthReportRetryRetryState {
	return &HealthReportRetryRetryState{steps: append([]error(nil), steps...)}
}

func (s *HealthReportRetryRetryState) Next() error {
	s.attempts++
	var err error
	if len(s.steps) > 0 {
		err = s.steps[0]
		s.steps = s.steps[1:]
	}
	if err != nil {
		// 保留原始错误，避免丢失错误类别（如 RetryFailure.Temporary）。
		return err
	}
	s.commits++
	return nil
}

func (s *HealthReportRetryRetryState) Attempts() int { return s.attempts }
func (s *HealthReportRetryRetryState) Commits() int  { return s.commits }
