package memory

import (
	"fmt"

	"github.com/rlongo/ambrosia/api"
)

type AmbrosiaStorageMemory struct {
	RecipesDB api.Recipes
}

// Open fakes a new DB connection
func Open(storageConnectionString string) (*AmbrosiaStorageMemory, error) {
	return &AmbrosiaStorageMemory{}, nil
}

func tagsFilter(tags []string, filter []string) bool {

	if filter == nil || len(filter) == 0 {
		return true
	}

	for _, f := range filter {
		found := false
		for _, t := range tags {
			if t == f {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (db *AmbrosiaStorageMemory) GetRecipes(filterTags []string, filterAuthor string) (api.Recipes, error) {
	var results api.Recipes

	for _, r := range db.RecipesDB {
		if (filterTags == nil || tagsFilter(r.Tags, filterTags)) && (len(filterAuthor) == 0 || filterAuthor == r.Author) {
			results = append(results, r)
		}
	}

	return results, nil
}

func (db *AmbrosiaStorageMemory) GetRecipe(id api.RecipeID) (api.Recipe, error) {
	return api.Recipe{}, fmt.Errorf("Not Implemented")
}

func (db *AmbrosiaStorageMemory) AddRecipe(recipe *api.Recipe) error {
	return fmt.Errorf("Not Implemented")
}

func (db *AmbrosiaStorageMemory) UpdateRecipe(recipe *api.Recipe) error {
	return fmt.Errorf("Not Implemented")
}
