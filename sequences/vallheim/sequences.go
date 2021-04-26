package vallheim

import (
	"rust-roamer/key"
	"rust-roamer/sequencer"
	"rust-roamer/sequences/general"
)

func Run() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{Key: key.VK_W},
		general.KeyDown{Key: key.VK_LSHIFT},
	}
}

func Walk() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{Key: key.VK_W},
	}
}

func Grillmaster() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{Key: key.VK_E},
		sequencer.Wait{Duration: general.HumanizedMillis(200)},
		general.KeyUp{Key: key.VK_E},
		sequencer.Wait{Duration: general.HumanizedMillis(300)},

		general.KeyDown{Key: key.VK_E},
		sequencer.Wait{Duration: general.HumanizedMillis(200)},
		general.KeyUp{Key: key.VK_E},
		sequencer.Wait{Duration: general.HumanizedMillis(300)},

		sequencer.Wait{Duration: general.HumanizedMillis(26000)},

		general.KeyDown{Key: key.VK_E},
		sequencer.Wait{Duration: general.HumanizedMillis(200)},
		general.KeyUp{Key: key.VK_E},
		sequencer.Wait{Duration: general.HumanizedMillis(300)},

		sequencer.Wait{Duration: general.HumanizedMillis(200)},
		general.KeyDown{Key: key.VK_E},
		sequencer.Wait{Duration: general.HumanizedMillis(300)},
		general.KeyUp{Key: key.VK_E},

		sequencer.Wait{Duration: general.HumanizedMillis(600)},

		sequencer.Loop{},
	}
}
