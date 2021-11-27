package script

import (
	"reflect"
	"strconv"
	"testing"
)

func TestLex(t *testing.T) {
	testCases := []struct {
		Script string
		Tokens TokenStream
	}{
		{
			Script: "W",
			Tokens: TokenStream{Tokens: []Token{
				{Type: literal, Value: "W"},
			}},
		},
		{
			Script: "KD E;",
			Tokens: TokenStream{Tokens: []Token{
				{Type: literal, Value: "KD"},
				{Type: argumentSeparator},
				{Type: literal, Value: "E"},
				{Type: commandSeparator},
			}},
		},
		{
			Script: "R 25 [LD;W 60ms;LU;W 800ms]",
			Tokens: TokenStream{Tokens: []Token{
				{Type: literal, Value: "R"},
				{Type: argumentSeparator},
				{Type: literal, Value: "25"},
				{Type: argumentSeparator},
				{Type: blockOpen},
				{Type: literal, Value: "LD"},
				{Type: commandSeparator},
				{Type: literal, Value: "W"},
				{Type: argumentSeparator},
				{Type: literal, Value: "60ms"},
				{Type: commandSeparator},
				{Type: literal, Value: "LU"},
				{Type: commandSeparator},
				{Type: literal, Value: "W"},
				{Type: argumentSeparator},
				{Type: literal, Value: "800ms"},
				{Type: blockClose},
			}},
		},
	}

	var i = 0

	for _, testData := range testCases {
		var td = testData
		i++

		t.Run("case-"+strconv.Itoa(i), func(t *testing.T) {
			got, err := lex(td.Script)
			if err != nil {
				t.Fatalf("did not expect an error, but got: %v", err)
			}

			if !reflect.DeepEqual(&td.Tokens, got) {
				t.Fatalf("objects did not match:\ngot : %#v\nwant: %#v\n", got, td.Tokens)
			}
		})
	}
}
