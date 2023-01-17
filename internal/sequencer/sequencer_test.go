package sequencer_test

import (
	"context"
	"errors"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	sequencer.New(context.Background(), 0)
}

func TestSequencer_BeforeSequence(t *testing.T) {
	var (
		s         = sequencer.New(context.Background(), 1)
		wg        = sync.WaitGroup{}
		gotCalled bool
	)

	s.BeforeSequence(func(s *sequencer.Sequencer) {
		gotCalled = true
	})

	if gotCalled {
		t.Fatalf("expected before sequence not to be called before the sequence is enqueued, but it was")
	}

	wg.Add(1)
	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				if !gotCalled {
					t.Fatalf("expected sequence to be called after the sequence was enqueued, but it was not")
				}
				wg.Done()

				return nil
			}),
		}
	})

	wg.Wait()
}

func TestSequencer_AfterSequence(t *testing.T) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		s           = sequencer.New(ctx, 2)
		wg          = sync.WaitGroup{}
		gotCalled   int
	)
	defer cancel()
	wg.Add(1)

	s.AfterSequence(func(s *sequencer.Sequencer) {
		gotCalled++
		wg.Done()
	})

	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				gotCalled++

				return nil
			}),
		}
	})
	wg.Wait()

	if gotCalled == 0 {
		t.Fatalf("expected element to ber called, but it was not")
	}
}

func TestSequencer_HasSequence(t *testing.T) {
	var (
		s      = sequencer.New(context.Background(), 1)
		wg     = sync.WaitGroup{}
		freeze = sync.WaitGroup{}
	)

	if s.HasSequence() == true {
		t.Fatalf("expected HasSequence to return false, but it returned true")
	}

	freeze.Add(1)
	s.AfterSequence(func(s *sequencer.Sequencer) {
		freeze.Wait() // freeze sequencer in after callback to test on the hasSequence flag
	})

	wg.Add(1)
	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				wg.Done()
				return nil
			}),
		}
	})

	wg.Wait()

	if s.HasSequence() == false {
		t.Fatalf("expected HasSequence to return true, but it returned false")
	}

	freeze.Done()
}

func TestSequencer_Looping(t *testing.T) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		s           = sequencer.New(ctx, 1)
		wg          = sync.WaitGroup{}
		called      int
	)

	wg.Add(1)
	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				called++
				if called == 2 {
					wg.Done()
				}
				return nil
			}),
			sequencer.Loop{},
		}
	})

	wg.Wait()
	cancel()
}

func TestSequencer_IsPlaying(t *testing.T) {
	var (
		ctx, cancel      = context.WithCancel(context.Background())
		s                = sequencer.New(ctx, 1)
		waitForPLaying   = sync.WaitGroup{}
		waitForAssertion = sync.WaitGroup{}
	)

	waitForAssertion.Add(1)
	waitForPLaying.Add(1)
	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				waitForPLaying.Done()
				waitForAssertion.Wait()
				return nil
			}),
		}
	})

	waitForPLaying.Wait()

	var (
		want = true
		got  = s.IsPlaying()
	)
	if got != want {
		t.Errorf("want: %v, got: %v", want, got)
	}

	waitForAssertion.Done()
	cancel()
}

func TestSequencer_WaitElement(t *testing.T) {
	var (
		ctx, cancel  = context.WithCancel(context.Background())
		s            = sequencer.New(ctx, 1)
		start        time.Time
		playDuration time.Duration
		waitTime     = time.Millisecond * 30
		waitPlayEnd  = sync.WaitGroup{}
	)

	waitPlayEnd.Add(1)
	s.AfterSequence(func(s *sequencer.Sequencer) {
		playDuration = time.Since(start)
		waitPlayEnd.Done()
	})
	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				start = time.Now()
				return nil
			}),
			sequencer.Wait{
				Duration: waitTime,
			},
		}
	})

	waitPlayEnd.Wait()
	min := waitTime
	max := waitTime * 2 // roughly measured
	if playDuration < min || playDuration > max {
		t.Fatalf("expected play duration between %v and %v, but it was %v", min, max, playDuration)
	}
	cancel()
}

func TestSequencer_Pause(t *testing.T) {
	var (
		ctx, cancel  = context.WithCancel(context.Background())
		s            = sequencer.New(ctx, 1)
		firstElem    = make(chan struct{})
		secondElem   = make(chan struct{})
		thirdElem    = make(chan struct{})
		waitForPause = make(chan struct{})
	)

	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				close(firstElem)
				<-waitForPause
				return nil
			}),
			sequencer.ElemFunc(func() error {
				close(secondElem)
				return nil
			}),
			sequencer.ElemFunc(func() error {
				close(thirdElem)
				return nil
			}),
		}
	})

	<-firstElem
	err := s.Pause()
	if err != nil {
		t.Fatalf("did not expect Pause (1) to return an error, but got: %v", err)
	}
	// calling pause a second time will return an error, sequencer still is paused
	err = s.Pause()
	if !errors.Is(err, sequencer.ErrPauseAlreadyInProgress) {
		t.Fatalf("did expect Pause to return %v, but got: %v", sequencer.ErrPauseAlreadyInProgress, err)
	}
	close(waitForPause)

	var waitTime = time.NewTimer(time.Millisecond * 100)
	select {
	case <-waitTime.C:
	case <-secondElem:
		t.Fatal("sequence should be paused, but second element was played")
	}

	err = s.Pause()
	if err != nil {
		t.Fatalf("did not expect Pause (2) to return an error, but got: %v", err)
	}

	<-secondElem
	<-thirdElem

	cancel()
}
