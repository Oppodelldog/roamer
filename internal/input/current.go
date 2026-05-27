package input

import "sync"

var (
	currentMu sync.RWMutex
	current   Executor = WinAPIExecutor{}
)

func Current() Executor {
	currentMu.RLock()
	defer currentMu.RUnlock()

	return current
}

func UseLive() {
	currentMu.Lock()
	defer currentMu.Unlock()

	current = WinAPIExecutor{}
}

func UseDryRun() {
	currentMu.Lock()
	defer currentMu.Unlock()

	current = NewDryRunExecutor()
}
