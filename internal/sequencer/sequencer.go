package sequencer

import (
	"errors"
	"fmt"
	"time"
)

// Sequencer plays sequences of Elem by calling their Do method.
// It runs as a detached worker using go routines to playback sequences.
type Sequencer struct {
	sequence       chan func() []Elem
	seq            chan Elem
	pause          chan struct{}
	paused         bool
	aborted        bool
	hasSequence    bool
	beforeSequence func()
}

// New creates a new sequencer.
// Parameter queueSize defines how many Sequences may be queued before blocking.
func New(queueSize int) *Sequencer {
	s := &Sequencer{
		sequence: make(chan func() []Elem, queueSize),
		seq:      make(chan Elem),
		pause:    make(chan struct{}, 1),
		paused:   false,
	}

	go s.playSequence()
	go s.playElement()

	return s
}

func (s *Sequencer) BeforeSequence(f func()) {
	s.beforeSequence = f
}

// IsPaused indicated if playback is paused or not.
func (s *Sequencer) IsPaused() bool {
	return s.paused
}

// HasSequences indicates if sequencer has a sequence for playback or not.
func (s *Sequencer) HasSequence() bool {
	return s.hasSequence
}

// Enqueue adds a new sequence to the sequence queue.
// This method blocks by the queueSize defined in New.
func (s *Sequencer) EnqueueSequence(newSeq func() []Elem) {
	var newSequence = newSeq
	select {
	case s.sequence <- newSequence:
		fmt.Println("sequence enqueued")
	default:
		fmt.Println("dropped waiting sequence")
		<-s.sequence
		s.sequence <- newSequence
	}
}

// Pause toggles playback (pause, resume).
func (s *Sequencer) Pause() error {
	select {
	case s.pause <- struct{}{}:
		s.paused = !s.paused

		return nil
	default:
		return errors.New("pause request is already in progress")
	}
}

func (s *Sequencer) waitForResume() {
	fmt.Println("pausing")
	<-s.pause
	fmt.Println("resuming")
}

func (s *Sequencer) playSequence() {
waitForNext:
	s.hasSequence = false
	s.aborted = false
	newSequence := <-s.sequence
	s.hasSequence = true

	if s.beforeSequence != nil {
		s.beforeSequence()
	}

loop:
	newSeq := newSequence()
	fmt.Println("play sequence length: ", len(newSeq))

	for _, e := range newSeq {
		s.seq <- e
		if s.aborted {
			fmt.Println("aborted sequence")
			goto waitForNext
		}
	}

	select {
	case newSequence = <-s.sequence:
		fmt.Println("got a new sequence")
		goto loop
	default:
		fmt.Println("got no new sequence")
		var newSeq = newSequence()
		if len(newSeq) > 0 {
			if _, isLoop := newSeq[len(newSeq)-1].(Loop); isLoop {
				fmt.Println("looping")
				goto loop
			}
		}
	}

	goto waitForNext
}

func (s *Sequencer) playElement() {
	for {
		select {
		case <-s.pause:
			s.waitForResume()
		case el := <-s.seq:

			switch v := el.(type) {
			case Wait:
				fmt.Println("wait ", v.Duration)
				s.sleep(v.Duration)
			case Loop:
			default:
				err := el.Do()
				if err != nil {
					fmt.Printf("error in %T.Do: %v\n", el, err)
				}
			}
		}
	}
}

func (s *Sequencer) sleep(d time.Duration) {
	t := time.NewTimer(d)
	select {
	case <-t.C:
	case <-s.pause:
		s.waitForResume()
		if s.aborted {
			return
		}
		<-t.C
	}
}

func (s *Sequencer) Abort() {
	if s.paused {
		err := s.Pause()
		if err != nil {
			panic(fmt.Sprintf("error while aborting pause: %v", err))
		}

	}
	s.aborted = true
}
