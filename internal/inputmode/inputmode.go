package inputmode

import (
	"sync/atomic"

	"github.com/Oppodelldog/roamer/internal/input"
)

var dryRun atomic.Bool

func IsDryRun() bool {
	return dryRun.Load()
}

func SetDryRun(v bool) {
	dryRun.Store(v)
	if v {
		input.UseDryRun()
		return
	}

	input.UseLive()
}

func Name() string {
	if IsDryRun() {
		return "dry-run"
	}

	return "live"
}
