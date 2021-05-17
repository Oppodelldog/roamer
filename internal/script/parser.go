package script

import (
	"fmt"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
)

var commandMappings = map[string]func() interface{}{
	"L":  func() interface{} { return sequencer.Loop{} },
	"R":  func() interface{} { return sequencer.Repeat{} },
	"W":  func() interface{} { return sequencer.Wait{} },
	"KD": func() interface{} { return general.KeyDown{} },
	"KU": func() interface{} { return general.KeyUp{} },
	"LD": func() interface{} { return general.LeftMouseButtonDown{} },
	"LU": func() interface{} { return general.LeftMouseButtonUp{} },
	"MP": func() interface{} { return general.LookupMousePos{} },
	"MM": func() interface{} { return general.MouseMove{} },
	"RD": func() interface{} { return general.RightMouseButtonDown{} },
	"RU": func() interface{} { return general.RightMouseButtonUp{} },
	"SM": func() interface{} { return general.SetMousePos{} },
}

func Parse(script string) ([]sequencer.Elem, error) {
	var (
		t   *TokenStream
		seq []sequencer.Elem
		err error
	)

	t, err = lex(script)
	if err != nil {
		return nil, err
	}

	seq, err = parse(t)
	if err != nil {
		return nil, err
	}

	return seq, nil
}

func parse(t *TokenStream) ([]sequencer.Elem, error) {
	var seq []sequencer.Elem

	for !t.isEOF() {
		if t.Peek().Type == blockClose {
			t.Consume()

			return seq, nil
		}

		if t.Pos > 0 && t.Peek().Type == commandSeparator {
			t.Consume()

			continue
		}

		var (
			originalPos  = t.Pos
			command, err = parseCommand(t)
		)

		if err != nil {
			return nil, fmt.Errorf("[pos %v] error parsing command: %w", originalPos, err)
		}

		switch c := command.(type) {
		case sequencer.Repeat:
			for i := 0; i < c.Times; i++ {
				seq = append(seq, c.Sequence...)
			}
		default:
			seq = append(seq, command)
		}
	}

	return seq, nil
}

func parseCommand(t *TokenStream) (sequencer.Elem, error) {
	var token = t.Consume()
	if token.Type != literal {
		return nil, fmt.Errorf("expected command literal, but got '%s'", token.Type)
	}

	var commandElem, commandFound = commandMappings[token.Value]
	if !commandFound {
		return nil, fmt.Errorf("unknown command '%s'", token.Value)
	}

	elem, err := parseArguments(commandElem().(sequencer.Elem), t)
	if err != nil {
		return nil, fmt.Errorf("error parsing '%s': %w", token.Value, err)
	}

	return elem, nil
}

func parseArguments(elem sequencer.Elem, t *TokenStream) (sequencer.Elem, error) {
	var err error

	switch v := elem.(type) {
	case sequencer.Wait:
		return parseWait(v, t)
	case general.KeyDown:
		return parseKey(v, t)
	case general.KeyUp:
		return parseKey(v, t)
	case sequencer.Repeat:
		return parseRepeat(v, t)
	case general.MouseMove:
		return parseMouseMove(v, t)
	case general.SetMousePos:
		return parseSetMousePos(v, t)
	case general.LookupMousePos:
		return v, err
	case general.LeftMouseButtonDown:
		return v, err
	case general.LeftMouseButtonUp:
		return v, err
	case general.RightMouseButtonDown:
		return v, err
	case general.RightMouseButtonUp:
		return v, err
	case sequencer.Loop:
		return v, err
	}

	return elem, fmt.Errorf("unknown command %T", elem)
}
