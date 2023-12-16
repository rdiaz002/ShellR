package parser

import (
	"shellr/tokenizer"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	tests := []struct {
		Name   string
		Tokens []tokenizer.Token
		Want   []Command
	}{
		{
			Name: "Single command",
			Tokens: []tokenizer.Token{
				{
					T:    tokenizer.STRING,
					Data: "echo",
				},
				{
					T:    tokenizer.STRING,
					Data: "Hello",
				},
			},
			Want: []Command{
				{
					[]Entry{
						{T: COMMAND, Val: "echo"},
						{T: ARGUMENT, Val: "Hello"},
					},
				},
			},
		},
		{
			Name: "Single command with quotes",
			Tokens: []tokenizer.Token{
				{
					T:    tokenizer.STRING,
					Data: "echo",
				},
				{
					T:    tokenizer.STRING,
					Data: "\"Hello\"",
				},
			},
			Want: []Command{
				{
					[]Entry{
						{T: COMMAND, Val: "echo"},
						{T: QUOTED_ARGUMENT, Val: "\"Hello\""},
					},
				},
			},
		},
		{
			Name: "Multiple commands with pipe",
			Tokens: []tokenizer.Token{
				{
					T:    tokenizer.STRING,
					Data: "echo",
				},
				{
					T:    tokenizer.STRING,
					Data: "\"Hello\"",
				},
				{
					T:    tokenizer.PIPE,
					Data: "|",
				},
				{
					T:    tokenizer.STRING,
					Data: "ping",
				},
			},
			Want: []Command{
				{
					[]Entry{
						{T: COMMAND, Val: "echo"},
						{T: QUOTED_ARGUMENT, Val: "\"Hello\""},
					},
				},
				{
					[]Entry{
						{T: PIPE, Val: "|"},
						{T: COMMAND, Val: "ping"},
					},
				},
			},
		},
		{
			Name: "Multiple commands with semicolon",
			Tokens: []tokenizer.Token{
				{
					T:    tokenizer.STRING,
					Data: "echo",
				},
				{
					T:    tokenizer.STRING,
					Data: "\"Hello\"",
				},
				{
					T:    tokenizer.SEMICOLON,
					Data: ";",
				},
				{
					T:    tokenizer.STRING,
					Data: "ping",
				},
			},
			Want: []Command{
				{
					[]Entry{
						{T: COMMAND, Val: "echo"},
						{T: QUOTED_ARGUMENT, Val: "\"Hello\""},
					},
				},
				{
					[]Entry{
						{T: COMMAND, Val: "ping"},
					},
				},
			},
		},
		{
			Name: "Multiple commands with ampersand",
			Tokens: []tokenizer.Token{
				{
					T:    tokenizer.STRING,
					Data: "echo",
				},
				{
					T:    tokenizer.STRING,
					Data: "\"Hello\"",
				},
				{
					T:    tokenizer.BACKGROUND,
					Data: "&",
				},
				{
					T:    tokenizer.SEMICOLON,
					Data: ";",
				},
				{
					T:    tokenizer.STRING,
					Data: "ping",
				},
			},
			Want: []Command{
				{
					[]Entry{
						{T: BACKGROUND, Val: "&"},
						{T: COMMAND, Val: "echo"},
						{T: QUOTED_ARGUMENT, Val: "\"Hello\""},
					},
				},
				{
					[]Entry{
						{T: COMMAND, Val: "ping"},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cmdTable := Parse(test.Tokens)

			if !cmp.Equal(cmdTable, test.Want) {
				t.Errorf("Parse(%v) return value is unexpected\n got: %v want: %v", test.Tokens, cmdTable, test.Want)
			}
		})
	}
}
