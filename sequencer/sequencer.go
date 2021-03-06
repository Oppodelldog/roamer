package sequencer

import (
	"errors"
	"fmt"
	"rust-roamer/key"
	"rust-roamer/mouse"
	"time"
)

type Sequencer struct {
	sequence chan func() []Elem
	seq      chan Elem
	pause    chan struct{}
	paused   bool
}

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

func (s *Sequencer) IsPaused() bool {
	return s.paused
}

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

func (s *Sequencer) Pause() error {
	select {
	case s.pause <- struct{}{}:
		s.paused = !s.paused
		return nil
	default:
		return errors.New("pause request is already in progress")
	}
}

func (s *Sequencer) WaitForResume() {
	fmt.Println("pausing")
	<-s.pause
	fmt.Println("resuming")
}

func (s *Sequencer) playSequence() {

waitForNext:
	newSequence := <-s.sequence

loop:

	newSeq := newSequence()
	fmt.Println("play sequence length: ", len(newSeq))
	for _, e := range newSeq {
		s.seq <- e
	}

	select {
	case newSequence = <-s.sequence:
		fmt.Println("got a new sequence")
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
			s.WaitForResume()
		case el := <-s.seq:
			switch v := el.(type) {
			case Wait:
				fmt.Println("wait ", v.Duration)
				s.sleep(v.Duration)
			case KeyDown:
				fmt.Println("down ", v.Key)
				key.Down(v.Key)
			case KeyUp:
				fmt.Println("up ", v.Key)
				key.Up(v.Key)
			case LeftMouseButtonDown:
				fmt.Println("lmb-down")
				mouse.LeftDown()
			case LeftMouseButtonUp:
				fmt.Println("lmb-up")
				mouse.LeftUp()
			case RightMouseButtonDown:
				fmt.Println("rmb-down")
				mouse.RightDown()
			case RightMouseButtonUp:
				fmt.Println("rmb-up")
				mouse.RightUp()
			case LookupMousePos:
				pos := mouse.GetCursorPos()
				fmt.Printf("current mouse pos: %#v\n", pos)
			case SetMousePos:
				mouse.SetPosition(v.Pos)
				fmt.Printf("set mouse pos to: %#v\n", v.Pos)
			}
		}
	}
}

func (s *Sequencer) sleep(d time.Duration) {
	t := time.NewTimer(d)
	select {
	case <-t.C:
	case <-s.pause:
		s.WaitForResume()
		<-t.C
	}
}
