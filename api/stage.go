package api

// The Stage struct represents a subset of a recipe
type Stage struct {
	Name        string      `json:"name" bson:"name"`
	Notes       string      `json:"notes" bson:"notes"`
	Ingredients Ingredients `json:"ingredients" bson:"ingredients"`
	Steps       []string    `json:"steps" bson:"steps"`
}

type Stages []Stage
