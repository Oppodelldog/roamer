package input

import (
	"sync"

	"github.com/Oppodelldog/roamer/internal/mouse"
)

type DryRunExecutor struct {
	mu        sync.Mutex
	keys      map[int]bool
	leftDown  bool
	rightDown bool
	pos       mouse.Pos
}

func NewDryRunExecutor() *DryRunExecutor {
	return &DryRunExecutor{
		keys: map[int]bool{},
	}
}

func (d *DryRunExecutor) KeyDown(key int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.keys[key] = true

	return nil
}

func (d *DryRunExecutor) KeyUp(key int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.keys[key] = false

	return nil
}

func (d *DryRunExecutor) LeftDown() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.leftDown = true

	return nil
}

func (d *DryRunExecutor) LeftUp() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.leftDown = false

	return nil
}

func (d *DryRunExecutor) RightDown() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.rightDown = true

	return nil
}

func (d *DryRunExecutor) RightUp() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.rightDown = false

	return nil
}

func (d *DryRunExecutor) SetPosition(pos mouse.Pos) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.pos = pos

	return nil
}

func (d *DryRunExecutor) GetCursorPos() (mouse.Pos, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.pos, nil
}

func (d *DryRunExecutor) Move(x, y int32) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.pos.X += x
	d.pos.Y += y

	return nil
}

func (d *DryRunExecutor) ResetPressed() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for key := range d.keys {
		d.keys[key] = false
	}
	d.leftDown = false
	d.rightDown = false
}
