package sequences

import (
	"math/rand"
	"rust-roamer/key"
	"rust-roamer/sequencer"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

const paddleTime = 2000
const paddleWait = 800

func GetSequenceFunc(name string) func() []sequencer.Elem {
	return func() []sequencer.Elem {
		var sequences = map[string][]sequencer.Elem{
			"none": {},
			"click-left": {
				sequencer.LeftMouseButtonDown{},
				sequencer.LeftMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(300)},
				sequencer.Loop{},
			},
			"run-forward": {
				sequencer.KeyDown{Key: key.VK_LSHIFT},
				sequencer.KeyDown{Key: key.VK_W},
			},
			"run-forward-stop": {
				sequencer.KeyUp{Key: key.VK_LSHIFT},
				sequencer.KeyUp{Key: key.VK_W},
			},
			"kayak-forward": {
				sequencer.KeyDown{Key: key.VK_A},
				sequencer.Wait{Duration: humanizedMillis(paddleTime + 240)},
				sequencer.KeyUp{Key: key.VK_A},
				sequencer.Wait{Duration: humanizedMillis(paddleWait + 240)},
				sequencer.KeyDown{Key: key.VK_D},
				sequencer.Wait{Duration: humanizedMillis(paddleTime)},
				sequencer.KeyUp{Key: key.VK_D},
				sequencer.Wait{Duration: humanizedMillis(paddleWait)},
				sequencer.Loop{},
			},
			"kayak-backward": {
				sequencer.KeyDown{Key: key.VK_S},
				sequencer.KeyDown{Key: key.VK_A},
				sequencer.Wait{Duration: humanizedMillis(paddleTime + 240)},
				sequencer.KeyUp{Key: key.VK_A},
				sequencer.Wait{Duration: humanizedMillis(paddleWait + 240)},
				sequencer.KeyDown{Key: key.VK_D},
				sequencer.Wait{Duration: humanizedMillis(paddleTime)},
				sequencer.KeyUp{Key: key.VK_D},
				sequencer.Wait{Duration: humanizedMillis(paddleWait)},
				sequencer.KeyUp{Key: key.VK_S},
				sequencer.Loop{},
			},
		}

		return sequences[name]
	}
}

func humanizedMillis(v int) time.Duration {
	var d = v / 10
	var v1 = v - d

	var v2 = v1 + r.Intn(d)*2

	return time.Millisecond * time.Duration(v2)
}
