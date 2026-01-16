package gplay

import (
	"context"
	"sync"
	"time"
)

type throttleState struct {
	mu          sync.Mutex
	startedAt   time.Time
	timesCalled int
	inThrottle  bool
}

func newThrottleState() *throttleState {
	return &throttleState{}
}

func (t *throttleState) wait(ctx context.Context, limitPerSecond int) error {
	if t == nil || limitPerSecond <= 0 {
		return nil
	}

	for {
		t.mu.Lock()
		if t.startedAt.IsZero() {
			t.startedAt = time.Now()
		}

		elapsed := time.Since(t.startedAt)
		if t.timesCalled < limitPerSecond && elapsed < time.Second {
			t.timesCalled++
			t.mu.Unlock()
			return nil
		}

		if !t.inThrottle {
			t.inThrottle = true
			t.mu.Unlock()

			timer := time.NewTimer(time.Second)
			select {
			case <-ctx.Done():
				timer.Stop()
				t.mu.Lock()
				t.inThrottle = false
				t.mu.Unlock()
				return ctx.Err()
			case <-timer.C:
			}

			t.mu.Lock()
			t.timesCalled = 0
			t.startedAt = time.Now()
			t.inThrottle = false
			t.mu.Unlock()

			continue
		}

		t.mu.Unlock()

		ticker := time.NewTicker(2 * time.Millisecond)
		select {
		case <-ctx.Done():
			ticker.Stop()
			return ctx.Err()
		case <-ticker.C:
			ticker.Stop()
			continue
		}
	}
}
