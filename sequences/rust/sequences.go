package rust

import (
	"rust-roamer/config"
	"rust-roamer/key"
	"rust-roamer/mouse"
	"rust-roamer/sequencer"
	"rust-roamer/sequences/general"
	"rust-roamer/sequences/version"
)

const game = "rust"
const slotsBelt = "belt"
const slotsInventory = "inventory"
const slotsClothing = "clothing"

const paddleTime = 1200
const paddleWait = 400
const inventoryOpenWaitShort = 200
const inventoryMoveWait = 60
const mouseMoveWait = 60
const mouseClickWait = 0
const mouseDragWait = 16
const itemFieldMargin = 3
const itemFieldSize = 92

func MoveInventoryRow(verticalFieldOffset int32) []sequencer.Elem {
	pos := mouse.GetCursorPos()
	var seqs [][]sequencer.Elem

	for i := int32(0); i < 6; i++ {
		from := relativeHorizontalFieldPos(pos, i)
		to := relativeVerticalFieldPos(from, verticalFieldOffset)
		seqs = append(seqs, [][]sequencer.Elem{Drag(from, to),
			{sequencer.Wait{Duration: general.HumanizedMillis(80)}}}...)
	}

	return general.Flatten(seqs)
}

func TransferInventoryRow() []sequencer.Elem {
	pos := mouse.GetCursorPos()
	var seqs [][]sequencer.Elem

	for i := int32(0); i < 6; i++ {
		seqs = append(seqs, [][]sequencer.Elem{
			collect(relativeHorizontalFieldPos(pos, i)),
			{sequencer.Wait{Duration: general.HumanizedMillis(80)}}}...)
	}

	return general.Flatten(seqs)
}

func relativeHorizontalFieldPos(from mouse.Pos, fieldNo int32) mouse.Pos {
	return mouse.Pos{X: from.X + (fieldNo * (itemFieldSize + itemFieldMargin)), Y: from.Y}
}

func relativeVerticalFieldPos(from mouse.Pos, fieldNo int32) mouse.Pos {
	return mouse.Pos{X: from.X, Y: from.Y + (fieldNo * (itemFieldSize + itemFieldMargin))}
}

func FillInventoryRow() []sequencer.Elem {
	pos := mouse.GetCursorPos()

	var seqs [][]sequencer.Elem

	for i := int32(1); i < 6; i++ {
		seqs = append(seqs, [][]sequencer.Elem{
			Unstack(pos, relativeHorizontalFieldPos(pos, i)),
			{sequencer.Wait{Duration: general.HumanizedMillis(80)}}}...)
	}

	return general.Flatten(seqs)
}

func collect(from mouse.Pos) []sequencer.Elem {
	return []sequencer.Elem{
		general.SetMousePos{Pos: from},
		general.RightMouseButtonDown{},
		sequencer.Wait{Duration: general.HumanizedMillis(80)},
		general.RightMouseButtonUp{},
		sequencer.Wait{Duration: general.HumanizedMillis(80)},
	}
}

func Unstack(from, to mouse.Pos) []sequencer.Elem {
	return []sequencer.Elem{
		general.SetMousePos{Pos: from},
		general.RightMouseButtonDown{},
		sequencer.Wait{Duration: general.HumanizedMillis(80)},
		general.SetMousePos{Pos: to},
		sequencer.Wait{Duration: general.HumanizedMillis(80)},
		general.RightMouseButtonUp{},
		sequencer.Wait{Duration: general.HumanizedMillis(80)},
	}
}

func Drag(from, to mouse.Pos) []sequencer.Elem {
	return []sequencer.Elem{
		general.SetMousePos{Pos: from},
		general.LeftMouseButtonDown{},
		sequencer.Wait{Duration: general.HumanizedMillis(80)},
		general.SetMousePos{Pos: to},
		sequencer.Wait{Duration: general.HumanizedMillis(80)},
		general.LeftMouseButtonUp{},
		sequencer.Wait{Duration: general.HumanizedMillis(80)},
	}
}

func DivingTankOff() []sequencer.Elem {
	return general.Flatten(
		[][]sequencer.Elem{
			openInventory(),
			dragClothingToInventory(7, 19),
			closeInventory(),
		},
	)
}

func DivingTankOn() []sequencer.Elem {
	return general.Flatten(
		[][]sequencer.Elem{
			openInventory(),
			dragInventoryToClothing(19, 7),
			closeInventory(),
		},
	)
}

func SmartBreath() []sequencer.Elem {
	var s = [][]sequencer.Elem{
		DivingTankOn(),
		{
			sequencer.Wait{Duration: general.HumanizedMillis(6000)},
		},
		DivingTankOff(),
		{
			sequencer.Wait{Duration: general.HumanizedMillis(2000)},
		},
		{
			sequencer.Loop{},
		},
	}

	return general.Flatten(s)
}

func dragInventoryToClothing(inventorySlot, clothingSlot int) []sequencer.Elem {
	var inventorySlots = getSlots(slotsInventory)
	var clothingSlots = getSlots(slotsInventory)
	return dragInventory(
		mousePos(inventorySlots.At(inventorySlot)),
		mousePos(clothingSlots.At(clothingSlot)),
	)
}

func dragClothingToInventory(clothingSlot, inventorySlot int) []sequencer.Elem {
	var inventorySlots = getSlots(slotsInventory)
	var clothingSlots = getSlots(slotsInventory)
	return dragInventory(
		mousePos(clothingSlots.At(clothingSlot)),
		mousePos(inventorySlots.At(inventorySlot)),
	)
}

func toggleInventory() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{Key: key.VK_TAB},
		general.KeyUp{Key: key.VK_TAB},
	}
}

func openInventory() []sequencer.Elem {
	return general.Flatten([][]sequencer.Elem{
		toggleInventory(),
		{
			sequencer.Wait{Duration: general.HumanizedMillis(inventoryOpenWaitShort)},
		},
	})
}

func closeInventory() []sequencer.Elem {
	return general.Flatten([][]sequencer.Elem{
		toggleInventory(),
	})
}

func dragInventory(from mouse.Pos, to mouse.Pos) []sequencer.Elem {
	return []sequencer.Elem{
		general.SetMousePos{Pos: from},
		general.LeftMouseButtonDown{},
		sequencer.Wait{Duration: general.HumanizedMillis(mouseDragWait)},
		general.SetMousePos{Pos: to},
		sequencer.Wait{Duration: general.HumanizedMillis(mouseMoveWait)},
		general.LeftMouseButtonUp{},
		sequencer.Wait{Duration: general.HumanizedMillis(mouseDragWait)},
	}
}

func ClickingLeft() []sequencer.Elem {
	return []sequencer.Elem{
		general.LeftMouseButtonDown{},
		general.LeftMouseButtonUp{},
		sequencer.Wait{Duration: general.HumanizedMillis(300)},
		sequencer.Loop{},
	}
}

func RunForward() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{Key: key.VK_LSHIFT},
		general.KeyDown{Key: key.VK_W},
		sequencer.Wait{Duration: general.HumanizedMillis(3000)},
		general.KeyUp{Key: key.VK_LSHIFT},
	}
}

func RunForwardStop() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyUp{Key: key.VK_LSHIFT},
		general.KeyUp{Key: key.VK_W},
	}
}

func Arm() []sequencer.Elem {
	var inventorySlots = getSlots(slotsInventory)
	var beltSlots = getSlots(slotsBelt)

	return general.Flatten(
		[][]sequencer.Elem{
			openInventory(),
			Drag(mousePos(inventorySlots.At(19)), mousePos(beltSlots.At(1))),
			Drag(mousePos(inventorySlots.At(20)), mousePos(beltSlots.At(2))),
			Drag(mousePos(inventorySlots.At(21)), mousePos(beltSlots.At(3))),
			Drag(mousePos(inventorySlots.At(22)), mousePos(beltSlots.At(4))),
			Drag(mousePos(inventorySlots.At(23)), mousePos(beltSlots.At(5))),
			Drag(mousePos(inventorySlots.At(24)), mousePos(beltSlots.At(6))),
			{sequencer.Wait{Duration: general.HumanizedMillis(300)}},
			closeInventory(),
		},
	)
}

func Unarm() []sequencer.Elem {
	var inventorySlots = getSlots(slotsInventory)
	var beltSlots = getSlots(slotsBelt)

	return general.Flatten(
		[][]sequencer.Elem{
			openInventory(),
			Drag(mousePos(beltSlots.At(1)), mousePos(inventorySlots.At(19))),
			Drag(mousePos(beltSlots.At(2)), mousePos(inventorySlots.At(20))),
			Drag(mousePos(beltSlots.At(3)), mousePos(inventorySlots.At(21))),
			Drag(mousePos(beltSlots.At(4)), mousePos(inventorySlots.At(22))),
			Drag(mousePos(beltSlots.At(5)), mousePos(inventorySlots.At(23))),
			Drag(mousePos(beltSlots.At(6)), mousePos(inventorySlots.At(24))),
			{sequencer.Wait{Duration: general.HumanizedMillis(300)}},
			closeInventory(),
		},
	)
}

func DuckGatherTree() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{Key: key.VK_LCONTROL},
		general.LeftMouseButtonDown{},
		sequencer.Wait{Duration: general.HumanizedMillis(7000)},
		general.LeftMouseButtonUp{},
		general.KeyUp{Key: key.VK_LCONTROL},
	}
}

func KayakBackward() []sequencer.Elem {
	return general.Flatten(
		[][]sequencer.Elem{
			{general.KeyDown{Key: key.VK_S}},
			KayakPaddle(),
			{general.KeyUp{Key: key.VK_S}},
			{sequencer.Loop{}},
		},
	)
}

func KayakForward() []sequencer.Elem {
	return general.Flatten(
		[][]sequencer.Elem{
			KayakPaddle(),
			{sequencer.Loop{}},
		},
	)
}

func KayakPaddle() []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{Key: key.VK_A},
		sequencer.Wait{Duration: general.HumanizedMillis(paddleTime + 240)},
		general.KeyUp{Key: key.VK_A},
		sequencer.Wait{Duration: general.HumanizedMillis(paddleWait + 240)},
		general.KeyDown{Key: key.VK_D},
		sequencer.Wait{Duration: general.HumanizedMillis(paddleTime)},
		general.KeyUp{Key: key.VK_D},
		sequencer.Wait{Duration: general.HumanizedMillis(paddleWait)},
	}
}

func GetMousePos() []sequencer.Elem {
	return []sequencer.Elem{
		general.LookupMousePos{},
	}
}

func getSlots(kind string) config.Slots {
	return config.GetSlots(version.Get(), game, kind)
}

func mousePos(pos config.Pos) mouse.Pos {
	return mouse.Pos{
		X: int32(pos.X),
		Y: int32(pos.Y),
	}
}
