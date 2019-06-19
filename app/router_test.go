package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	guid "github.com/google/uuid"

	"github.com/rlongo/ambrosia/api"
	"github.com/rlongo/ambrosia/storage"
)

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertBodyRecipes(t *testing.T, got []byte, want api.Recipes) {
	t.Helper()
	var results api.Recipes

	if err := json.Unmarshal(got, &results); err != nil {
		t.Errorf("Response was invalid JSON")
	}

	if len(results) != len(want) {
		t.Errorf("Response size is wrong. Expected: %d, Got: %d", len(want), len(results))
	}

	set := make(map[string]bool)
	for _, v := range want {
		set[v.Name] = false
	}

	for _, v := range results {
		if _, ok := set[v.Name]; !ok {
			t.Errorf("Unexpected result '%s'", v.Name)
		}
		set[v.Name] = true
	}

	for k, v := range set {
		if !v {
			t.Errorf("Expected result not found '%s'", k)
		}
	}
}

func assertBodyRecipe(t *testing.T, got []byte, want api.Recipe) {
	t.Helper()
	var result api.Recipe

	if err := json.Unmarshal(got, &result); err != nil {
		t.Errorf("Response was invalid JSON")
	}

	if result.Name != want.Name {
		t.Errorf("Response mismatch. Expected: '%s', Got: '%s'",
			want.Name, result.Name)
	}
}

func TestSearchRecipes(t *testing.T) {
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

	ambrosiaDB := storage.AmbrosiaStorageMemory{recipesDB}

	testRunner := func(path string, expected api.Recipes) {
		router := NewRouter(&ambrosiaDB)
		request, _ := http.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBodyRecipes(t, response.Body.Bytes(), expected)
	}

	t.Run("no query params", func(t *testing.T) {
		testRunner("/api/v1/recipes", recipesDB)
	})

	t.Run("filter one tag", func(t *testing.T) {
		expected := api.Recipes{recipesDB[1], recipesDB[2], recipesDB[7]}
		testRunner("/api/v1/recipes?tag=xmas", expected)
	})

	t.Run("filter multi tag", func(t *testing.T) {
		expected := api.Recipes{recipesDB[7]}
		testRunner("/api/v1/recipes?tag=xmas&tag=cookie", expected)
	})

	t.Run("filter author", func(t *testing.T) {
		expected := api.Recipes{recipesDB[5], recipesDB[6]}
		testRunner("/api/v1/recipes?author=a3", expected)
	})

	t.Run("filter author + tag", func(t *testing.T) {
		expected := api.Recipes{recipesDB[0], recipesDB[8]}
		testRunner("/api/v1/recipes?author=a1&tag=easter", expected)
	})
}

func TestPOSTRecipes(t *testing.T) {
	recipe := api.Recipe{Name: "cake1", Author: "a1", Rating: 1, Tags: []string{"cake", "easter"}}
	ambrosiaDB := storage.AmbrosiaStorageMemory{}

	t.Run("post succeeds", func(t *testing.T) {
		router := NewRouter(&ambrosiaDB)

		expectedTestJSON, _ := json.Marshal(recipe)
		b := bytes.NewBuffer(expectedTestJSON)

		request, _ := http.NewRequest(http.MethodPost, "/api/v1/recipe", b)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusCreated)

		request, _ = http.NewRequest(http.MethodGet, "/api/v1/recipes", nil)
		response = httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBodyRecipes(t, response.Body.Bytes(), api.Recipes{recipe})
	})
}

func TestPUTRecipes(t *testing.T) {
	uuid, _ := guid.NewUUID()
	recipe := api.Recipe{ID: api.RecipeID(uuid), Name: "cake1", Author: "a1", Rating: 1, Tags: []string{"cake", "easter"}}
	ambrosiaDB := storage.AmbrosiaStorageMemory{RecipesDB: api.Recipes{recipe}}

	t.Run("put succeeds", func(t *testing.T) {
		router := NewRouter(&ambrosiaDB)

		recipe.Author = "a2"

		expectedTestJSON, _ := json.Marshal(recipe)
		b := bytes.NewBuffer(expectedTestJSON)

		idString := uuid.String()

		request, _ := http.NewRequest(http.MethodPut, "/api/v1/recipe/"+idString, b)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusCreated)

		request, _ = http.NewRequest(http.MethodGet, "/api/v1/recipe/"+idString, nil)
		response = httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		var result api.Recipe

		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Errorf("Response was invalid JSON")
		}

		if result.Author != "a2" {
			t.Errorf("Response mismatch. Expected: '%s', Got: '%s'",
				"a2", result.Author)
		}

	})
}
