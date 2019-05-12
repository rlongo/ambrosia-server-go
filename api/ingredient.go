package api

// The Ingredient struct represents a single ingredient in the any step of the recipe
type Ingredient struct {
	Name     string `json:"name"`
	Unit     string `json:"unit"`
	Quantity uint32 `json:"quantity"`
}
