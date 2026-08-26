package service

import (
	"errors"

	"serviceregistry/internal/store"
)

type DependencyRebuildFlow struct {
	state *store.DependencyRebuildLeaseState
}

func NewDependencyRebuildFlow(state *store.DependencyRebuildLeaseState) *DependencyRebuildFlow {
	return &DependencyRebuildFlow{state: state}
}

func (f *DependencyRebuildFlow) Process(key string, fail bool) error {
	token, ok := f.state.Acquire(key)
	if !ok {
		return errors.New("lease busy")
	}
	// 注册在所有后续返回路径（含 fail 分支）上释放租约，
	// 保证依赖重建出错返回前活动租约一定被清理。
	defer f.state.Release(key, token)
	if fail {
		return errors.New("operation failed")
	}
	return nil
}
