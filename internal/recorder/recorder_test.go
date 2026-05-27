package recorder

import (
	"strings"
	"testing"
	"time"

	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/mouse"
)

type fakeSampler struct {
	state InputState
}

func (s *fakeSampler) Sample() InputState {
	return s.state
}

func TestRecorderRecordsKeyTransitions(t *testing.T) {
	sampler := &fakeSampler{state: InputState{Keys: map[int]bool{}}}
	rec := New(sampler)

	if err := rec.Start(); err != nil {
		t.Fatalf("expected recorder to start: %v", err)
	}

	time.Sleep(25 * time.Millisecond)
	sampler.state = InputState{Keys: map[int]bool{key.VkW: true}}
	rec.sample()
	time.Sleep(25 * time.Millisecond)
	sampler.state = InputState{Keys: map[int]bool{key.VkW: false}}
	rec.sample()

	state, err := rec.Stop()
	if err != nil {
		t.Fatalf("expected recorder to stop: %v", err)
	}

	if state.Sequence == "" {
		t.Fatal("expected recorded sequence")
	}

	if state.Count < 4 {
		t.Fatalf("expected wait and key transitions, got count %d with %q", state.Count, state.Sequence)
	}
}

func TestRecorderRecordsMousePositionBeforeClick(t *testing.T) {
	sampler := &fakeSampler{state: InputState{Keys: map[int]bool{}, Pos: mouse.Pos{X: 10, Y: 20}}}
	rec := New(sampler)

	if err := rec.Start(); err != nil {
		t.Fatalf("expected recorder to start: %v", err)
	}

	sampler.state = InputState{Keys: map[int]bool{}, Left: true, Pos: mouse.Pos{X: 30, Y: 40}}
	rec.sample()
	rec.mu.Lock()
	for idx := range rec.commandAt {
		rec.commandAt[idx] = time.Now().Add(-time.Second)
	}
	rec.mu.Unlock()

	state, err := rec.Stop()
	if err != nil {
		t.Fatalf("expected recorder to stop: %v", err)
	}

	if state.Sequence != "SM 30 40;LD" {
		t.Fatalf("expected mouse position and left down, got %q", state.Sequence)
	}
}

func TestRecorderEmitsEvents(t *testing.T) {
	sampler := &fakeSampler{state: InputState{Keys: map[int]bool{}}}
	rec := New(sampler)

	if err := rec.Start(); err != nil {
		t.Fatalf("expected recorder to start: %v", err)
	}

	sampler.state = InputState{Keys: map[int]bool{key.VkW: true}}
	rec.sample()

	select {
	case event := <-rec.Events():
		if event.Label != "KD W" || event.Kind != "key" {
			t.Fatalf("unexpected event: %#v", event)
		}
	default:
		t.Fatal("expected recorder event")
	}
}

func TestRecorderTrimsFreshTrailingStopClick(t *testing.T) {
	rec := New(&fakeSampler{state: InputState{Keys: map[int]bool{}}})
	now := time.Now()
	rec.commands = []string{"KD W", "W 100ms", "SM 30 40", "LD", "LU"}
	rec.commandAt = []time.Time{
		now.Add(-2 * time.Second),
		now.Add(-time.Second),
		now.Add(-20 * time.Millisecond),
		now.Add(-20 * time.Millisecond),
		now.Add(-10 * time.Millisecond),
	}

	rec.trimTrailingStopClick(now)

	got := strings.Join(rec.commands, ";")
	if got != "KD W" {
		t.Fatalf("expected trailing stop click to be trimmed, got %q", got)
	}
}
