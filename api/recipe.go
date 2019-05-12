package api

type RecipeID uint64

// The Recipe struct is a wrapper for the entire recipe
type Recipe struct {
	ID     RecipeID `json:"_id"`
	Name   string   `json:"name"`
	Author string   `json:"author"`
	Rating uint8    `json:"rating"`
	Notes  string   `json:"notes"`
	Tags   []string `json:"tags"`
	Staves []Stage  `json:"stages"`
}

// The Recipes is a convenience wrapper for all recipe types
type Recipes []*Recipe

// The StorageServiceRecipes is an interface for the backend storage service
type StorageServiceRecipes interface {
	GetRecipes(filterTags []string, filterAuthor string, filterMinRating uint8) (Recipes, error)
	GetRecipe(id RecipeID) (Recipe, error)
	AddRecipe(recipe *Recipe) error
	UpdateRecipe(recipe *Recipe) error
}
