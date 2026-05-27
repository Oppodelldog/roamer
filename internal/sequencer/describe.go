package sequencer

import (
	"fmt"
	"time"
)

type ElementEvent struct {
	Label      string
	DurationMs int64
}

func DescribeElement(elem Elem) ElementEvent {
	switch v := elem.(type) {
	case Wait:
		return ElementEvent{Label: "W " + formatDuration(v.Duration), DurationMs: v.Duration.Milliseconds()}
	case Loop:
		return ElementEvent{Label: "L"}
	case NoOperation:
		return ElementEvent{Label: "NOP"}
	case fmt.Stringer:
		return ElementEvent{Label: v.String()}
	default:
		return ElementEvent{Label: fmt.Sprintf("%T", elem)}
	}
}

func formatDuration(d time.Duration) string {
	if d%time.Minute == 0 && d >= time.Minute {
		return fmt.Sprintf("%dm", int(d/time.Minute))
	}

	if d%time.Second == 0 && d >= time.Second {
		return fmt.Sprintf("%ds", int(d/time.Second))
	}

	return fmt.Sprintf("%dms", d.Milliseconds())
}
