package kcron

import (
	"context"
	"sync"
	"time"
)

type Lock interface {
	Lock(ctx context.Context, ttl time.Duration) error
	Unlock(ctx context.Context) error
	Refresh(ctx context.Context, ttl time.Duration) error
}

type mockLock struct {
	mtx sync.Mutex

	optMtx sync.Mutex
	locked bool
}

func (m *mockLock) Lock(ctx context.Context, ttl time.Duration) error {
	m.optMtx.Lock()
	defer m.optMtx.Unlock()

	m.mtx.Lock()
	m.locked = true
	return nil
}

func (m *mockLock) Unlock(ctx context.Context) error {
	m.optMtx.Lock()
	defer m.optMtx.Unlock()

	if m.locked {
		m.mtx.Unlock()
		m.locked = false
	}
	return nil
}

func (m *mockLock) Refresh(ctx context.Context, ttl time.Duration) error {
	return nil
}
