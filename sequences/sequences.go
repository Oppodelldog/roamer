package sequences

import (
	"math/rand"
	"rust-roamer/key"
	"rust-roamer/mouse"
	"rust-roamer/sequencer"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

const paddleTime = 1200
const paddleWait = 400
const inventoryOpenWaitShort = 200
const inventoryMoveWait = 60
const mouseMoveWait = 60
const mouseClickWait = 0
const mouseDragWait = 16
const itemFieldMargin = 3
const itemFieldSize = 92

func transferInventoryRow() []sequencer.Elem {
	pos := mouse.GetCursorPos()
	var seqs [][]sequencer.Elem

	for i := int32(0); i < 6; i++ {
		seqs = append(seqs, [][]sequencer.Elem{
			collect(relativeFieldPos(pos, i)),
			{sequencer.Wait{Duration: humanizedMillis(80)}}}...)
	}

	return flatten(seqs)
}

func relativeFieldPos(from mouse.Pos, fieldNo int32) mouse.Pos {
	return mouse.Pos{X: from.X + (fieldNo * (itemFieldSize + itemFieldMargin)), Y: from.Y}
}

func fillInventoryRow() []sequencer.Elem {
	pos := mouse.GetCursorPos()

	var seqs [][]sequencer.Elem

	for i := int32(1); i < 6; i++ {
		seqs = append(seqs, [][]sequencer.Elem{
			unstack(pos, relativeFieldPos(pos, i)),
			{sequencer.Wait{Duration: humanizedMillis(80)}}}...)
	}

	return flatten(seqs)
}

func collect(from mouse.Pos) []sequencer.Elem {
	return []sequencer.Elem{
		SetMousePos{Pos: from},
		RightMouseButtonDown{},
		sequencer.Wait{Duration: humanizedMillis(60)},
		RightMouseButtonUp{},
		sequencer.Wait{Duration: humanizedMillis(60)},
	}
}

func unstack(from, to mouse.Pos) []sequencer.Elem {
	return []sequencer.Elem{
		SetMousePos{Pos: from},
		RightMouseButtonDown{},
		sequencer.Wait{Duration: humanizedMillis(60)},
		SetMousePos{Pos: to},
		sequencer.Wait{Duration: humanizedMillis(60)},
		RightMouseButtonUp{},
		sequencer.Wait{Duration: humanizedMillis(60)},
	}
}

func fromCurrent(x, y int32) mouse.Pos {
	pos := mouse.GetCursorPos()

	return mouse.Pos{
		X: pos.X + x,
		Y: pos.Y + y,
	}
}
func divingTankOff() []sequencer.Elem {
	return flatten(
		[][]sequencer.Elem{
			openInventory(),
			dragClothingToInventory(7, 19),
			closeInventory(),
		},
	)
}

func divingTankOn() []sequencer.Elem {
	return flatten(
		[][]sequencer.Elem{
			openInventory(),
			dragInventoryToClothing(19, 7),
			closeInventory(),
		},
	)
}

func smartBreath() []sequencer.Elem {
	var s = [][]sequencer.Elem{
		divingTankOn(),
		{
			sequencer.Wait{Duration: humanizedMillis(6000)},
		},
		divingTankOff(),
		{
			sequencer.Wait{Duration: humanizedMillis(2000)},
		},
		{
			sequencer.Loop{},
		},
	}

	return flatten(s)
}

func dragInventoryToClothing(inventorySlot, clothingSlot int) []sequencer.Elem {
	return dragInventory(InventorySlotPos(inventorySlot), clothingSlotPos(clothingSlot))
}

func dragClothingToInventory(clothingSlot, inventorySlot int) []sequencer.Elem {
	return dragInventory(clothingSlotPos(clothingSlot), InventorySlotPos(inventorySlot))
}

func toggleInventory() []sequencer.Elem {
	return []sequencer.Elem{
		KeyDown{Key: key.VK_TAB},
		KeyUp{Key: key.VK_TAB},
	}
}

func openInventory() []sequencer.Elem {
	return flatten([][]sequencer.Elem{
		toggleInventory(),
		{
			sequencer.Wait{Duration: humanizedMillis(inventoryOpenWaitShort)},
		},
	})
}

func closeInventory() []sequencer.Elem {
	return flatten([][]sequencer.Elem{
		toggleInventory(),
	})
}

func dragInventory(from mouse.Pos, to mouse.Pos) []sequencer.Elem {
	return []sequencer.Elem{
		SetMousePos{Pos: from},
		LeftMouseButtonDown{},
		sequencer.Wait{Duration: humanizedMillis(mouseDragWait)},
		SetMousePos{Pos: to},
		sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
		LeftMouseButtonUp{},
		sequencer.Wait{Duration: humanizedMillis(mouseDragWait)},
	}
}

func clickingLeft() []sequencer.Elem {
	return []sequencer.Elem{
		LeftMouseButtonDown{},
		LeftMouseButtonUp{},
		sequencer.Wait{Duration: humanizedMillis(300)},
		sequencer.Loop{},
	}
}

func runForward() []sequencer.Elem {
	return []sequencer.Elem{
		KeyDown{Key: key.VK_LSHIFT},
		KeyDown{Key: key.VK_W},
		sequencer.Wait{Duration: humanizedMillis(3000)},
		KeyUp{Key: key.VK_LSHIFT},
	}
}

func runForwardStop() []sequencer.Elem {
	return []sequencer.Elem{
		KeyUp{Key: key.VK_LSHIFT},
		KeyUp{Key: key.VK_W},
	}
}

func arm() []sequencer.Elem {
	return flatten(
		[][]sequencer.Elem{
			openInventory(),
			detachInventoryItem(1),
			detachInventoryItem(2),
			detachInventoryItem(3),
			detachInventoryItem(4),
			detachInventoryItem(5),
			detachInventoryItem(6),
			{sequencer.Wait{Duration: humanizedMillis(600)}},
			closeInventory(),
		},
	)
}

func unarm() []sequencer.Elem {
	return flatten(
		[][]sequencer.Elem{
			openInventory(),
			detachItem(1),
			detachItem(2),
			detachItem(3),
			detachItem(4),
			detachItem(5),
			detachItem(6),
			{sequencer.Wait{Duration: humanizedMillis(600)}},
			closeInventory(),
		},
	)
}

func detachInventoryItem(inventorySlot int) []sequencer.Elem {
	return []sequencer.Elem{
		SetMousePos{Pos: InventorySlotPos(inventorySlot)},
		sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
		RightMouseButtonDown{},
		sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
		RightMouseButtonUp{},
		sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},
	}
}

func detachItem(itemSlot int) []sequencer.Elem {
	return []sequencer.Elem{
		SetMousePos{Pos: ItemSlotPos(itemSlot)},
		sequencer.Wait{Duration: humanizedMillis(mouseMoveWait)},
		RightMouseButtonDown{},
		sequencer.Wait{Duration: humanizedMillis(mouseClickWait)},
		RightMouseButtonUp{},
		sequencer.Wait{Duration: humanizedMillis(inventoryMoveWait)},
	}
}

func duckGatherTree() []sequencer.Elem {
	return []sequencer.Elem{
		KeyDown{Key: key.VK_LCONTROL},
		LeftMouseButtonDown{},
		sequencer.Wait{Duration: humanizedMillis(7000)},
		LeftMouseButtonUp{},
		KeyUp{Key: key.VK_LCONTROL},
	}
}

func kayakBackward() []sequencer.Elem {
	return flatten(
		[][]sequencer.Elem{
			{KeyDown{Key: key.VK_S}},
			kayakPaddle(),
			{KeyUp{Key: key.VK_S}},
			{sequencer.Loop{}},
		},
	)
}

func kayakForward() []sequencer.Elem {
	return flatten(
		[][]sequencer.Elem{
			kayakPaddle(),
			{sequencer.Loop{}},
		},
	)
}

func kayakPaddle() []sequencer.Elem {
	return []sequencer.Elem{
		KeyDown{Key: key.VK_A},
		sequencer.Wait{Duration: humanizedMillis(paddleTime + 240)},
		KeyUp{Key: key.VK_A},
		sequencer.Wait{Duration: humanizedMillis(paddleWait + 240)},
		KeyDown{Key: key.VK_D},
		sequencer.Wait{Duration: humanizedMillis(paddleTime)},
		KeyUp{Key: key.VK_D},
		sequencer.Wait{Duration: humanizedMillis(paddleWait)},
	}
}

func getMousePos() []sequencer.Elem {
	return []sequencer.Elem{
		LookupMousePos{},
	}
}

func flatten(seqList [][]sequencer.Elem) []sequencer.Elem {
	var s []sequencer.Elem

	for _, elements := range seqList {
		s = append(s, elements...)
	}

	return s
}
