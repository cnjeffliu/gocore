/*
 * @Author: Jeffrey Liu
 * @Date: 2022-08-26 11:42:04
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-12-15 14:38:53
 * @Description:
 */
package timex

import (
	"math/rand"
	"time"
)

// Jitter return duration + rand * maxFactor * duration
// For example for 10s and jitter 1, it will return a time within [10s, 20s])
func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	wait := duration + time.Duration(rand.Float64()*maxFactor*float64(duration))
	return wait
}

// jitterUp return duration which added a random jitter
//
// This adds or subtracts time from the duration within a given jitter fraction.
// For example for 10s and jitter 0.1, it will return a time within [9s, 11s])
func jitterUp(duration time.Duration, jitter float64) time.Duration {
	multiplier := jitter * (rand.Float64()*2 - 1)
	return time.Duration(float64(duration) * (1 + multiplier))
}

type BackoffManager interface {
	Backoff() Timer
}

type jitteredBackoffManagerImpl struct {
	clock    Clock
	timer    Timer
	duration time.Duration
	jitter   float64
}

// NewJitteredBackoffManager returns a BackoffManager that backoffs with given duration plus given jitter. If the jitter
// is negative, backoff will not be jittered.
func NewJitteredBackoffManager(duration time.Duration, jitter float64, c Clock) BackoffManager {
	return &jitteredBackoffManagerImpl{
		clock:    c,
		duration: duration,
		jitter:   jitter,
		timer:    nil,
	}
}

func (j *jitteredBackoffManagerImpl) getNextBackoff() time.Duration {
	jitteredPeriod := j.duration
	if j.jitter > 0.0 {
		jitteredPeriod = Jitter(j.duration, j.jitter)
	}
	return jitteredPeriod
}

// Backoff implements BackoffManager.Backoff, it returns a timer so caller can block on the timer for jittered backoff.
// The returned timer must be drained before calling Backoff() the second time
func (j *jitteredBackoffManagerImpl) Backoff() Timer {
	backoff := j.getNextBackoff()
	if j.timer == nil {
		j.timer = j.clock.NewTimer(backoff)
	} else {
		j.timer.Reset(backoff)
	}
	return j.timer
}

func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
	JitterUntil(f, period, 0.0, true, stopCh)
}

func JitterUntil(f func(), period time.Duration, jitterFactor float64, sliding bool, stopCh <-chan struct{}) {
	BackoffUntil(f, NewJitteredBackoffManager(period, jitterFactor, &RealClock{}), sliding, stopCh)
}

func BackoffUntil(f func(), backoff BackoffManager, sliding bool, stopCh <-chan struct{}) {
	var t Timer
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		if !sliding {
			t = backoff.Backoff()
		}

		f()

		if sliding {
			t = backoff.Backoff()
		}

		// NOTE: b/c there is no priority selection in golang
		// it is possible for this to race, meaning we could
		// trigger t.C and stopCh, and t.C select falls through.
		// In order to mitigate we re-check stopCh at the beginning
		// of every loop to prevent extra executions of f().
		select {
		case <-stopCh:
			if !t.Stop() {
				<-t.Chan()
			}
			return
		case <-t.Chan():
		}
	}
}
