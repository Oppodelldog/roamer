package input_test

import (
	"testing"

	"github.com/Oppodelldog/roamer/internal/input"
	"github.com/Oppodelldog/roamer/internal/mouse"
)

func TestDryRunExecutorTracksCursorPosition(t *testing.T) {
	executor := input.NewDryRunExecutor()

	if err := executor.SetPosition(mouse.Pos{X: 10, Y: 20}); err != nil {
		t.Fatalf("did not expect SetPosition error, got %v", err)
	}

	if err := executor.Move(5, -2); err != nil {
		t.Fatalf("did not expect Move error, got %v", err)
	}

	pos, err := executor.GetCursorPos()
	if err != nil {
		t.Fatalf("did not expect GetCursorPos error, got %v", err)
	}

	if pos.X != 15 || pos.Y != 18 {
		t.Fatalf("expected cursor position 15/18, got %#v", pos)
	}
}

func TestCurrentExecutorCanSwitchToDryRun(t *testing.T) {
	input.UseDryRun()
	if _, ok := input.Current().(*input.DryRunExecutor); !ok {
		t.Fatalf("expected dry-run executor, got %T", input.Current())
	}

	input.UseLive()
	if _, ok := input.Current().(input.WinAPIExecutor); !ok {
		t.Fatalf("expected WinAPI executor, got %T", input.Current())
	}
}
