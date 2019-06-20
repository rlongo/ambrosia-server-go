package storage

import (
	"fmt"
	"os"
	"testing"

	guid "github.com/google/uuid"
	"github.com/rlongo/ambrosia-server-go/api"
)

const (
	collection = "test_collection"
	dbURL      = "mongodb://127.0.0.1:27017"
)

func getConnectionURL() string {
	if v := os.Getenv("DB_URL"); len(v) > 0 {
		return v
	}
	return dbURL
}

func TestGracefullyHandelsInvalidConections(t *testing.T) {
	_, err := OpenMongo("fake", collection)
	if err == nil {
		t.Errorf("Connection should have failed")
	}
}

func TestGetMissingRecipeFails(t *testing.T) {
	m, err := OpenMongo(getConnectionURL(), collection)
	if err != nil {
		t.Errorf("Failed to connect to mongo: %s", err.Error())
	}
	defer m.Close()

	uuid, _ := guid.NewUUID()
	r, err := m.GetRecipe(api.RecipeID(uuid))

	if err == nil {
		t.Errorf("Shouldn't have found any recipes, got %s", r.Name)
	}
}

func TestSetGetRecipe(t *testing.T) {
	m, err := OpenMongo(getConnectionURL(), collection)
	if err != nil {
		t.Errorf("Failed to connect to mongo: %s", err.Error())
	}
	defer m.Close()

	r := api.Recipe{Name: "cake1", Author: "testAuthor", Rating: 1, Tags: []string{"1", "2"}}

	err = m.AddRecipe(&r)
	if err != nil {
		t.Errorf("Failed to insert recipe: %s", err.Error())
	}

	r2, err := m.GetRecipe(r.ID)
	if err != nil {
		t.Errorf("Failed to get the recipe: %s", err.Error())
	}

	if r.ID != r2.ID || r.Name != r2.Name {
		t.Errorf("Recipe in DB didn't match")
	}
}

func TestGetRecipes(t *testing.T) {
	recipesDB := api.Recipes{
		api.Recipe{Name: "cake1", Author: "a1", Rating: 1, Tags: []string{"cake", "easter"}},
		api.Recipe{Name: "cake2", Author: "a1", Rating: 1, Tags: []string{"cake", "xmas"}},
		api.Recipe{Name: "cake3", Author: "a2", Rating: 1, Tags: []string{"cake", "xmas", "nye"}},
		api.Recipe{Name: "pie1", Author: "a1", Rating: 1, Tags: []string{"pie", "chocolate"}},
		api.Recipe{Name: "pie2", Author: "a2", Rating: 1, Tags: []string{"pie", "easter"}},
		api.Recipe{Name: "pie3", Author: "a3", Rating: 1, Tags: []string{"pie", "bday"}},
		api.Recipe{Name: "cookie1", Author: "a3", Rating: 1, Tags: []string{"cookie", "bday"}},
		api.Recipe{Name: "cookie2", Author: "a2", Rating: 1, Tags: []string{"cookie", "xmas", "nye"}},
		api.Recipe{Name: "cookie3", Author: "a1", Rating: 1, Tags: []string{"cookie", "easter"}},
	}

	m, err := OpenMongo(getConnectionURL(), collection)
	if err != nil {
		t.Errorf("Failed to connect to mongo: %s", err.Error())
	}
	defer m.Close()

	for _, recipe := range recipesDB {
		err = m.AddRecipe(&recipe)
		if err != nil {
			fmt.Errorf("Failed to insert recipe %s", recipe.Name)
		}
	}

	// Single Tag
	recipes, err := m.GetRecipes([]string{"xmas"}, "")
	if len(recipes) != 3 {
		t.Errorf("Failed to get required recipes for single tag. Pulled %d results", len(recipes))
	}

	// Multi Tag
	recipes, err = m.GetRecipes([]string{"xmas", "cookie"}, "")
	if len(recipes) != 1 {
		t.Errorf("Failed to get required recipes for multi-tag. Pulled %d results", len(recipes))
	}

	// Multi Tag + Author
	recipes, err = m.GetRecipes([]string{"easter"}, "a1")
	if len(recipes) != 2 {
		t.Errorf("Failed to get required recipes for single tag + author. Pulled %d results", len(recipes))
	}

	// Author
	recipes, err = m.GetRecipes(nil, "a1")
	if len(recipes) != 4 {
		t.Errorf("Failed to get required recipes for author. Pulled %d results", len(recipes))
	}
}

func TestUpdateMustExist(t *testing.T) {
	m, err := OpenMongo(getConnectionURL(), collection)
	if err != nil {
		t.Errorf("Failed to connect to mongo: %s", err.Error())
	}
	defer m.Close()

	r := api.Recipe{Name: "cake1", Author: "testAuthor", Rating: 1, Tags: []string{"1", "2"}}

	err = m.UpdateRecipe(&r)
	if err == nil {
		t.Errorf("Update should have failed")
	}
}

func TestUpdate(t *testing.T) {
	m, err := OpenMongo(getConnectionURL(), collection)
	if err != nil {
		t.Errorf("Failed to connect to mongo: %s", err.Error())
	}
	defer m.Close()

	r := api.Recipe{Name: "cake1", Author: "testAuthor", Rating: 1, Tags: []string{"1", "2"}}

	err = m.AddRecipe(&r)
	if err != nil {
		t.Errorf("Failed to insert recipe: %s", err.Error())
	}

	r.Name = "cake2"
	err = m.UpdateRecipe(&r)
	if err != nil {
		t.Errorf("Failed to update recipe: %s", err.Error())
	}

	r2, err := m.GetRecipe(r.ID)
	if err != nil {
		t.Errorf("Failed to get the recipe: %s", err.Error())
	}

	if r.ID != r2.ID || r.Name != r2.Name {
		t.Errorf("Recipe in DB didn't match")
	}
}
