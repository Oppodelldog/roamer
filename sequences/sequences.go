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
const inventoryOpenWaitShort = 200
const inventoryMoveWait = 60
const mouseMoveWait = 60
const mouseClickWait = 0
const mouseDragWait = 160

func GetSequenceFunc(name string) func() []sequencer.Elem {
	return func() []sequencer.Elem {
		var sequences = map[string][]sequencer.Elem{
			"none": {},
			"get-mouse-pos": {
				LookupMousePos{},
			},
			"click-left": {
				LeftMouseButtonDown{},
				LeftMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(300)},
				sequencer.Loop{},
			},
			"run-forward": {
				KeyDown{Key: key.VK_LSHIFT},
				KeyDown{Key: key.VK_W},
				sequencer.Wait{Duration: humanizedMillis(3000)},
				KeyUp{Key: key.VK_LSHIFT},
			},
			"run-forward-stop": {
				KeyUp{Key: key.VK_LSHIFT},
				KeyUp{Key: key.VK_W},
			},
			"kayak-forward": {
				KeyDown{Key: key.VK_A},
				sequencer.Wait{Duration: humanizedMillis(paddleTime + 240)},
				KeyUp{Key: key.VK_A},
				sequencer.Wait{Duration: humanizedMillis(paddleWait + 240)},
				KeyDown{Key: key.VK_D},
				sequencer.Wait{Duration: humanizedMillis(paddleTime)},
				KeyUp{Key: key.VK_D},
				sequencer.Wait{Duration: humanizedMillis(paddleWait)},
				sequencer.Loop{},
			},
			"kayak-backward": {
				KeyDown{Key: key.VK_S},
				KeyDown{Key: key.VK_A},
				sequencer.Wait{Duration: humanizedMillis(paddleTime + 240)},
				KeyUp{Key: key.VK_A},
				sequencer.Wait{Duration: humanizedMillis(paddleWait + 240)},
				KeyDown{Key: key.VK_D},
				sequencer.Wait{Duration: humanizedMillis(paddleTime)},
				KeyUp{Key: key.VK_D},
				sequencer.Wait{Duration: humanizedMillis(paddleWait)},
				KeyUp{Key: key.VK_S},
				sequencer.Loop{},
			},
			"duck-gather-tree": {
				KeyDown{Key: key.VK_LCONTROL},
				LeftMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(7000)},
				LeftMouseButtonUp{},
				KeyUp{Key: key.VK_LCONTROL},
			},
			"unarm": {
				KeyDown{Key: key.VK_TAB},
				KeyUp{Key: key.VK_TAB},
				sequencer.Wait{Duration: humanizedMillis(inventoryOpenWait)},
				SetMousePos{Pos: ItemSlotPos(1)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				SetMousePos{Pos: ItemSlotPos(2)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				SetMousePos{Pos: ItemSlotPos(3)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				SetMousePos{Pos: ItemSlotPos(4)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				SetMousePos{Pos: ItemSlotPos(5)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				SetMousePos{Pos: ItemSlotPos(6)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait * 10)},

				KeyDown{Key: key.VK_TAB},
				KeyUp{Key: key.VK_TAB},
			},
			"arm": {
				KeyDown{Key: key.VK_TAB},
				KeyUp{Key: key.VK_TAB},
				sequencer.Wait{Duration: humanizedMillis(inventoryOpenWait)},
				SetMousePos{Pos: InventorySlotPos(1)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				SetMousePos{Pos: InventorySlotPos(2)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				SetMousePos{Pos: InventorySlotPos(3)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				SetMousePos{Pos: InventorySlotPos(4)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				SetMousePos{Pos: InventorySlotPos(5)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},

				SetMousePos{Pos: InventorySlotPos(6)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				RightMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
				RightMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait * 10)},

				KeyDown{Key: key.VK_TAB},
				KeyUp{Key: key.VK_TAB},
			},
			"diving-tank-on": {
				KeyDown{Key: key.VK_TAB},
				KeyUp{Key: key.VK_TAB},
				sequencer.Wait{Duration: humanizedMillis(inventoryOpenWaitShort)},
				SetMousePos{Pos: InventorySlotPos(19)},
				LeftMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseDragWait)},
				SetMousePos{Pos: clothingSlotPos(7)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				LeftMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryOpenWaitShort)},
				KeyDown{Key: key.VK_TAB},
				KeyUp{Key: key.VK_TAB},
			},
			"diving-tank-off": {
				KeyDown{Key: key.VK_TAB},
				KeyUp{Key: key.VK_TAB},
				sequencer.Wait{Duration: humanizedMillis(inventoryOpenWaitShort)},
				SetMousePos{Pos: clothingSlotPos(7)},
				LeftMouseButtonDown{},
				sequencer.Wait{Duration: humanizedMillis(mouseDragWait)},
				SetMousePos{Pos: InventorySlotPos(19)},
				sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
				LeftMouseButtonUp{},
				sequencer.Wait{Duration: humanizedMillis(inventoryOpenWaitShort)},
				KeyDown{Key: key.VK_TAB},
				KeyUp{Key: key.VK_TAB},
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
