package script

import (
	"errors"
	"fmt"

	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
)

var ErrUnknownCommand = errors.New("unknown command")
var ErrExpectedCommandLiteral = errors.New("expected command literal")

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

		if t.Peek().Type == argumentSeparator {
			t.Consume()

			continue
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
		return nil, fmt.Errorf("%w, but got '%s'", ErrExpectedCommandLiteral, token.Type)
	}

	var commandElem, commandFound = commandMappings[token.Value]
	if !commandFound {
		return nil, fmt.Errorf("%w '%s'", ErrUnknownCommand, token.Value)
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

	return elem, fmt.Errorf("%w %T", ErrUnknownCommand, elem)
}
