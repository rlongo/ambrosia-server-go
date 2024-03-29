package storage

import (
	"fmt"
	"log"

	guid "github.com/google/uuid"
	"github.com/rlongo/ambrosia-server-go/api"
)

type AmbrosiaStorageMemory struct {
	RecipesDB api.Recipes
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
	for _, r := range db.RecipesDB {
		if r.ID == id {
			return r, nil
		}
	}

	return api.Recipe{}, fmt.Errorf("Recipe not found")
}

func (db *AmbrosiaStorageMemory) AddRecipe(recipe *api.Recipe) error {
	log.Println("!!!Warning: AmbrosiaStorageMemory.AddRecipe not safe for production!")
	uuid, err := guid.NewUUID()
	if err != nil {
		return err
	}
	recipe.ID = api.RecipeID(uuid)
	db.RecipesDB = append(db.RecipesDB, *recipe)
	return nil
}

func (db *AmbrosiaStorageMemory) UpdateRecipe(recipe *api.Recipe) error {
	log.Println("!!!Warning: AmbrosiaStorageMemory.UpdateRecipe not safe for production!")
	for i, r := range db.RecipesDB {
		if r.ID == recipe.ID {
			db.RecipesDB[i] = *recipe
			return nil
		}
	}

	return fmt.Errorf("Recipe not found")
}
