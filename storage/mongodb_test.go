package storage

import (
	"fmt"
	"testing"

	guid "github.com/google/uuid"
	"github.com/rlongo/ambrosia/api"
)

const (
	collection = "test_collection"
	dbURL      = "mongodb://127.0.0.1:27017"
)

func getConnectionURL() string {
	return dbURL
}

func TestSetGetRecipes(t *testing.T) {
	m, err := OpenMongo(getConnectionURL(), collection)
	if err != nil {
		t.Errorf("Failed to connect to mongo: %s", err.Error())
	}

	r := api.Recipe{Name: "cake1", Author: "a1", Rating: 1, Tags: []string{"cake", "easter"}}

	err = m.AddRecipe(&r)
	if err != nil {
		t.Errorf("Failed to insert recipe: %s", err.Error())
	}

	fmt.Printf("Inserted recipe with UUID %s\n", guid.UUID(r.ID).String())
	r2, err := m.GetRecipe(r.ID)
	if err != nil {
		t.Errorf("Failed to get the recipe: %s", err.Error())
	}

	if r.ID != r2.ID || r.Name != r2.Name {
		t.Errorf("Recipe in DB didn't match")
	}

}
