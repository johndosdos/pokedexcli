package actions

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/johndosdos/pokedexcli/internal/pokedex"
	"github.com/johndosdos/pokedexcli/internal/pokemon"
)

func Catch(pokemonData pokemon.PokemonData, catchRateData pokemon.PokemonSpecies, pokedex *pokedex.Pokedex) error {
	if pokemonData.Name == "" {
		return fmt.Errorf("No pokemon data found")
	}

	catchRate := catchRateData.CaptureRate

	// The formulas used are from the GEN 3 pokemon games.
	// Source: https://bulbapedia.bulbagarden.net/wiki/Catch_rate#Capture_method_(Generation_III-IV)

	// Probability = (((3 * Max HP - 2 * Current HP) * Modified Catch Rate * Ball Bonus) / (3 * Max HP)) * Status Bonus
	// Simplified catch rate formula.
	modifiedCatchRate := (float64(1.0) / 3) * float64(catchRate)

	// Shake probability. A shake is where the pokeball shakes before a successful or failed capture.
	shakeProb := int(1048560 / (math.Sqrt(math.Sqrt((16711680 / modifiedCatchRate)))))

	fmt.Printf("\tTrying to catch %v\n", pokemonData.Name)

	isSuccess := shakeCheck(shakeProb)
	if isSuccess {
		pokedex.Add(pokemonData)
		fmt.Printf("\t%v was caught! Nice work.", pokemonData.Name)
	} else {
		fmt.Printf("\tFailed to catch %v. Better luck next time!", pokemonData.Name)
	}

	return nil
}

func shakeCheck(shakeProb int) bool {
	shakeSuccessCount := 0

	// Four checks needed for successful pokemon capture.
	for i := 1; i < 5; i++ {

		// Emulate pokeball shake.
		fmt.Printf("\t%v...\n", i)
		time.Sleep(1 * time.Second)

		// Generate random number between 0 and 65535 inclusive.
		randNum := rand.Intn(65535 + 1)

		// Compare rand and shakeProb. If rand >= shakeProb, the check fails.
		if randNum >= shakeProb {
			break
		}

		shakeSuccessCount++
		if shakeSuccessCount == 4 {
			return true
		}
	}

	return false
}
