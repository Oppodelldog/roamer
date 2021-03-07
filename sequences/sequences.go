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

const paddleTime = 1200
const paddleWait = 400
const inventoryOpenWait = 600
const inventoryMoveWait = 60
const mouseMoveWait = 60
const mouseClickWait = 0

func GetSequenceFunc(name string) func() []sequencer.Elem {
	return func() []sequencer.Elem {
		var sequences = map[string][]sequencer.Elem{
			"none": {},
			"get-mouse-pos": {
				sequencer.LookupMousePos{},
			},
			"click-left": {
				sequencer.LeftMouseButtonDown{},
				sequencer.LeftMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(300)},
				sequencer.Loop{},
			},
			"run-forward": {
				sequencer.KeyDown{Key: key.VK_LSHIFT},
				sequencer.KeyDown{Key: key.VK_W},
				sequencer.Wait{Duration: humanizedMillis(3000)},
				sequencer.KeyUp{Key: key.VK_LSHIFT},
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
			"duck-gather-tree": {
				sequencer.KeyDown{Key: key.VK_LCONTROL},
				sequencer.LeftMouseButtonDown{},
				sequencer.Wait{Duration: 7600 * time.Millisecond},
				sequencer.LeftMouseButtonUp{},
				sequencer.KeyUp{Key: key.VK_LCONTROL},
			},
			"unarm": {
				sequencer.KeyDown{Key: key.VK_TAB},
				sequencer.KeyUp{Key: key.VK_TAB},
				sequencer.Wait{Duration: humanizedMillis(inventoryOpenWait)},
				sequencer.SetMousePos{Pos: sequencer.ItemSlotPos(1)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				sequencer.SetMousePos{Pos: sequencer.ItemSlotPos(2)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				sequencer.SetMousePos{Pos: sequencer.ItemSlotPos(3)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				sequencer.SetMousePos{Pos: sequencer.ItemSlotPos(4)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				sequencer.SetMousePos{Pos: sequencer.ItemSlotPos(5)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				sequencer.SetMousePos{Pos: sequencer.ItemSlotPos(6)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait * 10)},

				sequencer.KeyDown{Key: key.VK_TAB},
				sequencer.KeyUp{Key: key.VK_TAB},
			},
			"arm": {
				sequencer.KeyDown{Key: key.VK_TAB},
				sequencer.KeyUp{Key: key.VK_TAB},
				sequencer.Wait{Duration: humanizedMillis(inventoryOpenWait)},
				sequencer.SetMousePos{Pos: sequencer.InventorySlotPos(1)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				sequencer.SetMousePos{Pos: sequencer.InventorySlotPos(2)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				sequencer.SetMousePos{Pos: sequencer.InventorySlotPos(3)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				sequencer.SetMousePos{Pos: sequencer.InventorySlotPos(4)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				sequencer.SetMousePos{Pos: sequencer.InventorySlotPos(5)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				sequencer.SetMousePos{Pos: sequencer.InventorySlotPos(6)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				sequencer.RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				sequencer.RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait * 10)},

				sequencer.KeyDown{Key: key.VK_TAB},
				sequencer.KeyUp{Key: key.VK_TAB},
			},
		}

		return sequences[name]
	}
}

func humanizedMillis(v int) time.Duration {
	if v == 0 {
		return 0
	}

	var d = v / 10
	var v1 = v - d

	var v2 = v1 + r.Intn(d)*2

	return time.Millisecond * time.Duration(v2)
}
