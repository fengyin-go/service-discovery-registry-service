package service

import (
	"errors"

	"serviceregistry/internal/store"
)

type ProbeLeaseFlow struct{ state *store.ProbeLeaseLeaseState }

func NewProbeLeaseFlow(state *store.ProbeLeaseLeaseState) *ProbeLeaseFlow {
	return &ProbeLeaseFlow{state: state}
}

func (f *ProbeLeaseFlow) Process(key string, fail bool) error {
	token, ok := f.state.Acquire(key)
	if !ok {
		return errors.New("lease busy")
	}
	// 先注册释放，确保失败分支也会归还租约，否则同一实例再探测会一直被占用。
	defer f.state.Release(key, token)
	if fail {
		return errors.New("operation failed")
	}
	return nil
}
