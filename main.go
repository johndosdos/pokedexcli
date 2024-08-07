package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		command := scanner.Text()
		fmt.Println()
		fmt.Println(command)
		fmt.Println()
	}
}
