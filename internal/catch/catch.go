package catch

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
