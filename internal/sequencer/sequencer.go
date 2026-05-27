package sequencer

import (
	"context"
	"errors"
	"fmt"
	"github.com/Oppodelldog/roamer/internal/logger"
	"time"
)

var ErrPauseAlreadyInProgress = errors.New("pause request is already in progress")

type State int

const (
	Idle State = iota
	Playing
	Paused
)

func (s State) String() string {
	switch s {
	case Playing:
		return "playing"
	case Paused:
		return "paused"
	default:
		return "idle"
	}
}

// Sequencer plays sequences of Elem by calling their Do method.
// It runs as a detached worker using go routines to playback sequences.
type Sequencer struct {
	sequence       chan func() []Elem
	seq            chan Elem
	pause          chan struct{}
	state          State
	hasSequence    bool
	beforeSequence func(s *Sequencer)
	afterSequence  func(s *Sequencer)
	beforeElement  func(s *Sequencer, elem Elem, event ElementEvent)
	elementError   func(s *Sequencer, elem Elem, err error)
	debug          bool
	seqWait        chan struct{}
}

// New creates a new sequencer.
// Parameter queueSize defines how many Sequences may be queued before dropping.
func New(ctx context.Context, queueSize int) *Sequencer {
	var s = &Sequencer{
		sequence: make(chan func() []Elem, queueSize),
		seq:      make(chan Elem),
		pause:    make(chan struct{}, 1),
		state:    Idle,
		debug:    false,
	}

	go s.playSequence(ctx)

	go s.playElement(ctx)

	return s
}

// BeforeSequence registers a function that is called before a sequence has finished.
func (s *Sequencer) BeforeSequence(f func(s *Sequencer)) {
	s.beforeSequence = f
}

// AfterSequence registers a function that is called after a sequence has finished.
func (s *Sequencer) AfterSequence(f func(s *Sequencer)) {
	s.afterSequence = f
}

func (s *Sequencer) ElementError(f func(s *Sequencer, elem Elem, err error)) {
	s.elementError = f
}

func (s *Sequencer) BeforeElement(f func(s *Sequencer, elem Elem, event ElementEvent)) {
	s.beforeElement = f
}

// IsPlaying indicated if playback is running or not.
func (s *Sequencer) IsPlaying() bool {
	return s.state == Playing
}

func (s *Sequencer) State() State {
	return s.state
}

// HasSequence indicates if sequencer has a sequence for playback or not.
func (s *Sequencer) HasSequence() bool {
	return s.hasSequence
}

// EnqueueSequence adds a new sequence to the sequence queue.
// This method blocks by the queueSize defined in New.
func (s *Sequencer) EnqueueSequence(newSeq func() []Elem) {
	var newSequence = newSeq
	select {
	case s.sequence <- newSequence:
		logger.Println("sequence enqueued")
	default:
		logger.Println("dropped waiting sequence")
		<-s.sequence
		s.sequence <- newSequence
	}
}

// Pause toggles playback (pause, resume).
func (s *Sequencer) Pause() error {
	select {
	case s.pause <- struct{}{}:
		if s.state == Paused {
			s.state = Playing
		} else if s.state == Playing {
			s.state = Paused
		}

		return nil
	default:
		return ErrPauseAlreadyInProgress
	}
}

func (s *Sequencer) waitForResume(ctx context.Context) {
	if s.state == Idle {
		logger.Println("pause ignored after abort")
		return
	}

	logger.Println("pausing")

	s.state = Paused

	select {
	case <-ctx.Done():
		logger.Println("stopped waiting for resume")
		return
	case <-s.pause:
	}

	if s.state == Idle {
		logger.Println("resume ignored after abort")
		return
	}

	logger.Println("resuming")
}

func (s *Sequencer) playSequence(ctx context.Context) {
	var newSequence func() []Elem
waitForNext:
	s.hasSequence = false

	select {
	case <-ctx.Done():
		logger.Println("stopped sequence queue")
		return
	case newSequence = <-s.sequence:
	}

	s.hasSequence = true

loop:
	s.state = Playing
	if s.beforeSequence != nil {
		s.beforeSequence(s)
	}

	newSeq := newSequence()
	fmt.Printf("play sequence length: %v\n", len(newSeq))

	s.seqWait = make(chan struct{})
	for _, e := range newSeq {
		s.seq <- e
		<-s.seqWait

		if s.state == Idle {
			logger.Println("is not playing sequence")
			goto waitForNext
		}
	}

	close(s.seqWait)

	select {
	case <-ctx.Done():
		logger.Println("stopped sequence queue")
		return
	case newSequence = <-s.sequence:
		logger.Println("got a new sequence")
		if s.afterSequence != nil {
			s.afterSequence(s)
		}

		goto loop
	default:
		logger.Println("got no new sequence")

		var newSeq = newSequence()
		if len(newSeq) > 0 {
			if _, isLoop := newSeq[len(newSeq)-1].(Loop); isLoop {
				logger.Println("looping sequence")
				if s.afterSequence != nil {
					s.afterSequence(s)
				}
				goto loop
			}
		}
	}

	s.state = Idle
	s.hasSequence = false

	if s.afterSequence != nil {
		s.afterSequence(s)
	}

	logger.Println("wait for next sequence")

	goto waitForNext
}

func (s *Sequencer) playElement(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Println("stopped element playback")
			return
		case <-s.pause:
			s.waitForResume(ctx)
			continue
		default:
		}

		select {
		case <-ctx.Done():
			logger.Println("stopped element playback")
			return
		case <-s.pause:
			s.waitForResume(ctx)
		case el := <-s.seq:
			if s.beforeElement != nil {
				s.beforeElement(s, el, DescribeElement(el))
			}

			switch v := el.(type) {
			case Wait:
				logger.Println("wait ", v.Duration)
				s.sleep(ctx, v.Duration)
			case Loop:
			case NoOperation:
				logger.Println("no operation")
			default:
				if s.debug {
					fmt.Printf("%T\n", el)
					break
				}

				if err := el.Do(); err != nil {
					fmt.Printf("error in %T.Do: %v\n", el, err)
					if s.elementError != nil {
						s.elementError(s, el, err)
					}
				}
			}
			s.seqWait <- struct{}{}
		}
	}
}

func (s *Sequencer) sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	select {
	case <-t.C:
		logger.Println("the waiting had an end")
	case <-s.pause:
		s.waitForResume(ctx)

		if s.state == Idle {
			return
		}

		<-t.C
	}
}

func (s *Sequencer) Abort() {
	s.state = Idle
	s.hasSequence = false

	select {
	case s.pause <- struct{}{}:
	default:
	}
}

func (s *Sequencer) Debug(v bool) {
	s.debug = v
}
