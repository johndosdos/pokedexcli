package pokedex

import (
	"fmt"

	"github.com/johndosdos/pokedexcli/internal/pokemon"
)

type Pokedex struct {
	dex map[string]pokemon.PokemonData
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		dex: make(map[string]pokemon.PokemonData),
	}
}

func (p *Pokedex) Add(pokemonData pokemon.PokemonData) {
	pokemonName := pokemonData.Name
	p.dex[pokemonName] = pokemonData
}

func (p *Pokedex) Read(pokemonName string) (pokemon.PokemonData, error) {
	data, found := p.dex[pokemonName]
	if !found {
		return pokemon.PokemonData{}, fmt.Errorf("Pokemon not found.")
	}

	return data, nil
}

func (p *Pokedex) ListCaught() error {
	if len(p.dex) == 0 {
		return fmt.Errorf("Pokedex is empty, cannot print.")
	}

	fmt.Println("\tCaught Pokemons:")

	counter := 1
	for pokemonName := range p.dex {
		fmt.Printf("\t%v: %v\n", counter, pokemonName)
		counter++
	}

	return nil
}
