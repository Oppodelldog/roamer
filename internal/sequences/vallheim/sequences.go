package vallheim

import (
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
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

func Jump() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: general.HumanizedMillis(30)},
		general.KeyUp{Key: key.VK_SPACE},
		sequencer.Wait{Duration: general.HumanizedMillis(800)},
		sequencer.Loop{},
	}
}

func Gather() []sequencer.Elem {
	return []sequencer.Elem{
		general.LeftMouseButtonDown{},
		sequencer.Wait{Duration: general.HumanizedMillis(60)},
		general.LeftMouseButtonUp{},
		sequencer.Wait{Duration: general.HumanizedMillis(800)},
		sequencer.Loop{},
	}
}

func Grillmaster() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{Key: key.VK_E},
		sequencer.Wait{Duration: general.HumanizedMillis(300)},
		general.KeyUp{Key: key.VK_E},
		sequencer.Wait{Duration: general.HumanizedMillis(300)},

		general.KeyDown{Key: key.VK_E},
		sequencer.Wait{Duration: general.HumanizedMillis(300)},
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
