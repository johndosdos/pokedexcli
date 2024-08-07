package main

import (
	"bufio"
	"fmt"
	"os"
)

type Command struct {
	Name        string
	Description string
	Execute     func()
}

func main() {
	commands := make(map[string]Command)
	commands = map[string]Command{
		"help": {
			Name:        "help",
			Description: "Display help screen",
			Execute: func() {
				fmt.Println("\tWelcome to PokeDex CLI!")
				fmt.Println()
				fmt.Println("\tAvailable commands: ")

				for input, cmd := range commands {
					fmt.Printf("\t\t%v -- %v", input, cmd.Description)
					fmt.Println()
				}
			},
		},

		"exit": {
			Name:        "exit",
			Description: "Exit from the program",
			Execute: func() {
				os.Exit(0)
			},
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		fmt.Println()
		command := scanner.Text()

		if cmd, ok := commands[command]; ok {
			cmd.Execute()
		} else {
			fmt.Println("\tError: Command not found. Type 'help' to see available commands")
		}

		fmt.Println()
	}
}
