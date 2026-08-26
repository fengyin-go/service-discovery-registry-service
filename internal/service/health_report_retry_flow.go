package service

import (
	"errors"

	"serviceregistry/internal/store"
)

type HealthReportRetryFlow struct {
	state *store.HealthReportRetryRetryState
}

func NewHealthReportRetryFlow(state *store.HealthReportRetryRetryState) *HealthReportRetryFlow {
	return &HealthReportRetryFlow{state: state}
}

func (f *HealthReportRetryFlow) Execute() error {
	err := f.state.Next()
	if err == nil {
		return nil
	}
	// 仅临时错误才再试一次；不可重试的拒绝直接返回，保留原始错误类别。
	var rf *store.RetryFailure
	if errors.As(err, &rf) && rf.Temporary {
		return f.state.Next()
	}
	return err
}
