package sequences

import (
	"rust-roamer/sequencer"
	"rust-roamer/sequences/rust"
	"rust-roamer/sequences/sevendaystodie"
)

func NewSequenceFunc(name string) func() []sequencer.Elem {
	return func() []sequencer.Elem {
		var sequences = map[string][]sequencer.Elem{
			"none":                    {},
			"get-mouse-pos":           rust.GetMousePos(),
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
		}

		return sequences[name]
	}
}
