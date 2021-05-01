package sequences

import (
	sequencer2 "github.com/Oppodelldog/roamer/internal/sequencer"
	altf42 "github.com/Oppodelldog/roamer/internal/sequences/altf4"
	rust2 "github.com/Oppodelldog/roamer/internal/sequences/rust"
	sevendaystodie2 "github.com/Oppodelldog/roamer/internal/sequences/sevendaystodie"
	vallheim2 "github.com/Oppodelldog/roamer/internal/sequences/vallheim"
)

func NewSequenceFunc(name string) func() []sequencer2.Elem {
	return func() []sequencer2.Elem {
		var sequences = map[string][]sequencer2.Elem{
			"none":                    {},
			"get-mouse-pos":           rust2.GetMousePos(),
			"click-left":              rust2.ClickingLeft(),
			"run-forward":             rust2.RunForward(),
			"run-forward-stop":        rust2.RunForwardStop(),
			"kayak-forward":           rust2.KayakForward(),
			"kayak-backward":          rust2.KayakBackward(),
			"duck-gather-tree":        rust2.DuckGatherTree(),
			"unarm":                   rust2.Unarm(),
			"arm":                     rust2.Arm(),
			"diving-tank-on":          rust2.DivingTankOn(),
			"diving-tank-off":         rust2.DivingTankOff(),
			"smart-breath":            rust2.SmartBreath(),
			"fill-inventory-row":      rust2.FillInventoryRow(),
			"transfer-inventory-row":  rust2.TransferInventoryRow(),
			"move-down-inventory-row": rust2.MoveInventoryRow(1),
			"move-up-inventory-row":   rust2.MoveInventoryRow(-1),
			"7d2d_repair-1":           sevendaystodie2.Repair(1),
			"7d2d_repair-2":           sevendaystodie2.Repair(2),
			"7d2d_repair-3":           sevendaystodie2.Repair(3),
			"7d2d_repair-4":           sevendaystodie2.Repair(4),
			"7d2d_repair-5":           sevendaystodie2.Repair(5),
			"7d2d_repair-6":           sevendaystodie2.Repair(6),
			"7d2d_repair-7":           sevendaystodie2.Repair(7),
			"7d2d_repair-8":           sevendaystodie2.Repair(8),
			"7d2d_repair-9":           sevendaystodie2.Repair(9),
			"7d2d_repair-0":           sevendaystodie2.Repair(0),
			"7d2d_walk":               sevendaystodie2.Walk(),
			"7d2d_run":                sevendaystodie2.Run(),
			"7d2d_walk_run_stop":      sevendaystodie2.WalkRunStop(),
			"7d2d_click_left_fast":    sevendaystodie2.ClickingLeft(30),
			"7d2d_click_left":         sevendaystodie2.ClickingLeft(200),
			"7d2d_click_left_slow":    sevendaystodie2.ClickingLeft(1000),
			"altf4-run-all":           altf42.RunAll(),
			"altf4-run-all-2":         altf42.RunAll2(),
			"altf4-run-0":             altf42.Reset(),
			"altf4-run-1":             altf42.Run1(),
			"altf4-run-2":             altf42.Run2(),
			"altf4-run-3":             altf42.Run3(),
			"altf4-run-4":             altf42.Run4(),
			"altf4-run-5":             altf42.Run5(),
			"vallheim-run":            vallheim2.Run(),
			"vallheim-walk":           vallheim2.Walk(),
			"vallheim-grillmaster":    vallheim2.Grillmaster(),
		}

		return sequences[name]
	}
}
