package api

// The Ingredient struct represents a single ingredient in the any step of the recipe
type Ingredient struct {
	Name     string `json:"name" bson:"name"`
	Unit     string `json:"unit" bson:"unit"`
	Quantity float32 `json:"quantity" bson:"quantity"`
}

type Ingredients []Ingredient
