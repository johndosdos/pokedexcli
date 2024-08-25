package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/johndosdos/pokedexcli/internal/pokecache"
)

type Command struct {
	Name        string
	Description string
	Execute     func()
	MapExecute  func(locations *LocationAreaResponse)
}

type LocationAreaResponse struct {
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func processURL(url string) LocationAreaResponse {
	res, errGet := http.Get(url)
	if errGet != nil {
		log.Printf("Error fetching URL %s: %v", url, errGet)
		return LocationAreaResponse{}
	}
	defer res.Body.Close()

	locations := LocationAreaResponse{}
	if err := json.NewDecoder(res.Body).Decode(&locations); err != nil {
		fmt.Println(err)
	}

	return locations
}

func main() {
	locations := LocationAreaResponse{}
	cache := pokecache.NewCache(5)

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
			MapExecute: func(locations *LocationAreaResponse) {
				if locations.Results == nil && locations.Next == "" {
					baseURL := "https://pokeapi.co/api/v2/location-area/"
					*locations = processURL(baseURL)
					cache.Add(baseURL, locations.Results)
				} else {
					if result, ok := cache.Get(locations.Next); ok {
						locations.Results = result
					} else {
						*locations = processURL(locations.Next)
					}
				}

				for _, location := range locations.Results {
					fmt.Println(location.Name)
				}
			},
		},

		"mapb": {
			Name:        "mapb",
			Description: "Display the previous 20 locations",
			MapExecute: func(locations *LocationAreaResponse) {
				if prevURL, ok := locations.Previous.(string); ok {
					if result, ok := cache.Get(prevURL); ok {
						locations.Results = result
					} else {
						*locations = processURL(prevURL)
					}

					for _, location := range locations.Results {
						fmt.Println(location.Name)
					}
				} else {
					fmt.Println("This is the first page.")
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
			if command == "map" || command == "mapb" {
				cmd.MapExecute(&locations)
			} else {
				cmd.Execute()
			}
		} else {
			fmt.Println("\tError: Command not found. Type 'help' to see available commands")
		}

		fmt.Println()
	}
}
