package sevendaystodie

import (
	config2 "github.com/Oppodelldog/roamer/internal/config"
	key2 "github.com/Oppodelldog/roamer/internal/key"
	mouse2 "github.com/Oppodelldog/roamer/internal/mouse"
	sequencer2 "github.com/Oppodelldog/roamer/internal/sequencer"
	general2 "github.com/Oppodelldog/roamer/internal/sequences/general"
	version2 "github.com/Oppodelldog/roamer/internal/sequences/version"
	"time"
)

const game = "7d2d"
const slotsBelt = "belt"

func Walk() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyDown{key2.VK_W},
	}
}

func Run() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyDown{key2.VK_SHIFT},
		general2.KeyDown{key2.VK_W},
	}
}

func WalkRunStop() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyUp{key2.VK_SHIFT},
		general2.KeyUp{key2.VK_W},
	}
}

func Repair(slotNo int) []sequencer2.Elem {
	var slots = getSlots(slotsBelt)

	return general2.Flatten([][]sequencer2.Elem{
		{
			general2.KeyDown{key2.VK_TAB},
			general2.KeyUp{key2.VK_TAB},
			sequencer2.Wait{Duration: general2.HumanizedMillis(300)},
			general2.SetMousePos{mousePos(slots.Get(slotNo, 0))},
			sequencer2.Wait{Duration: general2.HumanizedMillis(100)},
			general2.LeftMouseButtonDown{},
			sequencer2.Wait{Duration: general2.HumanizedMillis(100)},
			general2.LeftMouseButtonUp{},
			sequencer2.Wait{Duration: general2.HumanizedMillis(100)},
			general2.KeyDown{key2.VK_A},
			general2.KeyUp{key2.VK_A},
			sequencer2.Wait{Duration: general2.HumanizedMillis(300)},
			general2.KeyDown{key2.VK_TAB},
			general2.KeyUp{key2.VK_TAB},
		},
	})
}

func ClickingLeft(delay int) []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.LeftMouseButtonDown{},
		general2.LeftMouseButtonUp{},
		sequencer2.Wait{Duration: time.Millisecond * time.Duration(delay)},
		sequencer2.Loop{},
	}
}

func getSlots(kind string) config2.Slots {
	return config2.GetSlots(version2.Get(), game, kind)
}

func mousePos(pos config2.Pos) mouse2.Pos {
	return mouse2.Pos{
		X: int32(pos.X),
		Y: int32(pos.Y),
	}
}
