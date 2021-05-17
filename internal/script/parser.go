package script

import (
	"fmt"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
)

var commandMappings = map[string]func() interface{}{
	"W":  func() interface{} { return sequencer.Wait{} },
	"L":  func() interface{} { return sequencer.Loop{} },
	"KD": func() interface{} { return general.KeyDown{} },
	"KU": func() interface{} { return general.KeyUp{} },
	"LD": func() interface{} { return general.LeftMouseButtonDown{} },
	"LU": func() interface{} { return general.LeftMouseButtonUp{} },
	"RD": func() interface{} { return general.RightMouseButtonDown{} },
	"RU": func() interface{} { return general.RightMouseButtonUp{} },
	"MM": func() interface{} { return general.MouseMove{} },
	"SM": func() interface{} { return general.SetMousePos{} },
	"R":  func() interface{} { return sequencer.Repeat{} },
}

func NewCustomSequenceFunc(script string) func() []sequencer.Elem {
	return func() []sequencer.Elem {
		seq, err := Parse(script)
		if err != nil {
			panic(err)
		}

		return seq
	}
}

func Parse(script string) ([]sequencer.Elem, error) {
	t, err := lex(script)
	if err != nil {
		return nil, err
	}

	seq, err := parse(t)
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

		if t.Pos > 0 {
			sep := t.Peek()
			if sep.Type == commandSeparator {
				t.Consume()

				continue
			}
			prev := t.PeekAt(-1)
			if sep.Type != commandSeparator && prev.Type != blockOpen && prev.Type != commandSeparator {
				return nil, fmt.Errorf("[pos %v] expected command separator, but got: '%s'(='%s')", t.Pos, sep.Type, sep.Value)
			}
		}

		originalPos := t.Pos
		command, err := parseCommand(t)
		if err != nil {
			return nil, fmt.Errorf("[pos %v] error parsing command: %w", originalPos, err)
		}

		if repeat, isRepeat := command.(sequencer.Repeat); isRepeat {
			for i := 0; i < repeat.Times; i++ {
				seq = append(seq, repeat.Sequence...)
			}
		} else {
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

	return elem, fmt.Errorf("unkown command %T", elem)
}
