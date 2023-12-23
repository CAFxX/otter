package spinlock

import (
	"runtime"
	"sync/atomic"
)

const maxSpins = 16

// SpinLock is an implementation of spinlock.
type SpinLock struct {
	v atomic.Bool
}

// Lock locks sl. If the lock is already in use, the calling goroutine blocks until the spinlock is available.
func (sl *SpinLock) Lock() {
	spins := 0
	for {
		for sl.v.Load() {
			spins++
			if spins > maxSpins {
				spins = 0
				runtime.Gosched()
			} else {
				// TODO: add a PAUSE instruction
			}
		}

		if !sl.v.Swap(true) {
			return
		}
	}
}

// Unlock unlocks sl. A locked SpinLock is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a SpinLock and then arrange for another goroutine to unlock it.
func (sl *SpinLock) Unlock() {
	sl.v.Store(false)
}
