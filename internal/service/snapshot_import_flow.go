package service

import (
	"errors"

	"serviceregistry/internal/store"
)

type SnapshotImportFlow struct {
	state *store.SnapshotImportLeaseState
}

func NewSnapshotImportFlow(state *store.SnapshotImportLeaseState) *SnapshotImportFlow {
	return &SnapshotImportFlow{state: state}
}

func (f *SnapshotImportFlow) Process(key string, fail bool) error {
	token, ok := f.state.Acquire(key)
	if !ok {
		return errors.New("lease busy")
	}
	// Release on every exit path, including failure, so a retry of the
	// same snapshot is not told the lease is still busy.
	defer f.state.Release(key, token)
	if fail {
		return errors.New("operation failed")
	}
	return nil
}
