package script

import (
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
)

const valBlockClose = ']'
const valBlockOpen = '['
const valCommandSep = ';'
const valArgSep = ' '

var commandMappings = map[string]func() interface{}{
	"NOP": func() interface{} { return sequencer.NoOperation{} },
	"L":   func() interface{} { return sequencer.Loop{} },
	"R":   func() interface{} { return sequencer.Repeat{} },
	"W":   func() interface{} { return sequencer.Wait{} },
	"KD":  func() interface{} { return general.KeyDown{} },
	"KU":  func() interface{} { return general.KeyUp{} },
	"LD":  func() interface{} { return general.LeftMouseButtonDown{} },
	"LU":  func() interface{} { return general.LeftMouseButtonUp{} },
	"MP":  func() interface{} { return general.LookupMousePos{} },
	"MM":  func() interface{} { return general.MouseMove{} },
	"RD":  func() interface{} { return general.RightMouseButtonDown{} },
	"RU":  func() interface{} { return general.RightMouseButtonUp{} },
	"SM":  func() interface{} { return general.SetMousePos{} },
}
