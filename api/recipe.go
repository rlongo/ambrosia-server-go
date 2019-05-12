package api

type RecipeID uint64

// The Recipe struct is a wrapper for the entire recipe
type Recipe struct {
	ID     RecipeID `json:"_id" bson:"_id"`
	Name   string   `json:"name" bson:"name"`
	Author string   `json:"author" bson:"author"`
	Rating uint8    `json:"rating" bson:"rating"`
	Notes  string   `json:"notes" bson:"notes"`
	Tags   []string `json:"tags" bson:"tags"`
	Staves Stages   `json:"stages" bson:"stages"`
}

// The Recipes is a convenience wrapper for all recipe types
type Recipes []Recipe

// The StorageServiceRecipes is an interface for the backend storage service
type StorageServiceRecipes interface {
	GetRecipes(filterTags []string, filterAuthor string) (Recipes, error)
	GetRecipe(id RecipeID) (Recipe, error)
	AddRecipe(recipe *Recipe) error
	UpdateRecipe(recipe *Recipe) error
}
