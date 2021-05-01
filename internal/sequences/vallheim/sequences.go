package vallheim

import (
	key2 "github.com/Oppodelldog/roamer/internal/key"
	sequencer2 "github.com/Oppodelldog/roamer/internal/sequencer"
	general2 "github.com/Oppodelldog/roamer/internal/sequences/general"
)

func Run() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyDown{Key: key2.VK_W},
		general2.KeyDown{Key: key2.VK_LSHIFT},
	}
}

func Walk() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyDown{Key: key2.VK_W},
	}
}

func Grillmaster() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyDown{Key: key2.VK_E},
		sequencer2.Wait{Duration: general2.HumanizedMillis(300)},
		general2.KeyUp{Key: key2.VK_E},
		sequencer2.Wait{Duration: general2.HumanizedMillis(300)},

		general2.KeyDown{Key: key2.VK_E},
		sequencer2.Wait{Duration: general2.HumanizedMillis(300)},
		general2.KeyUp{Key: key2.VK_E},
		sequencer2.Wait{Duration: general2.HumanizedMillis(300)},

		sequencer2.Wait{Duration: general2.HumanizedMillis(26000)},

		general2.KeyDown{Key: key2.VK_E},
		sequencer2.Wait{Duration: general2.HumanizedMillis(200)},
		general2.KeyUp{Key: key2.VK_E},
		sequencer2.Wait{Duration: general2.HumanizedMillis(300)},

		sequencer2.Wait{Duration: general2.HumanizedMillis(200)},
		general2.KeyDown{Key: key2.VK_E},
		sequencer2.Wait{Duration: general2.HumanizedMillis(300)},
		general2.KeyUp{Key: key2.VK_E},

		sequencer2.Wait{Duration: general2.HumanizedMillis(600)},

		sequencer2.Loop{},
	}
}
