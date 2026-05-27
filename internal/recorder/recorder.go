package recorder

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/mouse"
)

var ErrAlreadyRecording = errors.New("recorder is already running")
var ErrNotRecording = errors.New("recorder is not running")

type State struct {
	Active   bool
	Sequence string
	Count    int
}

type Event struct {
	Label string
	Kind  string
}

type InputState struct {
	Keys  map[int]bool
	Left  bool
	Right bool
	Pos   mouse.Pos
}

type Sampler interface {
	Sample() InputState
}

type Recorder struct {
	mu        sync.Mutex
	sampler   Sampler
	active    bool
	stop      chan struct{}
	commands  []string
	commandAt []time.Time
	lastState InputState
	lastAt    time.Time
	sequence  string
	events    chan Event
}

func New(sampler Sampler) *Recorder {
	return &Recorder{
		sampler: sampler,
		events:  make(chan Event, 64),
	}
}

func (r *Recorder) Events() <-chan Event {
	return r.events
}

func (r *Recorder) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.active {
		return ErrAlreadyRecording
	}

	r.active = true
	r.commands = []string{}
	r.commandAt = []time.Time{}
	r.sequence = ""
	r.lastState = r.sampler.Sample()
	r.lastAt = time.Now()
	r.stop = make(chan struct{})

	go r.loop(r.stop)

	return nil
}

func (r *Recorder) Stop() (State, error) {
	r.mu.Lock()
	if !r.active {
		state := r.stateLocked()
		r.mu.Unlock()
		return state, ErrNotRecording
	}

	stop := r.stop
	r.active = false
	r.mu.Unlock()

	close(stop)

	r.mu.Lock()
	defer r.mu.Unlock()

	r.trimTrailingStopClick(time.Now())
	r.sequence = strings.Join(r.commands, ";")
	return r.stateLocked(), nil
}

func (r *Recorder) State() State {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.stateLocked()
}

func (r *Recorder) loop(stop <-chan struct{}) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			r.sample()
		}
	}
}

func (r *Recorder) sample() {
	now := time.Now()
	next := r.sampler.Sample()

	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.active {
		return
	}

	for _, recordable := range key.RecordableKeys() {
		previous := r.lastState.Keys[recordable.Code]
		current := next.Keys[recordable.Code]
		if previous == current {
			continue
		}

		if current {
			r.appendCommand(now, fmt.Sprintf("KD %s", key.Name(recordable.Code)))
		} else {
			r.appendCommand(now, fmt.Sprintf("KU %s", key.Name(recordable.Code)))
		}
	}

	if r.lastState.Left != next.Left {
		if next.Left {
			r.appendCommand(now, fmt.Sprintf("SM %d %d", next.Pos.X, next.Pos.Y))
			r.appendCommand(now, "LD")
		} else {
			r.appendCommand(now, "LU")
		}
	}

	if r.lastState.Right != next.Right {
		if next.Right {
			r.appendCommand(now, fmt.Sprintf("SM %d %d", next.Pos.X, next.Pos.Y))
			r.appendCommand(now, "RD")
		} else {
			r.appendCommand(now, "RU")
		}
	}

	r.lastState = next
}

func (r *Recorder) appendCommand(now time.Time, command string) {
	if wait := roundedWait(now.Sub(r.lastAt)); wait >= 20*time.Millisecond {
		waitCommand := "W " + wait.String()
		r.commands = append(r.commands, waitCommand)
		r.commandAt = append(r.commandAt, now)
		r.emit(Event{Label: waitCommand, Kind: "wait"})
	}

	r.commands = append(r.commands, command)
	r.commandAt = append(r.commandAt, now)
	r.emit(Event{Label: command, Kind: recorderEventKind(command)})
	r.sequence = strings.Join(r.commands, ";")
	r.lastAt = now
}

func (r *Recorder) emit(event Event) {
	select {
	case r.events <- event:
	default:
	}
}

func recorderEventKind(command string) string {
	switch {
	case strings.HasPrefix(command, "KD "), strings.HasPrefix(command, "KU "):
		return "key"
	case strings.HasPrefix(command, "SM "):
		return "position"
	case command == "LD" || command == "LU" || command == "RD" || command == "RU":
		return "mouse"
	default:
		return "command"
	}
}

func (r *Recorder) trimTrailingStopClick(now time.Time) {
	if len(r.commands) < 2 || len(r.commandAt) != len(r.commands) {
		return
	}

	lastIdx := len(r.commands) - 1
	if now.Sub(r.commandAt[lastIdx]) > 750*time.Millisecond {
		return
	}

	downCommand := ""
	switch r.commands[lastIdx] {
	case "LU":
		downCommand = "LD"
	case "RU":
		downCommand = "RD"
	default:
		return
	}

	downIdx := lastIdx - 1
	if downIdx >= 0 && strings.HasPrefix(r.commands[downIdx], "W ") {
		downIdx--
	}
	if downIdx < 0 || r.commands[downIdx] != downCommand {
		return
	}

	startIdx := downIdx
	if startIdx > 0 && strings.HasPrefix(r.commands[startIdx-1], "SM ") {
		startIdx--
	}
	if startIdx > 0 && strings.HasPrefix(r.commands[startIdx-1], "W ") {
		startIdx--
	}

	r.commands = r.commands[:startIdx]
	r.commandAt = r.commandAt[:startIdx]
}

func (r *Recorder) stateLocked() State {
	return State{
		Active:   r.active,
		Sequence: r.sequence,
		Count:    len(r.commands),
	}
}

func roundedWait(d time.Duration) time.Duration {
	return d.Round(10 * time.Millisecond)
}
