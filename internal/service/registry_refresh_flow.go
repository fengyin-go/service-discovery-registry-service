package service

import (
	"errors"

	"serviceregistry/internal/store"
)

type RegistryRefreshFlow struct {
	state *store.RegistryRefreshLeaseState
}

func NewRegistryRefreshFlow(state *store.RegistryRefreshLeaseState) *RegistryRefreshFlow {
	return &RegistryRefreshFlow{state: state}
}

func (f *RegistryRefreshFlow) Process(key string, fail bool) error {
	token, ok := f.state.Acquire(key)
	if !ok {
		return errors.New("lease busy")
	}
	// 失败也必须释放占用标记，否则残留租约会挡住同名重启刷新。
	defer f.state.Release(key, token)
	if fail {
		return errors.New("operation failed")
	}
	return nil
}
