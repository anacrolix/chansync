package chansync

import "testing"

func TestBroadcastCond(t *testing.T) {
	var bc BroadcastCond
	ch := bc.Signaled()
	if bc.Signaled() != ch {
		t.Fatal("Signaled should be cached until Broadcast")
	}
	bc.Broadcast()
	select {
	case <-ch:
	default:
		t.Fatal("Broadcast should have closed the signalled channel")
	}
	if bc.Signaled() == ch {
		t.Fatal("Signaled should return a fresh channel after Broadcast")
	}
	// Broadcast with no channel allocated must be a no-op.
	var empty BroadcastCond
	empty.Broadcast()
}

// Benchmarks the common cycle of obtaining the signal channel and then broadcasting on it.
func BenchmarkBroadcastCondSignaledBroadcast(b *testing.B) {
	var bc BroadcastCond
	for b.Loop() {
		_ = bc.Signaled()
		bc.Broadcast()
	}
}

// Benchmarks repeated Signaled calls without any intervening Broadcast, exercising the cached
// channel path.
func BenchmarkBroadcastCondSignaled(b *testing.B) {
	var bc BroadcastCond
	for b.Loop() {
		_ = bc.Signaled()
	}
}

// Benchmarks Broadcast with no waiters and no allocated channel, exercising the nil fast path.
func BenchmarkBroadcastCondBroadcast(b *testing.B) {
	var bc BroadcastCond
	for b.Loop() {
		bc.Broadcast()
	}
}

// Benchmarks many goroutines reading the cached signal channel concurrently, the workload the type
// is built for.
func BenchmarkBroadcastCondSignaledParallel(b *testing.B) {
	var bc BroadcastCond
	bc.Signaled()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = bc.Signaled()
		}
	})
}
