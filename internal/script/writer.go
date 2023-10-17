package script

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
)

func Write(seq []sequencer.Elem) string {
	var (
		sb = strings.Builder{}
		t  TokenStream
	)

	write(&t, seq)

	for _, tok := range t.Tokens {
		switch tok.Type {
		case literal:
			sb.WriteString(tok.Value)
		case blockOpen:
			sb.WriteByte(valBlockOpen)
		case blockClose:
			sb.WriteByte(valBlockClose)
		case commandSeparator:
			sb.WriteByte(valCommandSep)
		case argumentSeparator:
			sb.WriteByte(valArgSep)
		}
	}

	return sb.String()
}

func write(t *TokenStream, seq []sequencer.Elem) {
	for _, elem := range seq {
		writeElement(elem, t)
	}
}

func writeElement(c sequencer.Elem, t *TokenStream) {
	if len(t.Tokens) > 0 {
		writeCommandSeparator(t)
	}

	switch v := c.(type) {
	case sequencer.Wait:
		writeWait(t, v)
	case general.KeyDown:
		writeKeyDown(t, v)
	case general.KeyUp:
		writeKeyUp(t, v)
	case sequencer.Repeat:
	case general.MouseMove:
		writeMouseMove(t, v)
	case general.SetMousePos:
		writeSetMousePos(t, v)
	case general.LookupMousePos:
		writeLookupMousePos(t, v)
	case general.LeftMouseButtonDown:
		writeLeftMouseButtonDown(t, v)
	case general.LeftMouseButtonUp:
		writeLeftMouseButtonUp(t, v)
	case general.RightMouseButtonDown:
		writeRightMouseButtonDown(t, v)
	case general.RightMouseButtonUp:
		writeRightMouseButtonUp(t, v)
	case sequencer.Loop:
		writeLoop(t, v)
	case sequencer.NoOperation:
		writeNoOperation(t, v)
	default:
		panic(fmt.Sprintf("unknown element '%T'", c))
	}
}

func commandByType(elem sequencer.Elem) string {
	for name, t := range commandMappings {
		t1 := reflect.TypeOf(t()).String()
		t2 := reflect.TypeOf(elem).String()

		if t1 == t2 {
			return name
		}
	}

	return ""
}

func writeWait(t *TokenStream, v sequencer.Wait) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
	t.Tokens = append(t.Tokens, Token{Type: argumentSeparator})
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: v.Duration.String()})
}

func writeKeyDown(t *TokenStream, v general.KeyDown) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
	t.Tokens = append(t.Tokens, Token{Type: argumentSeparator})
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: stringFromKeyCode(v.Key)})
}

func writeKeyUp(t *TokenStream, v general.KeyUp) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
	t.Tokens = append(t.Tokens, Token{Type: argumentSeparator})
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: stringFromKeyCode(v.Key)})
}

func writeMouseMove(t *TokenStream, v general.MouseMove) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
	t.Tokens = append(t.Tokens, Token{Type: argumentSeparator})
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: strconv.Itoa(int(v.X))})
	t.Tokens = append(t.Tokens, Token{Type: argumentSeparator})
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: strconv.Itoa(int(v.Y))})
}

func writeSetMousePos(t *TokenStream, v general.SetMousePos) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
	t.Tokens = append(t.Tokens, Token{Type: argumentSeparator})
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: strconv.Itoa(int(v.Pos.X))})
	t.Tokens = append(t.Tokens, Token{Type: argumentSeparator})
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: strconv.Itoa(int(v.Pos.Y))})
}

func writeLookupMousePos(t *TokenStream, v general.LookupMousePos) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
}

func writeCommandSeparator(t *TokenStream) {
	t.Tokens = append(t.Tokens, Token{Type: commandSeparator})
}

func writeLeftMouseButtonDown(t *TokenStream, v general.LeftMouseButtonDown) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
}

func writeLeftMouseButtonUp(t *TokenStream, v general.LeftMouseButtonUp) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
}

func writeRightMouseButtonDown(t *TokenStream, v general.RightMouseButtonDown) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
}

func writeRightMouseButtonUp(t *TokenStream, v general.RightMouseButtonUp) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
}

func writeLoop(t *TokenStream, v sequencer.Loop) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
}

func writeNoOperation(t *TokenStream, v sequencer.NoOperation) {
	t.Tokens = append(t.Tokens, Token{Type: literal, Value: commandByType(v)})
}
