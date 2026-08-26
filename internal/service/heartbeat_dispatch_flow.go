package service

import (
	"context"

	"serviceregistry/internal/store"
)

type HeartbeatDispatchFlow struct{ state *store.HeartbeatDispatchState }

func NewHeartbeatDispatchFlow(state *store.HeartbeatDispatchState) *HeartbeatDispatchFlow {
	return &HeartbeatDispatchFlow{state: state}
}

func (f *HeartbeatDispatchFlow) Dispatch(ctx context.Context, key string) error {
	return f.state.Commit(context.Background(), key)
}
