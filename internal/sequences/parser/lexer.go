package parser

import (
	"errors"
	"strings"
)

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
		case isCommandSeperator(stream):
			stream.Consume()

			tokenStream.Tokens = append(tokenStream.Tokens, Token{Type: commandSeparator})

		case isArgumentSeparator(stream):
			stream.Consume()

			tokenStream.Tokens = append(tokenStream.Tokens, Token{Type: argumentSeparator})

		case isLiteral(stream):
			t, err := literalToken(stream)
			if err != nil {
				return nil, err
			}

			tokenStream.Tokens = append(tokenStream.Tokens, t)
		case isBlockOpen(stream):
			stream.Consume()

			tokenStream.Tokens = append(tokenStream.Tokens, Token{Type: blockOpen})

		case isBlockClose(stream):
			stream.Consume()

			tokenStream.Tokens = append(tokenStream.Tokens, Token{Type: blockClose})
		default:
			return nil, errors.New("syntax error")
		}
	}

	return tokenStream, nil
}

func isBlockClose(stream *InputReader) bool {
	return stream.Peek(0) == ']'
}

func isBlockOpen(stream *InputReader) bool {
	return stream.Peek(0) == '['
}

func isCommandSeperator(stream *InputReader) bool {
	return stream.Peek(0) == ';'
}

func isArgumentSeparator(stream *InputReader) bool {
	return stream.Peek(0) == ' '
}

func isLiteral(stream *InputReader) bool {
	var (
		char           = stream.Peek(0)
		isTimeUnitChar = char == 'n' || char == 's' || char == 'h' || char == 'm' || char == 'u'
		isUCaseChar    = char >= 'A' && char <= 'Z'
		isLCaseChar    = char >= '0' && char <= '9'
	)

	return isUCaseChar || isLCaseChar || isTimeUnitChar
}

func literalToken(stream *InputReader) (Token, error) {
	var sb = strings.Builder{}

	for !stream.isEOF() && isLiteral(stream) {
		sb.WriteByte(stream.Consume())
	}

	return Token{Type: literal, Value: sb.String()}, nil
}
