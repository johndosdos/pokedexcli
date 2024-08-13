package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Command struct {
	Name        string
	Description string
	Execute     func()
}

type LocationAreaResponse struct {
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
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

		"map": {
			Name:        "map",
			Description: "Display 20 locations at a time",
			Execute: func() {
				res, errGet := http.Get("https://pokeapi.co/api/v2/location-area/")
				if errGet != nil {
					fmt.Println(errGet)
				}
				defer res.Body.Close()

				locations := LocationAreaResponse{}
				if err := json.NewDecoder(res.Body).Decode(&locations); err != nil {
					fmt.Println(err)
				}

				for _, location := range locations.Results {
					fmt.Println(location.Name)
				}
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
