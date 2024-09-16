package pokedex

import (
	"github.com/johndosdos/pokedexcli/internal/catch"
)

var dex = make(map[string]catch.PokemonData)

func Add(pokemonData catch.PokemonData) {
	pokemon := pokemonData.Name

	dex[pokemon] = pokemonData
}
