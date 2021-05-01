package rust

import (
	config2 "github.com/Oppodelldog/roamer/internal/config"
	key2 "github.com/Oppodelldog/roamer/internal/key"
	mouse2 "github.com/Oppodelldog/roamer/internal/mouse"
	sequencer2 "github.com/Oppodelldog/roamer/internal/sequencer"
	general2 "github.com/Oppodelldog/roamer/internal/sequences/general"
	version2 "github.com/Oppodelldog/roamer/internal/sequences/version"
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

func MoveInventoryRow(verticalFieldOffset int32) []sequencer2.Elem {
	pos := mouse2.GetCursorPos()
	var seqs [][]sequencer2.Elem

	for i := int32(0); i < 6; i++ {
		from := relativeHorizontalFieldPos(pos, i)
		to := relativeVerticalFieldPos(from, verticalFieldOffset)
		seqs = append(seqs, [][]sequencer2.Elem{Drag(from, to),
			{sequencer2.Wait{Duration: general2.HumanizedMillis(80)}}}...)
	}

	return general2.Flatten(seqs)
}

func TransferInventoryRow() []sequencer2.Elem {
	pos := mouse2.GetCursorPos()
	var seqs [][]sequencer2.Elem

	for i := int32(0); i < 6; i++ {
		seqs = append(seqs, [][]sequencer2.Elem{
			collect(relativeHorizontalFieldPos(pos, i)),
			{sequencer2.Wait{Duration: general2.HumanizedMillis(80)}}}...)
	}

	return general2.Flatten(seqs)
}

func relativeHorizontalFieldPos(from mouse2.Pos, fieldNo int32) mouse2.Pos {
	return mouse2.Pos{X: from.X + (fieldNo * (itemFieldSize + itemFieldMargin)), Y: from.Y}
}

func relativeVerticalFieldPos(from mouse2.Pos, fieldNo int32) mouse2.Pos {
	return mouse2.Pos{X: from.X, Y: from.Y + (fieldNo * (itemFieldSize + itemFieldMargin))}
}

func FillInventoryRow() []sequencer2.Elem {
	pos := mouse2.GetCursorPos()

	var seqs [][]sequencer2.Elem

	for i := int32(1); i < 6; i++ {
		seqs = append(seqs, [][]sequencer2.Elem{
			Unstack(pos, relativeHorizontalFieldPos(pos, i)),
			{sequencer2.Wait{Duration: general2.HumanizedMillis(80)}}}...)
	}

	return general2.Flatten(seqs)
}

func collect(from mouse2.Pos) []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.SetMousePos{Pos: from},
		general2.RightMouseButtonDown{},
		sequencer2.Wait{Duration: general2.HumanizedMillis(80)},
		general2.RightMouseButtonUp{},
		sequencer2.Wait{Duration: general2.HumanizedMillis(80)},
	}
}

func Unstack(from, to mouse2.Pos) []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.SetMousePos{Pos: from},
		general2.RightMouseButtonDown{},
		sequencer2.Wait{Duration: general2.HumanizedMillis(80)},
		general2.SetMousePos{Pos: to},
		sequencer2.Wait{Duration: general2.HumanizedMillis(80)},
		general2.RightMouseButtonUp{},
		sequencer2.Wait{Duration: general2.HumanizedMillis(80)},
	}
}

func Drag(from, to mouse2.Pos) []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.SetMousePos{Pos: from},
		general2.LeftMouseButtonDown{},
		sequencer2.Wait{Duration: general2.HumanizedMillis(80)},
		general2.SetMousePos{Pos: to},
		sequencer2.Wait{Duration: general2.HumanizedMillis(80)},
		general2.LeftMouseButtonUp{},
		sequencer2.Wait{Duration: general2.HumanizedMillis(80)},
	}
}

func DivingTankOff() []sequencer2.Elem {
	return general2.Flatten(
		[][]sequencer2.Elem{
			openInventory(),
			dragClothingToInventory(7, 19),
			closeInventory(),
		},
	)
}

func DivingTankOn() []sequencer2.Elem {
	return general2.Flatten(
		[][]sequencer2.Elem{
			openInventory(),
			dragInventoryToClothing(19, 7),
			closeInventory(),
		},
	)
}

func SmartBreath() []sequencer2.Elem {
	var s = [][]sequencer2.Elem{
		DivingTankOn(),
		{
			sequencer2.Wait{Duration: general2.HumanizedMillis(6000)},
		},
		DivingTankOff(),
		{
			sequencer2.Wait{Duration: general2.HumanizedMillis(2000)},
		},
		{
			sequencer2.Loop{},
		},
	}

	return general2.Flatten(s)
}

func dragInventoryToClothing(inventorySlot, clothingSlot int) []sequencer2.Elem {
	var inventorySlots = getSlots(slotsInventory)
	var clothingSlots = getSlots(slotsInventory)
	return dragInventory(
		mousePos(inventorySlots.At(inventorySlot)),
		mousePos(clothingSlots.At(clothingSlot)),
	)
}

func dragClothingToInventory(clothingSlot, inventorySlot int) []sequencer2.Elem {
	var inventorySlots = getSlots(slotsInventory)
	var clothingSlots = getSlots(slotsInventory)
	return dragInventory(
		mousePos(clothingSlots.At(clothingSlot)),
		mousePos(inventorySlots.At(inventorySlot)),
	)
}

func toggleInventory() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyDown{Key: key2.VK_TAB},
		general2.KeyUp{Key: key2.VK_TAB},
	}
}

func openInventory() []sequencer2.Elem {
	return general2.Flatten([][]sequencer2.Elem{
		toggleInventory(),
		{
			sequencer2.Wait{Duration: general2.HumanizedMillis(inventoryOpenWaitShort)},
		},
	})
}

func closeInventory() []sequencer2.Elem {
	return general2.Flatten([][]sequencer2.Elem{
		toggleInventory(),
	})
}

func dragInventory(from mouse2.Pos, to mouse2.Pos) []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.SetMousePos{Pos: from},
		general2.LeftMouseButtonDown{},
		sequencer2.Wait{Duration: general2.HumanizedMillis(mouseDragWait)},
		general2.SetMousePos{Pos: to},
		sequencer2.Wait{Duration: general2.HumanizedMillis(mouseMoveWait)},
		general2.LeftMouseButtonUp{},
		sequencer2.Wait{Duration: general2.HumanizedMillis(mouseDragWait)},
	}
}

func ClickingLeft() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.LeftMouseButtonDown{},
		general2.LeftMouseButtonUp{},
		sequencer2.Wait{Duration: general2.HumanizedMillis(300)},
		sequencer2.Loop{},
	}
}

func RunForward() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyDown{Key: key2.VK_LSHIFT},
		general2.KeyDown{Key: key2.VK_W},
		sequencer2.Wait{Duration: general2.HumanizedMillis(3000)},
		general2.KeyUp{Key: key2.VK_LSHIFT},
	}
}

func RunForwardStop() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyUp{Key: key2.VK_LSHIFT},
		general2.KeyUp{Key: key2.VK_W},
	}
}

func Arm() []sequencer2.Elem {
	var inventorySlots = getSlots(slotsInventory)
	var beltSlots = getSlots(slotsBelt)

	return general2.Flatten(
		[][]sequencer2.Elem{
			openInventory(),
			Drag(mousePos(inventorySlots.At(19)), mousePos(beltSlots.At(1))),
			Drag(mousePos(inventorySlots.At(20)), mousePos(beltSlots.At(2))),
			Drag(mousePos(inventorySlots.At(21)), mousePos(beltSlots.At(3))),
			Drag(mousePos(inventorySlots.At(22)), mousePos(beltSlots.At(4))),
			Drag(mousePos(inventorySlots.At(23)), mousePos(beltSlots.At(5))),
			Drag(mousePos(inventorySlots.At(24)), mousePos(beltSlots.At(6))),
			{sequencer2.Wait{Duration: general2.HumanizedMillis(300)}},
			closeInventory(),
		},
	)
}

func Unarm() []sequencer2.Elem {
	var inventorySlots = getSlots(slotsInventory)
	var beltSlots = getSlots(slotsBelt)

	return general2.Flatten(
		[][]sequencer2.Elem{
			openInventory(),
			Drag(mousePos(beltSlots.At(1)), mousePos(inventorySlots.At(19))),
			Drag(mousePos(beltSlots.At(2)), mousePos(inventorySlots.At(20))),
			Drag(mousePos(beltSlots.At(3)), mousePos(inventorySlots.At(21))),
			Drag(mousePos(beltSlots.At(4)), mousePos(inventorySlots.At(22))),
			Drag(mousePos(beltSlots.At(5)), mousePos(inventorySlots.At(23))),
			Drag(mousePos(beltSlots.At(6)), mousePos(inventorySlots.At(24))),
			{sequencer2.Wait{Duration: general2.HumanizedMillis(300)}},
			closeInventory(),
		},
	)
}

func DuckGatherTree() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyDown{Key: key2.VK_LCONTROL},
		general2.LeftMouseButtonDown{},
		sequencer2.Wait{Duration: general2.HumanizedMillis(7000)},
		general2.LeftMouseButtonUp{},
		general2.KeyUp{Key: key2.VK_LCONTROL},
	}
}

func KayakBackward() []sequencer2.Elem {
	return general2.Flatten(
		[][]sequencer2.Elem{
			{general2.KeyDown{Key: key2.VK_S}},
			KayakPaddle(),
			{general2.KeyUp{Key: key2.VK_S}},
			{sequencer2.Loop{}},
		},
	)
}

func KayakForward() []sequencer2.Elem {
	return general2.Flatten(
		[][]sequencer2.Elem{
			KayakPaddle(),
			{sequencer2.Loop{}},
		},
	)
}

func KayakPaddle() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.KeyDown{Key: key2.VK_A},
		sequencer2.Wait{Duration: general2.HumanizedMillis(paddleTime + 240)},
		general2.KeyUp{Key: key2.VK_A},
		sequencer2.Wait{Duration: general2.HumanizedMillis(paddleWait + 240)},
		general2.KeyDown{Key: key2.VK_D},
		sequencer2.Wait{Duration: general2.HumanizedMillis(paddleTime)},
		general2.KeyUp{Key: key2.VK_D},
		sequencer2.Wait{Duration: general2.HumanizedMillis(paddleWait)},
	}
}

func GetMousePos() []sequencer2.Elem {
	return []sequencer2.Elem{
		general2.LookupMousePos{},
	}
}

func getSlots(kind string) config2.Slots {
	return config2.GetSlots(version2.Get(), game, kind)
}

func mousePos(pos config2.Pos) mouse2.Pos {
	return mouse2.Pos{
		X: int32(pos.X),
		Y: int32(pos.Y),
	}
}
