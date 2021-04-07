package sevendaystodie

import (
	"rust-roamer/config"
	"rust-roamer/key"
	"rust-roamer/mouse"
	"rust-roamer/sequencer"
	"rust-roamer/sequences/general"
	"rust-roamer/sequences/version"
	"time"
)

const game = "7d2d"
const slotsBelt = "belt"

func Walk() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{key.VK_W},
	}
}

func Run() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{key.VK_SHIFT},
		general.KeyDown{key.VK_W},
	}
}

func WalkRunStop() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyUp{key.VK_SHIFT},
		general.KeyUp{key.VK_W},
	}
}

func Repair(slotNo int) []sequencer.Elem {
	var slots = getSlots(slotsBelt)

	return general.Flatten([][]sequencer.Elem{
		{
			general.KeyDown{key.VK_TAB},
			general.KeyUp{key.VK_TAB},
			sequencer.Wait{Duration: general.HumanizedMillis(300)},
			general.SetMousePos{mousePos(slots.Get(slotNo, 0))},
			sequencer.Wait{Duration: general.HumanizedMillis(100)},
			general.LeftMouseButtonDown{},
			sequencer.Wait{Duration: general.HumanizedMillis(100)},
			general.LeftMouseButtonUp{},
			sequencer.Wait{Duration: general.HumanizedMillis(100)},
			general.KeyDown{key.VK_A},
			general.KeyUp{key.VK_A},
			sequencer.Wait{Duration: general.HumanizedMillis(300)},
			general.KeyDown{key.VK_TAB},
			general.KeyUp{key.VK_TAB},
		},
	})
}

func ClickingLeft(delay int) []sequencer.Elem {
	return []sequencer.Elem{
		general.LeftMouseButtonDown{},
		general.LeftMouseButtonUp{},
		sequencer.Wait{Duration: time.Millisecond * time.Duration(delay)},
		sequencer.Loop{},
	}
}

func getSlots(kind string) config.Slots {
	return config.GetSlots(version.Get(), game, kind)
}

func mousePos(pos config.Pos) mouse.Pos {
	return mouse.Pos{
		X: int32(pos.X),
		Y: int32(pos.Y),
	}
}
