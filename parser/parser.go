package parser

import (
	"shellr/tokenizer"
)

type ParserType int

const (
	NONE ParserType = iota
	PIPE
	BACKGROUND
	QUOTED_ARGUMENT
	ARGUMENT
	COMMAND
)

type Entry struct {
	T   ParserType
	Val string
}

type Command struct {
	E []Entry
}

// Parse returns a list of commands
func Parse(tokens []tokenizer.Token) []Command {
	cmdTable := []Command{}
	row := []Entry{}
	for _, t := range tokens {
		e := Entry{}
		switch t.T {
		case tokenizer.STRING:
			if t.Data[0] == '"' {
				e.T = QUOTED_ARGUMENT
			} else {
				if len(row) > 0 && (row[len(row)-1].T != PIPE || row[len(row)-1].T == COMMAND) {
					e.T = ARGUMENT
				} else {
					e.T = COMMAND
				}
			}
			e.Val = t.Data
			row = append(row, e)
		case tokenizer.PIPE:
			if len(row) > 0 {
				cmdTable = append(cmdTable, Command{E: row})
			}
			row = []Entry{
				{T: PIPE, Val: "|"},
			}
		case tokenizer.SEMICOLON:
			if len(row) > 0 {
				cmdTable = append(cmdTable, Command{E: row})
			}
			row = []Entry{}
		case tokenizer.BACKGROUND:
			row = append([]Entry{{T: BACKGROUND, Val: "&"}}, row...)
			if len(row) > 0 {
				cmdTable = append(cmdTable, Command{E: row})
			}
			row = []Entry{}
		default:
			continue
		}
	}
	cmdTable = append(cmdTable, Command{E: row})
	return cmdTable
}
