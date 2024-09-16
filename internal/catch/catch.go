package catch

import (
	"fmt"
	"math"
	"math/rand"
)

type PokemonData struct {
	BaseExperience         int     `json:"base_experience"`
	Height                 int     `json:"height"`
	LocationAreaEncounters string  `json:"location_area_encounters"`
	Name                   string  `json:"name"`
	Stats                  []Stats `json:"stats"`
	Types                  []Types `json:"types"`
	Weight                 int     `json:"weight"`
}
type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Stats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"stat"`
}
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

type PokemonSpecies struct {
	CaptureRate int `json:"capture_rate"`
}

func Catch(pokemonData PokemonData, catchRateData PokemonSpecies) error {
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
