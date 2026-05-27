package sequencer_test

import (
	"context"
	"errors"
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
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

func TestSequencer_BeforeSequenceCalledForQueuedSequence(t *testing.T) {
	var (
		ctx, cancel          = context.WithCancel(context.Background())
		s                    = sequencer.New(ctx, 1)
		beforeSequenceCalled = make(chan struct{}, 2)
		firstDone            = make(chan struct{})
		secondStarted        = make(chan struct{})
	)

	defer cancel()

	s.BeforeSequence(func(s *sequencer.Sequencer) {
		beforeSequenceCalled <- struct{}{}
	})

	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				<-firstDone
				return nil
			}),
		}
	})

	<-beforeSequenceCalled

	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				close(secondStarted)
				return nil
			}),
		}
	})

	close(firstDone)

	select {
	case <-beforeSequenceCalled:
	case <-time.After(time.Second):
		t.Fatal("expected before sequence callback for queued sequence")
	}

	select {
	case <-secondStarted:
	case <-time.After(time.Second):
		t.Fatal("expected queued sequence to start")
	}
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

func TestSequencer_ElementError(t *testing.T) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		s           = sequencer.New(ctx, 1)
		wantErr     = errors.New("boom")
		gotErr      error
		done        = make(chan struct{})
	)

	defer cancel()

	s.ElementError(func(s *sequencer.Sequencer, elem sequencer.Elem, err error) {
		gotErr = err
		close(done)
	})

	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				return wantErr
			}),
		}
	})

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected element error callback")
	}

	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, gotErr)
	}
}

func TestSequencer_BeforeElement(t *testing.T) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		s           = sequencer.New(ctx, 1)
		events      = make(chan sequencer.ElementEvent, 2)
	)

	defer cancel()

	s.BeforeElement(func(s *sequencer.Sequencer, elem sequencer.Elem, event sequencer.ElementEvent) {
		events <- event
	})

	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.NoOperation{},
			sequencer.Wait{Duration: 10 * time.Millisecond},
		}
	})

	first := <-events
	if first.Label != "NOP" {
		t.Fatalf("expected first event label NOP, got %q", first.Label)
	}

	second := <-events
	if second.Label != "W 10ms" {
		t.Fatalf("expected second event label W 10ms, got %q", second.Label)
	}
	if second.DurationMs != 10 {
		t.Fatalf("expected second event duration 10ms, got %d", second.DurationMs)
	}

	third := sequencer.DescribeElement(general.KeyDown{Key: key.VkW})
	if third.Label != "KD W" {
		t.Fatalf("expected third event label KD W, got %q", third.Label)
	}
}

func TestSequencer_HasSequence(t *testing.T) {
	var (
		s                = sequencer.New(context.Background(), 1)
		sequenceStarted  = make(chan struct{})
		waitForAssertion = make(chan struct{})
		sequenceFinished = make(chan struct{})
	)

	if s.HasSequence() == true {
		t.Fatalf("expected HasSequence to return false, but it returned true")
	}

	s.AfterSequence(func(s *sequencer.Sequencer) {
		close(sequenceFinished)
	})

	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				close(sequenceStarted)
				<-waitForAssertion
				return nil
			}),
		}
	})

	<-sequenceStarted

	if s.HasSequence() == false {
		t.Fatalf("expected HasSequence to return true, but it returned false")
	}

	close(waitForAssertion)
	<-sequenceFinished

	if s.HasSequence() == true {
		t.Fatalf("expected HasSequence to return false after the sequence finished, but it returned true")
	}
}

func TestSequencer_FiniteSequenceReturnsToIdle(t *testing.T) {
	var (
		ctx, cancel      = context.WithCancel(context.Background())
		s                = sequencer.New(ctx, 1)
		sequenceFinished = make(chan struct{})
	)

	defer cancel()

	s.AfterSequence(func(s *sequencer.Sequencer) {
		close(sequenceFinished)
	})

	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				return nil
			}),
		}
	})

	select {
	case <-sequenceFinished:
	case <-time.After(time.Second):
		t.Fatal("expected finite sequence to finish")
	}

	if s.State() != sequencer.Idle {
		t.Fatalf("expected finite sequence to return to idle, got %v", s.State())
	}

	if s.HasSequence() {
		t.Fatal("expected finite sequence to clear HasSequence")
	}
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

func TestSequencer_AbortClearsHasSequence(t *testing.T) {
	var (
		ctx, cancel      = context.WithCancel(context.Background())
		s                = sequencer.New(ctx, 1)
		sequenceStarted  = make(chan struct{})
		waitForAssertion = make(chan struct{})
	)

	defer cancel()

	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				close(sequenceStarted)
				<-waitForAssertion
				return nil
			}),
		}
	})

	select {
	case <-sequenceStarted:
	case <-time.After(time.Second):
		t.Fatal("expected sequence to start")
	}

	s.Abort()

	if s.State() != sequencer.Idle {
		t.Fatalf("expected state to be idle after abort, got %v", s.State())
	}

	if s.HasSequence() {
		t.Fatal("expected abort to clear HasSequence")
	}

	close(waitForAssertion)
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

	minWait := waitTime
	maxWait := waitTime * 2 // roughly measured

	if playDuration < minWait ||
		playDuration > maxWait {
		t.Fatalf("expected play duration between %v and %v, but it was %v", minWait, maxWait, playDuration)
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

func TestSequencer_AbortWhilePausedDoesNotBlockFutureSequences(t *testing.T) {
	var (
		ctx, cancel       = context.WithCancel(context.Background())
		s                 = sequencer.New(ctx, 1)
		firstElemStarted  = make(chan struct{})
		firstElemRelease  = make(chan struct{})
		secondElemStarted = make(chan struct{})
		futureElemStarted = make(chan struct{})
	)

	defer cancel()

	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				close(firstElemStarted)
				<-firstElemRelease
				return nil
			}),
			sequencer.ElemFunc(func() error {
				close(secondElemStarted)
				return nil
			}),
		}
	})

	<-firstElemStarted

	if err := s.Pause(); err != nil {
		t.Fatalf("did not expect Pause to return an error, but got: %v", err)
	}

	close(firstElemRelease)
	time.Sleep(time.Millisecond * 20)

	s.Abort()

	if s.State() != sequencer.Idle {
		t.Fatalf("expected state to be idle after abort, but got %v", s.State())
	}

	select {
	case <-secondElemStarted:
		t.Fatal("second element should not run after abort")
	default:
	}

	s.EnqueueSequence(func() []sequencer.Elem {
		return []sequencer.Elem{
			sequencer.ElemFunc(func() error {
				close(futureElemStarted)
				return nil
			}),
		}
	})

	select {
	case <-futureElemStarted:
	case <-time.After(time.Second):
		t.Fatal("expected future sequence to run after aborting pause")
	}
}
