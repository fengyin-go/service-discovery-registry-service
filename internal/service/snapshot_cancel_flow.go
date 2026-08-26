package service

import (
	"context"

	"serviceregistry/internal/store"
)

type SnapshotCancelFlow struct{ state *store.SnapshotCancelState }

func NewSnapshotCancelFlow(state *store.SnapshotCancelState) *SnapshotCancelFlow {
	return &SnapshotCancelFlow{state: state}
}

func (f *SnapshotCancelFlow) Dispatch(ctx context.Context, key string) error {
	return f.state.Commit(context.Background(), key)
}
