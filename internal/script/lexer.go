package script

import (
	"errors"
	"strings"
)

var ErrSyntaxError = errors.New("syntax error")

type InputReader struct {
	Data string
	Pos  int
}

func (s *InputReader) Peek(offset int) byte {
	return s.Data[s.Pos+offset]
}

func (s *InputReader) Consume() byte {
	var char = s.Data[s.Pos]
	s.Pos++

	return char
}

func (s *InputReader) isEOF() bool {
	return s.Pos >= len(s.Data)
}

func lex(script string) (*TokenStream, error) {
	var (
		stream      = &InputReader{Data: script}
		tokenStream = &TokenStream{}
	)

	for !stream.isEOF() {
		switch {
		case isCommandSeparator(stream):
			stream.Consume()

			tokenStream.Tokens = append(tokenStream.Tokens, Token{Type: commandSeparator})

		case isArgumentSeparator(stream):
			stream.Consume()

			tokenStream.Tokens = append(tokenStream.Tokens, Token{Type: argumentSeparator})

		case isLiteral(stream):
			t := literalToken(stream)

			tokenStream.Tokens = append(tokenStream.Tokens, t)
		case isBlockOpen(stream):
			stream.Consume()

			tokenStream.Tokens = append(tokenStream.Tokens, Token{Type: blockOpen})

		case isBlockClose(stream):
			stream.Consume()

			tokenStream.Tokens = append(tokenStream.Tokens, Token{Type: blockClose})
		default:
			return nil, ErrSyntaxError
		}
	}

	return tokenStream, nil
}

func isBlockClose(stream *InputReader) bool {
	return stream.Peek(0) == valBlockClose
}

func isBlockOpen(stream *InputReader) bool {
	return stream.Peek(0) == valBlockOpen
}

func isCommandSeparator(stream *InputReader) bool {
	return stream.Peek(0) == valCommandSep
}

func isArgumentSeparator(stream *InputReader) bool {
	return stream.Peek(0) == valArgSep
}

func isLiteral(stream *InputReader) bool {
	var (
		char            = stream.Peek(0)
		isTimeUnitChar  = char == 'n' || char == 's' || char == 'h' || char == 'm' || char == 'u'
		isUCaseChar     = char >= 'A' && char <= 'Z'
		isLCaseChar     = char >= '0' && char <= '9'
		isDecimalDelim  = char == '.'
		isNumericPrefix = char == '-'
	)

	return isUCaseChar || isLCaseChar || isTimeUnitChar || isDecimalDelim || isNumericPrefix
}

func literalToken(stream *InputReader) Token {
	var sb = strings.Builder{}

	for !stream.isEOF() && isLiteral(stream) {
		sb.WriteByte(stream.Consume())
	}

	return Token{Type: literal, Value: sb.String()}
}
