package remnant

import (
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
)

const keyDownDelay = 20
const keyUpDelay = 20

func Hunter() []sequencer.Elem {
	return general.Flatten([][]sequencer.Elem{
		tab(),
		equip(HandGun, HGRepeaterPistol),
		equip(LongGun, LGHuntingRifle),
		equip(Melee, MWGuardianAxe),
		equip(ArmorHead, ArmHunter),
		equip(ArmorChest, ArmHunter),
		equip(ArmorFeet, ArmHunter),
		equip(Amulet, AmuCharcoalNecklace),
		tab(),
	})
}

func Radiant() []sequencer.Elem {
	return general.Flatten([][]sequencer.Elem{
		tab(),
		equip(HandGun, HGRepeaterPistol),
		equip(LongGun, LGFusionRifle),
		equip(Melee, MWGuardianAxe),
		equip(ArmorHead, ArmRadiant),
		equip(ArmorChest, ArmRadiant),
		equip(ArmorFeet, ArmRadiant),
		equip(Amulet, AmuVengeanceIdol),
		tab(),
	})
}

func multi(num int, seq []sequencer.Elem) []sequencer.Elem {
	var mewSeq []sequencer.Elem
	for i := 0; i < num; i++ {
		mewSeq = append(mewSeq, seq...)
	}

	return mewSeq
}

func unequip() []sequencer.Elem {
	return multi(2, space())
}

type ItemType int

const (
	HandGun ItemType = iota
	LongGun
	Melee
	ArmorHead
	ArmorChest
	ArmorFeet
	Amulet
)

func equip(itemType ItemType, id string) []sequencer.Elem {

	var keys []sequencer.Elem
	var slots []string
	var postKeys []sequencer.Elem

	switch itemType {
	case HandGun:
		keys = append(keys, unequip()...)
		slots = HGSlots
	case LongGun:
		keys = append(keys, multi(1, down())...)
		keys = append(keys, unequip()...)
		slots = LGSlots
		postKeys = multi(1, up())

	case Melee:
		keys = append(keys, multi(2, down())...)
		keys = append(keys, unequip()...)
		slots = MSlots
		postKeys = multi(2, up())

	case ArmorHead:
		keys = append(keys, multi(3, down())...)
		keys = append(keys, unequip()...)
		slots = ArmHeadSlots
		postKeys = multi(3, up())
	case ArmorChest:
		keys = append(keys, multi(3, down())...)
		keys = append(keys, multi(1, right())...)
		keys = append(keys, unequip()...)
		slots = ArmChestSlots
		postKeys = multi(3, up())
	case ArmorFeet:
		keys = append(keys, multi(3, down())...)
		keys = append(keys, multi(2, right())...)
		keys = append(keys, unequip()...)
		slots = ArmFeetSlots
		postKeys = multi(3, up())

	case Amulet:
		keys = append(keys, multi(4, down())...)
		keys = append(keys, unequip()...)
		slots = AmuSlots
		postKeys = multi(4, up())
	}

	var x = getSlotSteps(slots, id)
	keys = append(keys, multi(x, down())...)

	keys = append(keys, space()...)
	keys = append(keys, escape()...)
	keys = append(keys, postKeys...)

	return keys
}

func tab() []sequencer.Elem {
	return pressKey(key.VK_TAB)
}

func down() []sequencer.Elem {
	return pressKey(key.VK_DOWN)
}
func up() []sequencer.Elem {
	return pressKey(key.VK_UP)
}
func left() []sequencer.Elem {
	return pressKey(key.VK_LEFT)
}
func right() []sequencer.Elem {
	return pressKey(key.VK_RIGHT)
}
func space() []sequencer.Elem {
	return pressKey(key.VK_SPACE)
}
func escape() []sequencer.Elem {
	return pressKey(key.VK_ESC)
}
func pressKey(key int) []sequencer.Elem {
	return []sequencer.Elem{
		general.KeyDown{Key: key},
		sequencer.Wait{Duration: general.HumanizedMillis(keyDownDelay)},
		general.KeyUp{Key: key},
		sequencer.Wait{Duration: general.HumanizedMillis(keyUpDelay)},
	}
}

func getSlotSteps(slots []string, id string) int {
	for i, slot := range slots {
		if slot == id {
			return i
		}
	}

	return -1
}
