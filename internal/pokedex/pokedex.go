package pokedex

import (
	"fmt"

	"github.com/johndosdos/pokedexcli/internal/catch"
)

type Pokedex struct {
	dex map[string]catch.PokemonData
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		dex: make(map[string]catch.PokemonData),
	}
}

func (p *Pokedex) Add(pokemonData catch.PokemonData) {
	pokemon := pokemonData.Name
	p.dex[pokemon] = pokemonData
}

func (p *Pokedex) Read(pokemon string) (catch.PokemonData, error) {
	data, found := p.dex[pokemon]
	if !found {
		return catch.PokemonData{}, fmt.Errorf("Pokemon not found.")
	}

	return data, nil
}
