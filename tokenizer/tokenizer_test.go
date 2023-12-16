package tokenizer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		Name        string
		CommandLine string
		Want        []Token
	}{
		{
			Name:        "no special token types",
			CommandLine: "echo Hello World",
			Want: []Token{
				{T: STRING, Data: "echo"},
				{T: STRING, Data: "Hello"},
				{T: STRING, Data: "World"},
			},
		},
		{
			Name:        "command line with quotes",
			CommandLine: "echo \"Hello World\"",
			Want: []Token{
				{T: STRING, Data: "echo"},
				{T: STRING, Data: "\"Hello World\""},
			},
		},
		{
			Name:        "command line with pipes",
			CommandLine: "echo Hello World | grep Hello",
			Want: []Token{
				{T: STRING, Data: "echo"},
				{T: STRING, Data: "Hello"},
				{T: STRING, Data: "World"},
				{T: PIPE, Data: "|"},
				{T: STRING, Data: "grep"},
				{T: STRING, Data: "Hello"},
			},
		},
		{
			Name:        "command line with ampersand",
			CommandLine: "echo Hello World & grep Hello",
			Want: []Token{
				{T: STRING, Data: "echo"},
				{T: STRING, Data: "Hello"},
				{T: STRING, Data: "World"},
				{T: BACKGROUND, Data: "&"},
				{T: STRING, Data: "grep"},
				{T: STRING, Data: "Hello"},
			},
		},
		{
			Name:        "command line with semicolon",
			CommandLine: "echo Hello World; grep Hello",
			Want: []Token{
				{T: STRING, Data: "echo"},
				{T: STRING, Data: "Hello"},
				{T: STRING, Data: "World"},
				{T: SEMICOLON, Data: ";"},
				{T: STRING, Data: "grep"},
				{T: STRING, Data: "Hello"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := Tokenize(test.CommandLine)
			if err != nil {
				t.Errorf("Tokenize(%s) returned an error: %v", test.CommandLine, err)
			}
			if !cmp.Equal(got, test.Want) {
				t.Errorf("Tokenize(%s) unexpected results\n got: %v want: %v", test.CommandLine, got, test.Want)
			}
		})
	}
}

func TestTokenizerError(t *testing.T) {
	tests := []struct {
		Name        string
		CommandLine string
		Want        error
	}{
		{
			Name:        "no end of string",
			CommandLine: "echo \"Hello World",
			Want:        ErrNoEndOfString,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := Tokenize(test.CommandLine)
			if err == nil {
				t.Errorf("Tokenize(%s) did not return an error", test.CommandLine)
			}
			if err != test.Want {
				t.Errorf("Tokenize(%s) did not return the correct error.\n got: %v want: %v", test.CommandLine, err, test.Want)
			}
		})
	}
}
