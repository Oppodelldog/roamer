package inputmode_test

import (
	"testing"

	"github.com/Oppodelldog/roamer/internal/input"
	"github.com/Oppodelldog/roamer/internal/inputmode"
)

func TestInputModeName(t *testing.T) {
	inputmode.SetDryRun(false)
	if inputmode.Name() != "live" {
		t.Fatalf("expected live mode, got %q", inputmode.Name())
	}

	inputmode.SetDryRun(true)
	if inputmode.Name() != "dry-run" {
		t.Fatalf("expected dry-run mode, got %q", inputmode.Name())
	}
	if _, ok := input.Current().(*input.DryRunExecutor); !ok {
		t.Fatalf("expected dry-run executor, got %T", input.Current())
	}

	inputmode.SetDryRun(false)
	if _, ok := input.Current().(input.WinAPIExecutor); !ok {
		t.Fatalf("expected WinAPI executor, got %T", input.Current())
	}
}
