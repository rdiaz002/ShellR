package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"shellr/parser"
	"shellr/tokenizer"
)

func main() {
	fmt.Println("Welcome to ShellR. Its my little Shell Project.")
	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(">>> ", line)
		tokens, err := tokenizer.Tokenize(line)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Tokens:")
		fmt.Println(tokens)
		cmdTable := parser.Parse(tokens)
		fmt.Println("cmdTable:")
		fmt.Println(cmdTable)
	}
}
