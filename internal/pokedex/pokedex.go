package pokedex

import (
	"fmt"

	"github.com/johndosdos/pokedexcli/internal/catch"
)

var dex = make(map[string]catch.PokemonData)

func Add(pokemonData catch.PokemonData) {
	pokemon := pokemonData.Name

	dex[pokemon] = pokemonData
}

func Read(pokemon string) (catch.PokemonData, error) {
	data, found := dex[pokemon]
	if !found {
		return catch.PokemonData{}, fmt.Errorf("Pokemon not found.")
	}

	return data, nil
}
