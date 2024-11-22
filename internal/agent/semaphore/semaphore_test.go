package semaphore

import (
	"testing"
)

func TestNewSemaphore(t *testing.T) {
	maxReq := 3
	sem := NewSemaphore(maxReq)
	if sem == nil {
		t.Fatal("NewSemaphore() returned nil, expected valid Semaphore instance")
	}
	if cap(sem.semaCh) != maxReq {
		t.Errorf("NewSemaphore() channel capacity = %d, want %d", cap(sem.semaCh), maxReq)
	}
}

func TestSemaphore_AcquireRelease(t *testing.T) {
	maxReq := 1
	sem := NewSemaphore(maxReq)

	// Acquire should succeed immediately when below max capacity
	sem.Acquire() // First acquire
	select {
	case sem.semaCh <- struct{}{}:
		t.Error("Acquire did not lock as expected at max capacity")
	default:
		// Expected behavior: channel is at capacity
	}

	// Release should free up a slot
	sem.Release()
	select {
	case sem.semaCh <- struct{}{}:
		// Expected behavior: slot freed by Release, no errors
	default:
		t.Error("Release did not free up the semaphore slot")
	}
}

func TestSemaphore_MaxConcurrentAcquires(t *testing.T) {
	maxReq := 2
	sem := NewSemaphore(maxReq)

	// Acquire up to maxReq times
	for i := 0; i < maxReq; i++ {
		sem.Acquire()
	}

	// Attempting another acquire should block (using non-blocking check)
	select {
	case sem.semaCh <- struct{}{}:
		t.Error("Acquire should have blocked, but it did not")
	default:
		// Expected behavior: channel at capacity
	}

	// Release should free up a slot, allowing another acquire
	sem.Release()
	select {
	case sem.semaCh <- struct{}{}:
		// Expected behavior: slot freed, acquire proceeds
	default:
		t.Error("Acquire did not proceed after Release was called")
	}
}
