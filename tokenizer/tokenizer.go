package tokenizer

import (
	"fmt"
	"io"
	"strings"
)

type TokenType int

const (
	NONE TokenType = iota
	PIPE
	SEMICOLON
	BACKGROUND
	STRING
)

type Token struct {
	T    TokenType
	Data string
}

var (
	PIPE_TOKEN       = Token{T: PIPE, Data: "|"}
	SEMICOLON_TOKEN  = Token{T: SEMICOLON, Data: ";"}
	BACKGROUND_TOKEN = Token{T: BACKGROUND, Data: "&"}

	ErrNoEndOfString = fmt.Errorf("end of string not found")
)

// readString will read the input until a '"' or EOF character is encountered.
func readString(input *strings.Reader) (string, error) {
	buff := ""
	for {
		ch, _, err := input.ReadRune()
		if err != nil {
			if err == io.EOF {
				return "", ErrNoEndOfString
			}
			return "", err
		}
		switch ch {
		case '"':
			return buff, nil
		default:
			buff += string(ch)
		}
	}
}

// Tokenize handles formatting a string into an array of tokens.
func Tokenize(input string) ([]Token, error) {
	tokens := []Token{}
	r := strings.NewReader(input)
	currToken := Token{T: NONE}
	for {
		var ch rune
		var err error
		if ch, _, err = r.ReadRune(); err != nil {
			if err == io.EOF {
				tokens = append(tokens, currToken)
				return tokens, nil
			}
			return nil, err
		}
		switch ch {
		case ';':
			if currToken.T != NONE {
				tokens = append(tokens, currToken)
				currToken = Token{}
			}
			tokens = append(tokens, SEMICOLON_TOKEN)
		case '&':
			if currToken.T != NONE {
				tokens = append(tokens, currToken)
				currToken = Token{}
			}
			tokens = append(tokens, BACKGROUND_TOKEN)
		case '|':
			if currToken.T != NONE {
				tokens = append(tokens, currToken)
				currToken = Token{}
			}
			tokens = append(tokens, PIPE_TOKEN)
		case ' ':
			if currToken.T != NONE {
				tokens = append(tokens, currToken)
				currToken = Token{}
			}
		case 10:
			continue
		case '"':
			if currToken.T != NONE {
				continue
			}
			if data, err := readString(r); err != nil {
				return nil, err
			} else {
				currToken.T = STRING
				currToken.Data += string('"') + data + string('"')
			}
		default:
			if currToken.T == NONE {
				currToken.T = STRING
			}
			currToken.Data += string(ch)
		}
	}
}
