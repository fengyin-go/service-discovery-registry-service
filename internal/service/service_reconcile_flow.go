package service

import (
	"errors"

	"serviceregistry/internal/store"
)

type ServiceReconcileFlow struct {
	state *store.ServiceReconcileLeaseState
}

func NewServiceReconcileFlow(state *store.ServiceReconcileLeaseState) *ServiceReconcileFlow {
	return &ServiceReconcileFlow{state: state}
}

func (f *ServiceReconcileFlow) Process(key string, fail bool) error {
	token, ok := f.state.Acquire(key)
	if !ok {
		return errors.New("lease busy")
	}
	// 无论后续操作成功还是失败，本轮租约都必须释放，
	// 否则同一服务的重试会被上一轮遗留的忙碌状态拒绝。
	defer f.state.Release(key, token)
	if fail {
		return errors.New("operation failed")
	}
	return nil
}
