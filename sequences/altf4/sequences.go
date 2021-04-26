package altf4

import (
	"rust-roamer/key"
	"rust-roamer/mouse"
	"rust-roamer/sequencer"
	"rust-roamer/sequences/general"
	"time"
)

const keyReleaseDelay = 60

func RunAll() []sequencer.Elem {
	return general.Flatten([][]sequencer.Elem{
		Run0(),
		{sequencer.Wait{Duration: 2000 * time.Millisecond}},
		Run1(),
		{sequencer.Wait{Duration: 0 * time.Millisecond}},
		Run2(),
		{sequencer.Wait{Duration: 0 * time.Millisecond}},
		Run3(),
		{sequencer.Wait{Duration: 0 * time.Millisecond}},
		Run4(),
	})
}

func Run0() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{Key: key.VK_ESC},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_ESC},
		sequencer.Wait{Duration: 600 * time.Millisecond},
		general.SetMousePos{Pos: mouse.Pos{X: 178, Y: 726}},
		sequencer.Wait{Duration: 100 * time.Millisecond},
		general.LeftMouseButtonDown{},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.LeftMouseButtonUp{},
		sequencer.Wait{Duration: 600 * time.Millisecond},
	}
}

func Run1() []sequencer.Elem {
	return []sequencer.Elem{
		//Turn towards edge and walk to branch
		general.MouseMove{X: 115, Y: 200},
		general.KeyDown{Key: key.VK_W},
		sequencer.Wait{Duration: 5300 * time.Millisecond},
		general.KeyUp{Key: key.VK_W},
		// Jump over branch
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_SPACE},
		// Wait for jump to finish
		sequencer.Wait{Duration: 1700 * time.Millisecond},

		// walk towards edge
		general.KeyDown{Key: key.VK_W},
		sequencer.Wait{Duration: 1200 * time.Millisecond},
		// Turn towards dark rock
		general.MouseMove{X: -94, Y: 0},
		sequencer.Wait{Duration: 30 * time.Millisecond},
		// Jump onto dark rock
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_SPACE},
		// Wait for jump to finish
		sequencer.Wait{Duration: 1700 * time.Millisecond},
		general.KeyUp{Key: key.VK_W},

		// Turn towards next edge in an angle to not get stock between dark rock and brighter rock
		general.MouseMove{X: -60, Y: 0},
		sequencer.Wait{Duration: 30 * time.Millisecond},
		// walk towards edge
		general.KeyDown{Key: key.VK_W},
		sequencer.Wait{Duration: 1250 * time.Millisecond},
		// Jump
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_SPACE},
		sequencer.Wait{Duration: 1700 * time.Millisecond},

		// Turn towards next 2nd rotating spikes
		general.MouseMove{X: 190, Y: 0},
		sequencer.Wait{Duration: 150 * time.Millisecond},

		// Jump over 2nd rotating spikes using the power of right strafe in air
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_SPACE},
		general.KeyDown{Key: key.VK_D},
		sequencer.Wait{Duration: 1630 * time.Millisecond},
		general.KeyUp{Key: key.VK_D},
		sequencer.Wait{Duration: 100 * time.Millisecond},
		general.KeyUp{Key: key.VK_W},

		// Turn to land right of 1st square spikes
		general.MouseMove{X: -40, Y: 0},
		sequencer.Wait{Duration: 50 * time.Millisecond},

		// Jump next to 1st square spikes
		general.KeyDown{Key: key.VK_W},
		sequencer.Wait{Duration: 400 * time.Millisecond},
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_SPACE},
		sequencer.Wait{Duration: 1650 * time.Millisecond},
		general.KeyUp{Key: key.VK_W},

		// Turn 90 right towards 2nd square spikes
		general.MouseMove{X: 535, Y: 0},
		sequencer.Wait{Duration: 50 * time.Millisecond},

		// Jump next to 2nd square spikes
		general.KeyDown{Key: key.VK_W},
		sequencer.Wait{Duration: 300 * time.Millisecond},
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_SPACE},
		sequencer.Wait{Duration: 1300 * time.Millisecond},

		// turn camera left for a good straight run and a bit up so player can time 2nd step
		general.MouseMove{Y: -140},
		sequencer.Wait{Duration: 100 * time.Millisecond},
		general.MouseMove{X: -120},
		sequencer.Wait{Duration: 200 * time.Millisecond},
		general.KeyUp{Key: key.VK_W},
	}
}

func Run2() []sequencer.Elem {
	return []sequencer.Elem{

		general.KeyDown{Key: key.VK_W},
		sequencer.Wait{Duration: 500 * time.Millisecond},
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_SPACE},
		sequencer.Wait{Duration: 2200 * time.Millisecond},
		general.MouseMove{Y: -100},
		sequencer.Wait{Duration: 100 * time.Millisecond},
		general.MouseMove{X: -80},
		sequencer.Wait{Duration: 100 * time.Millisecond},
		sequencer.Wait{Duration: 1600 * time.Millisecond},
		general.KeyUp{Key: key.VK_W},
	}
}

func Run3() []sequencer.Elem {
	return []sequencer.Elem{

		general.KeyDown{Key: key.VK_W},
		//general.KeyDown{Key: key.VK_D},
		sequencer.Wait{Duration: 600 * time.Millisecond},
		//general.KeyUp{Key: key.VK_D},
		general.MouseMove{X: -50, Y: 0},
		sequencer.Wait{Duration: 800 * time.Millisecond},
		general.MouseMove{X: -60, Y: 0},
		sequencer.Wait{Duration: 1000 * time.Millisecond},
		general.MouseMove{X: -40, Y: 100},
		sequencer.Wait{Duration: 600 * time.Millisecond},
		general.KeyUp{Key: key.VK_W},
	}
}

func Run4() []sequencer.Elem {
	return []sequencer.Elem{

		general.KeyDown{Key: key.VK_W},
		sequencer.Wait{Duration: 330 * time.Millisecond},
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_SPACE},
		sequencer.Wait{Duration: 500 * time.Millisecond},
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_SPACE},
		sequencer.Wait{Duration: 800 * time.Millisecond},

		sequencer.Wait{Duration: 100 * time.Millisecond},
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_SPACE},
		sequencer.Wait{Duration: 1200 * time.Millisecond},

		general.KeyUp{Key: key.VK_W},

		sequencer.Wait{Duration: 7000 * time.Millisecond},
		general.MouseMove{X: -240, Y: 0},
		general.KeyDown{Key: key.VK_W},
		sequencer.Wait{Duration: 2200 * time.Millisecond},
		general.KeyUp{Key: key.VK_W},
		general.MouseMove{X: -340, Y: 0},
	}
}

func Run5() []sequencer.Elem {

	return []sequencer.Elem{

		general.KeyDown{Key: key.VK_W},
		sequencer.Wait{Duration: 500 * time.Millisecond},
		general.KeyDown{Key: key.VK_SPACE},
		sequencer.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general.KeyUp{Key: key.VK_SPACE},
		sequencer.Wait{Duration: 800 * time.Millisecond},
		general.KeyDown{Key: key.VK_SPACE},

		general.KeyUp{Key: key.VK_W},
	}
}
