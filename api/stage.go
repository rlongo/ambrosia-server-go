package api

// The Stage struct represents a subset of a recipe
type Stage struct {
	Name        string       `json:"name"`
	Notes       string       `json:"notes"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
}
