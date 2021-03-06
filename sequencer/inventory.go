package sequencer

import (
	"fmt"
	"rust-roamer/mouse"
)

func ItemSlotPos(i int) mouse.Pos {
	switch i {
	case 1:
		return mouse.Pos{X: 707, Y: 1003}
	case 2:
		return mouse.Pos{X: 806, Y: 1003}
	case 3:
		return mouse.Pos{X: 898, Y: 1003}
	case 4:
		return mouse.Pos{X: 992, Y: 1003}
	case 5:
		return mouse.Pos{X: 1086, Y: 1003}
	case 6:
		return mouse.Pos{X: 1170, Y: 1003}
	}

	panic(fmt.Sprintf("invalid item slot: %v", i))
}

func InventorySlotPos(i int) mouse.Pos {
	switch i {
	case 1:
		return mouse.Pos{X: 707, Y: 615}
	case 2:
		return mouse.Pos{X: 799, Y: 615}
	case 3:
		return mouse.Pos{X: 901, Y: 615}
	case 4:
		return mouse.Pos{X: 992, Y: 615}
	case 5:
		return mouse.Pos{X: 1082, Y: 615}
	case 6:
		return mouse.Pos{X: 1178, Y: 615}
	case 7:
		return mouse.Pos{X: 707, Y: 712}
	case 8:
		return mouse.Pos{X: 799, Y: 712}
	case 9:
		return mouse.Pos{X: 901, Y: 712}
	case 10:
		return mouse.Pos{X: 992, Y: 712}
	case 11:
		return mouse.Pos{X: 1082, Y: 712}
	case 12:
		return mouse.Pos{X: 1178, Y: 712}
	case 13:
		return mouse.Pos{X: 707, Y: 810}
	case 14:
		return mouse.Pos{X: 799, Y: 810}
	case 15:
		return mouse.Pos{X: 901, Y: 810}
	case 16:
		return mouse.Pos{X: 992, Y: 810}
	case 17:
		return mouse.Pos{X: 1082, Y: 810}
	case 18:
		return mouse.Pos{X: 1178, Y: 810}
	case 19:
		return mouse.Pos{X: 707, Y: 899}
	case 20:
		return mouse.Pos{X: 799, Y: 899}
	case 21:
		return mouse.Pos{X: 901, Y: 899}
	case 22:
		return mouse.Pos{X: 992, Y: 899}
	case 32:
		return mouse.Pos{X: 1082, Y: 899}
	case 24:
		return mouse.Pos{X: 1178, Y: 899}
	}

	panic(fmt.Sprintf("invalid item slot: %v", i))
}
