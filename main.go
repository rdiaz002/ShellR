package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type TokenType int

const (
	NONE TokenType = iota
	COMMAND
	ARGUMENT
	QUOTED_ARGUMENT
	PIPE
	SEMICOLON
	BACKGROUND
)

type Token struct {
	T    TokenType
	Data string
}

var (
	PIPE_TOKEN       = Token{T: PIPE, Data: "|"}
	SEMICOLON_TOKEN  = Token{T: SEMICOLON, Data: ";"}
	BACKGROUND_TOKEN = Token{T: BACKGROUND, Data: "&"}
)

func ReadString(input *strings.Reader) (string, error) {
	buff := ""
	for {
		ch, _, err := input.ReadRune()
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("end of string not found")
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

func tokenize(input string) ([]Token, error) {
	tokens := []Token{}
	prevToken := Token{T: NONE}
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
			prevToken = SEMICOLON_TOKEN
		case '&':
			if currToken.T != NONE {
				tokens = append(tokens, currToken)
				currToken = Token{}
			}
			tokens = append(tokens, BACKGROUND_TOKEN)
			prevToken = BACKGROUND_TOKEN
		case '|':
			if currToken.T != NONE {
				tokens = append(tokens, currToken)
				currToken = Token{}
			}
			tokens = append(tokens, PIPE_TOKEN)
			prevToken = PIPE_TOKEN
		case ' ':
			if currToken.T != NONE {
				tokens = append(tokens, currToken)
				prevToken = currToken
				currToken = Token{}
			}
		case 10:
			continue
		case '"':
			if currToken.T != NONE {
				continue
			}
			if data, err := ReadString(r); err != nil {
				return nil, err
			} else {
				currToken.T = QUOTED_ARGUMENT
				currToken.Data += string('"') + data + string('"')
			}
		default:
			if currToken.T == NONE {
				if prevToken.T == COMMAND || prevToken.T == ARGUMENT {
					currToken.T = ARGUMENT
				} else {
					currToken.T = COMMAND
				}
			}
			currToken.Data += string(ch)
		}
	}
}

func main() {
	fmt.Println("Welcome to Ron Term. Its my little terminal project.")
	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(">>> ", line)
		tokens, err := tokenize(line)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(tokens)
	}
}
