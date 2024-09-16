package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/johndosdos/pokedexcli/internal/catch"
	"github.com/johndosdos/pokedexcli/internal/explore"
	"github.com/johndosdos/pokedexcli/internal/pokecache"
)

type Command struct {
	Name           string
	Description    string
	Execute        func()
	MapExecute     func(locations *LocationAreaResponse) error
	ExploreExecute func(location string) error
	PokemonExecute func(pokemon string) error
}

type LocationAreaResponse struct {
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func processURL(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching URL %s: %w", url, err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading from response body: %w", err)
	}

	return data, nil
}

func main() {
	locations := LocationAreaResponse{}
	cache := pokecache.NewCache(2 * time.Second)
	defer cache.Stop()

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
			MapExecute: func(locations *LocationAreaResponse) error {
				if locations.Results == nil && locations.Next == "" {
					baseURL := "https://pokeapi.co/api/v2/location-area/"

					rawData, err := processURL(baseURL)
					if err != nil {
						return fmt.Errorf("Error processing URL: %w", err)
					}

					if err := json.Unmarshal(rawData, &locations); err != nil {
						return fmt.Errorf("Parse error: %w", err)
					}

					cache.Add(baseURL, rawData)
				} else {
					nextURL := locations.Next

					if rawData, ok := cache.Get(nextURL); ok {
						if err := json.Unmarshal(rawData, &locations); err != nil {
							return fmt.Errorf("Parse error: %w", err)
						}
					} else {
						rawData, err := processURL(nextURL)

						if err != nil {
							return fmt.Errorf("Error processing URL: %w", err)
						}

						if err := json.Unmarshal(rawData, &locations); err != nil {
							return fmt.Errorf("Parse error: %w", err)
						}

						cache.Add(nextURL, rawData)
					}
				}

				for _, location := range locations.Results {
					fmt.Println(location.Name)
				}

				return nil
			},
		},

		"mapb": {
			Name:        "mapb",
			Description: "Display the previous 20 locations",
			MapExecute: func(locations *LocationAreaResponse) error {
				if locations.Previous != nil {
					prevURL := locations.Previous.(string)

					if rawData, ok := cache.Get(prevURL); ok {
						if err := json.Unmarshal(rawData, &locations); err != nil {
							return fmt.Errorf("Parse error: %w", err)
						}
					} else {
						rawData, err := processURL(prevURL)
						if err != nil {
							return fmt.Errorf("Error processing URL: %w", err)
						}

						if err := json.Unmarshal(rawData, &locations); err != nil {
							return fmt.Errorf("Parse error: %w", err)
						}

						cache.Add(prevURL, rawData)
					}

					for _, location := range locations.Results {
						fmt.Println(location.Name)
					}

				} else {
					fmt.Println("This is the first page.")
				}

				return nil
			},
		},

		"explore": {
			Name:        "explore",
			Description: "Explore locations from map(b)",
			ExploreExecute: func(location string) error {
				ExploreResponseData := explore.Response{}
				cache := pokecache.NewCache(5)

				url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", location)

				var (
					data  []byte
					err   error
					found bool
				)

				data, found = cache.Get(url)
				if !found {
					data, err = processURL(url)
					if err != nil {
						return fmt.Errorf("Error fetching URL: %v", err)
					}

					// Add data to cache
					cache.Add(url, data)
				}

				if err = json.Unmarshal(data, &ExploreResponseData); err != nil {
					return fmt.Errorf("Parse error: %v", err)
				}

				for _, item := range ExploreResponseData.PokemonEncounters {
					fmt.Printf("\t* %v", item.Pokemon.Name)
					fmt.Println()
				}

				return nil
			},
		},

		"catch": {
			Name:        "Catch",
			Description: "Catch pokemon present in that area",
			PokemonExecute: func(pokemon string) error {
				PokemonDataResponse := catch.PokemonData{}
				PokemonSpeciesResponse := catch.PokemonSpecies{}

				// process pokemon data
				url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v/", pokemon)
				data, err := processURL(url)
				if err != nil {
					return fmt.Errorf("Error fetching URL: %v", err)
				}
				if err := json.Unmarshal(data, &PokemonDataResponse); err != nil {
					return fmt.Errorf("Parse error: %v", err)
				}

				// process pokemon-species data, specifically the pokemon catch rate
				url = fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%v/", pokemon)
				data, err = processURL(url)
				if err != nil {
					return fmt.Errorf("Error fetching URL: %v", err)
				}
				if err := json.Unmarshal(data, &PokemonSpeciesResponse); err != nil {
					return fmt.Errorf("Parse error: %v", err)
				}

				// call to execute pokemon capture
				err = catch.Catch(PokemonDataResponse, PokemonSpeciesResponse)
				if err != nil {
					return fmt.Errorf("Error catching pokemon: %v", err)
				}

				return nil
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
		args := strings.Fields(command)
		mainArg := args[0]
		secondArg := ""
		if len(args) == 2 {
			secondArg = args[1]
		}

		if cmd, ok := commands[mainArg]; ok {
			if mainArg == "map" || mainArg == "mapb" {
				if err := cmd.MapExecute(&locations); err != nil {
					log.Printf("Error: %v", err)
				}
			} else if mainArg == "explore" {
				if secondArg == "" {
					fmt.Println("\tPlease enter the locaton you want to explore.")
					fmt.Println("\tExample: \"explore [location]\"")
				} else {
					if err := cmd.ExploreExecute(secondArg); err != nil {
						log.Printf("Error executing command: %v", err)
					}
				}
			} else {
				cmd.Execute()
			}
		} else {
			fmt.Println("\tError: Command not found. Type 'help' to see available commands")
		}

		fmt.Println()
	}
}
