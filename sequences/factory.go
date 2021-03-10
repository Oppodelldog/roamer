package sequences

import "rust-roamer/sequencer"

func NewSequenceFunc(name string) func() []sequencer.Elem {
	return func() []sequencer.Elem {
		var sequences = map[string][]sequencer.Elem{
			"none":                    {},
			"get-mouse-pos":           getMousePos(),
			"click-left":              clickingLeft(),
			"run-forward":             runForward(),
			"run-forward-stop":        runForwardStop(),
			"kayak-forward":           kayakForward(),
			"kayak-backward":          kayakBackward(),
			"duck-gather-tree":        duckGatherTree(),
			"unarm":                   unarm(),
			"arm":                     arm(),
			"diving-tank-on":          divingTankOn(),
			"diving-tank-off":         divingTankOff(),
			"smart-breath":            smartBreath(),
			"fill-inventory-row":      fillInventoryRow(),
			"transfer-inventory-row":  transferInventoryRow(),
			"move-down-inventory-row": moveInventoryRow(1),
			"move-up-inventory-row":   moveInventoryRow(-1),
		}

		return sequences[name]
	}
}
