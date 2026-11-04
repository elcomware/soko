package middleware

import (
	"context"
	"testing"
	"time"
)

func TestSessionLockWaitCancellation(t *testing.T) {
	release, err := acquireSession(context.Background(), "same-session")
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		unlock, err := acquireSession(ctx, "same-session")
		if unlock != nil {
			unlock()
		}
		result <- err
	}()
	cancel()
	select {
	case err := <-result:
		if err != context.Canceled {
			t.Fatalf("wanted cancellation, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("cancelled waiter remained blocked")
	}
}
func TestSessionLockDeadlineAndRelease(t *testing.T) {
	release, err := acquireSession(context.Background(), "deadline-session")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	unlock, err := acquireSession(ctx, "deadline-session")
	if unlock != nil {
		unlock()
	}
	release()
	if err != context.DeadlineExceeded {
		t.Fatalf("wanted deadline, got %v", err)
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()
	unlock, err = acquireSession(ctx2, "deadline-session")
	if err != nil {
		t.Fatal(err)
	}
	unlock()
}
