package altf4

import (
	key2 "github.com/Oppodelldog/roamer/internal/key"
	mouse2 "github.com/Oppodelldog/roamer/internal/mouse"
	sequencer2 "github.com/Oppodelldog/roamer/internal/sequencer"
	general2 "github.com/Oppodelldog/roamer/internal/sequences/general"
	"time"
)

const keyReleaseDelay = 60

func RunAll() []sequencer2.Elem {
	return general2.Flatten([][]sequencer2.Elem{
		Reset(),
		{sequencer2.Wait{Duration: 2000 * time.Millisecond}},
		Run1(),
		{sequencer2.Wait{Duration: 0 * time.Millisecond}},
		Run2(),
		{sequencer2.Wait{Duration: 0 * time.Millisecond}},
		Run3(),
		{sequencer2.Wait{Duration: 0 * time.Millisecond}},
		Run4(),
	})
}

func RunAll2() []sequencer2.Elem {
	return general2.Flatten([][]sequencer2.Elem{
		Reset(),
		{sequencer2.Wait{Duration: 1600 * time.Millisecond}},
		RunDownToPathBranch1(),
		{
			general2.MouseMove{X: -127},
			//sequencer.Wait{Duration:2700 * time.Millisecond},
			general2.KeyDown{Key: key2.VK_W},
			sequencer2.Wait{Duration: 800 * time.Millisecond},
			general2.KeyDown{Key: key2.VK_SPACE},
			sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
			sequencer2.Wait{Duration: 800 * time.Millisecond},
			general2.KeyUp{Key: key2.VK_SPACE},
			sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
			general2.KeyUp{Key: key2.VK_W},
			sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},

			sequencer2.Wait{Duration: 1200 * time.Millisecond},

			general2.KeyDown{Key: key2.VK_W},
			sequencer2.Wait{Duration: 1000 * time.Millisecond},
			general2.KeyDown{Key: key2.VK_SPACE},
			sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
			general2.KeyUp{Key: key2.VK_SPACE},

			sequencer2.Wait{Duration: 1800 * time.Millisecond},
			general2.KeyUp{Key: key2.VK_W},
		},
		{sequencer2.Wait{Duration: 0 * time.Millisecond}},
	})
}

func Reset() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyDown{Key: key2.VK_ESC},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.KeyUp{Key: key2.VK_ESC},
		sequencer2.Wait{Duration: 600 * time.Millisecond},
		general2.SetMousePos{Pos: mouse2.Pos{X: 178, Y: 726}},
		sequencer2.Wait{Duration: 100 * time.Millisecond},
		general2.LeftMouseButtonDown{},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.LeftMouseButtonUp{},
		sequencer2.Wait{Duration: 600 * time.Millisecond},
	}
}

func RunDownToPathBranch1() []sequencer2.Elem {
	return []sequencer2.Elem{
		//Turn towards edge and walk to branch
		general2.MouseMove{X: 115, Y: 200},
		general2.KeyDown{Key: key2.VK_W},
		sequencer2.Wait{Duration: 5300 * time.Millisecond},
		general2.KeyUp{Key: key2.VK_W},
		// Jump over branch
		general2.KeyDown{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.KeyUp{Key: key2.VK_SPACE},
		// Wait for jump to finish
		sequencer2.Wait{Duration: 1700 * time.Millisecond},

		// walk towards edge
		general2.KeyDown{Key: key2.VK_W},
		sequencer2.Wait{Duration: 1200 * time.Millisecond},
		// Turn towards dark rock
		general2.MouseMove{X: -94, Y: 0},
		sequencer2.Wait{Duration: 30 * time.Millisecond},
		// Jump onto dark rock
		general2.KeyDown{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.KeyUp{Key: key2.VK_SPACE},
		// Wait for jump to finish
		sequencer2.Wait{Duration: 1700 * time.Millisecond},
		general2.KeyUp{Key: key2.VK_W},

		// Turn towards next edge in an angle to not get stock between dark rock and brighter rock
		general2.MouseMove{X: -60, Y: 0},
		sequencer2.Wait{Duration: 30 * time.Millisecond},
		// walk towards edge
		general2.KeyDown{Key: key2.VK_W},
		sequencer2.Wait{Duration: 1250 * time.Millisecond},
		// Jump
		general2.KeyDown{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.KeyUp{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: 1700 * time.Millisecond},

		// Turn towards next 2nd rotating spikes
		general2.MouseMove{X: 190, Y: 0},
		sequencer2.Wait{Duration: 150 * time.Millisecond},

		// Jump over 2nd rotating spikes using the power of right strafe in air
		general2.KeyDown{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.KeyUp{Key: key2.VK_SPACE},
		general2.KeyDown{Key: key2.VK_D},
		sequencer2.Wait{Duration: 1630 * time.Millisecond},
		general2.KeyUp{Key: key2.VK_D},
		sequencer2.Wait{Duration: 100 * time.Millisecond},
		general2.KeyUp{Key: key2.VK_W},

		// Turn to land right of 1st square spikes
		general2.MouseMove{X: -40, Y: 0},
		sequencer2.Wait{Duration: 50 * time.Millisecond},

		// Jump next to 1st square spikes
		general2.KeyDown{Key: key2.VK_W},
		sequencer2.Wait{Duration: 400 * time.Millisecond},
		general2.KeyDown{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.KeyUp{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: 1650 * time.Millisecond},
		general2.KeyUp{Key: key2.VK_W},
	}
}

func Run1() []sequencer2.Elem {
	return general2.Flatten([][]sequencer2.Elem{
		RunDownToPathBranch1(),
		{
			// Turn 90 right towards 2nd square spikes
			general2.MouseMove{X: 535, Y: 0},
			sequencer2.Wait{Duration: 50 * time.Millisecond},

			// Jump next to 2nd square spikes
			general2.KeyDown{Key: key2.VK_W},
			sequencer2.Wait{Duration: 300 * time.Millisecond},
			general2.KeyDown{Key: key2.VK_SPACE},
			sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
			general2.KeyUp{Key: key2.VK_SPACE},
			sequencer2.Wait{Duration: 1300 * time.Millisecond},

			// turn camera left for a good straight run and a bit up so player can time 2nd step
			general2.MouseMove{Y: -140},
			sequencer2.Wait{Duration: 100 * time.Millisecond},
			general2.MouseMove{X: -120},
			sequencer2.Wait{Duration: 200 * time.Millisecond},
			general2.KeyUp{Key: key2.VK_W},
		},
	})
}

func Run2() []sequencer2.Elem {
	return []sequencer2.Elem{

		general2.KeyDown{Key: key2.VK_W},
		sequencer2.Wait{Duration: 500 * time.Millisecond},
		general2.KeyDown{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.KeyUp{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: 2200 * time.Millisecond},
		general2.MouseMove{Y: -100},
		sequencer2.Wait{Duration: 100 * time.Millisecond},
		general2.MouseMove{X: -60},
		sequencer2.Wait{Duration: 100 * time.Millisecond},
		sequencer2.Wait{Duration: 1600 * time.Millisecond},
		general2.KeyUp{Key: key2.VK_W},
	}
}

func Run3() []sequencer2.Elem {
	return []sequencer2.Elem{

		general2.KeyDown{Key: key2.VK_W},
		//general.KeyDown{Key: key.VK_D},
		sequencer2.Wait{Duration: 600 * time.Millisecond},
		//general.KeyUp{Key: key.VK_D},
		general2.MouseMove{X: -50, Y: 0},
		sequencer2.Wait{Duration: 800 * time.Millisecond},
		general2.MouseMove{X: -60, Y: 0},
		sequencer2.Wait{Duration: 800 * time.Millisecond},
		general2.MouseMove{X: -40, Y: 100},
		general2.KeyUp{Key: key2.VK_W},
		sequencer2.Wait{Duration: 800 * time.Millisecond},
		general2.KeyDown{Key: key2.VK_W},
		sequencer2.Wait{Duration: 700 * time.Millisecond},
		general2.KeyUp{Key: key2.VK_W},
	}
}

func Run4() []sequencer2.Elem {
	return []sequencer2.Elem{

		general2.KeyDown{Key: key2.VK_W},
		sequencer2.Wait{Duration: 330 * time.Millisecond},
		general2.KeyDown{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.KeyUp{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: 500 * time.Millisecond},
		general2.KeyDown{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.KeyUp{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: 800 * time.Millisecond},

		sequencer2.Wait{Duration: 100 * time.Millisecond},
		general2.KeyDown{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.KeyUp{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: 1200 * time.Millisecond},

		general2.KeyUp{Key: key2.VK_W},

		sequencer2.Wait{Duration: 7000 * time.Millisecond},
		general2.MouseMove{X: -250, Y: 0},
		general2.KeyDown{Key: key2.VK_W},
		sequencer2.Wait{Duration: 2200 * time.Millisecond},
		general2.KeyUp{Key: key2.VK_W},
		general2.MouseMove{X: -350, Y: 0},
	}
}

func Run5() []sequencer2.Elem {

	return []sequencer2.Elem{

		general2.KeyDown{Key: key2.VK_W},
		sequencer2.Wait{Duration: 500 * time.Millisecond},
		general2.KeyDown{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: keyReleaseDelay * time.Millisecond},
		general2.KeyUp{Key: key2.VK_SPACE},
		sequencer2.Wait{Duration: 800 * time.Millisecond},
		general2.KeyDown{Key: key2.VK_SPACE},

		general2.KeyUp{Key: key2.VK_W},
	}
}
