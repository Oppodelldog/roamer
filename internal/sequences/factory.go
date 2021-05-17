package sequences

import (
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/altf4"
	"github.com/Oppodelldog/roamer/internal/sequences/remnant"
	"github.com/Oppodelldog/roamer/internal/sequences/rust"
	"github.com/Oppodelldog/roamer/internal/sequences/sevendaystodie"
)

func NewBuildInSequenceFunc(name string) func() []sequencer.Elem {
	return func() []sequencer.Elem {
		var sequences = map[string][]sequencer.Elem{
			"none":                    {},
			"click-left":              rust.ClickingLeft(),
			"run-forward":             rust.RunForward(),
			"run-forward-stop":        rust.RunForwardStop(),
			"kayak-forward":           rust.KayakForward(),
			"kayak-backward":          rust.KayakBackward(),
			"duck-gather-tree":        rust.DuckGatherTree(),
			"unarm":                   rust.Unarm(),
			"arm":                     rust.Arm(),
			"diving-tank-on":          rust.DivingTankOn(),
			"diving-tank-off":         rust.DivingTankOff(),
			"smart-breath":            rust.SmartBreath(),
			"fill-inventory-row":      rust.FillInventoryRow(),
			"transfer-inventory-row":  rust.TransferInventoryRow(),
			"move-down-inventory-row": rust.MoveInventoryRow(1),
			"move-up-inventory-row":   rust.MoveInventoryRow(-1),
			"7d2d_repair-1":           sevendaystodie.Repair(1),
			"7d2d_repair-2":           sevendaystodie.Repair(2),
			"7d2d_repair-3":           sevendaystodie.Repair(3),
			"7d2d_repair-4":           sevendaystodie.Repair(4),
			"7d2d_repair-5":           sevendaystodie.Repair(5),
			"7d2d_repair-6":           sevendaystodie.Repair(6),
			"7d2d_repair-7":           sevendaystodie.Repair(7),
			"7d2d_repair-8":           sevendaystodie.Repair(8),
			"7d2d_repair-9":           sevendaystodie.Repair(9),
			"7d2d_repair-0":           sevendaystodie.Repair(0),
			"7d2d_walk":               sevendaystodie.Walk(),
			"7d2d_run":                sevendaystodie.Run(),
			"7d2d_walk_run_stop":      sevendaystodie.WalkRunStop(),
			"7d2d_click_left_fast":    sevendaystodie.ClickingLeft(30),
			"7d2d_click_left":         sevendaystodie.ClickingLeft(200),
			"7d2d_click_left_slow":    sevendaystodie.ClickingLeft(1000),
			"altf4-run-all":           altf4.RunAll(),
			"altf4-run-all-2":         altf4.RunAll2(),
			"altf4-run-0":             altf4.Reset(),
			"altf4-run-1":             altf4.Run1(),
			"altf4-run-2":             altf4.Run2(),
			"altf4-run-3":             altf4.Run3(),
			"altf4-run-4":             altf4.Run4(),
			"altf4-run-5":             altf4.Run5(),
			"remnant-equip-hunter":    remnant.Hunter(),
			"remnant-equip-radiant":   remnant.Radiant(),
		}

		return sequences[name]
	}
}
