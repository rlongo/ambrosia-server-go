package mock

import (
	"fmt"

	"github.com/rlongo/ambrosia/api"
)

type MockAmbrosiaStorage struct {
	RecipesDB api.Recipes
}

// Open fakes a new DB connection
func Open(storageConnectionString string) (*MockAmbrosiaStorage, error) {
	return &MockAmbrosiaStorage{}, nil
}

func (db *MockAmbrosiaStorage) GetRecipes(filterTags []string, filterAuthor string, filterMinRating uint8) (api.Recipes, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (db *MockAmbrosiaStorage) GetRecipe(id api.RecipeID) (api.Recipe, error) {
	return api.Recipe{}, fmt.Errorf("Not Implemented")
}

func (db *MockAmbrosiaStorage) AddRecipe(recipe *api.Recipe) error {
	return fmt.Errorf("Not Implemented")
}

func (db *MockAmbrosiaStorage) UpdateRecipe(recipe *api.Recipe) error {
	return fmt.Errorf("Not Implemented")
}
