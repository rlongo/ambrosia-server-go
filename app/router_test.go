package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urfave/negroni"

	"github.com/rlongo/ambrosia/api"
	"github.com/rlongo/ambrosia/storage/memory"
)

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBodyTests(t *testing.T, got []byte, want api.Recipes) {
	t.Helper()
	var results api.Recipes

	if err := json.Unmarshal(got, &results); err != nil {
		t.Errorf("Response was invalid JSON")
	}

	if len(results) != len(want) {
		t.Errorf("Response size is wrong. Expected: %d, Got: %d", len(want), len(results))
	}

	for i := range results {
		if results[i].ID != want[i].ID {
			t.Errorf("Response mismatch at index %d. Expected: '%d', Got: '%d'",
				i, want[i].ID, results[i].ID)
		}
	}
}

func assertResponseBodyTest(t *testing.T, got []byte, want api.Recipe) {
	t.Helper()
	var result api.Recipe

	if err := json.Unmarshal(got, &result); err != nil {
		t.Errorf("Response was invalid JSON")
	}

	if result.ID != want.ID {
		t.Errorf("Response mismatch. Expected: '%d', Got: '%d'",
			want.ID, result.ID)
	}
}

func TestSearchRecipes(t *testing.T) {
	recipesDB := api.Recipes{
		&api.Recipe{ID: 0, Name: "cake1", Author: "a1", Rating: 1, Tags: []string{"cake", "easter"}},
		&api.Recipe{ID: 1, Name: "cake2", Author: "a1", Rating: 1, Tags: []string{"cake", "xmas"}},
		&api.Recipe{ID: 2, Name: "cake3", Author: "a2", Rating: 1, Tags: []string{"cake", "xmas", "nye"}},
		&api.Recipe{ID: 3, Name: "pie1", Author: "a1", Rating: 1, Tags: []string{"pie", "chocolate"}},
		&api.Recipe{ID: 4, Name: "pie2", Author: "a2", Rating: 1, Tags: []string{"pie", "easter"}},
		&api.Recipe{ID: 5, Name: "pie3", Author: "a3", Rating: 1, Tags: []string{"pie", "bday"}},
		&api.Recipe{ID: 6, Name: "cookie1", Author: "a3", Rating: 1, Tags: []string{"cookie", "bday"}},
		&api.Recipe{ID: 7, Name: "cookie2", Author: "a2", Rating: 1, Tags: []string{"cookie", "xmas", "nye"}},
		&api.Recipe{ID: 8, Name: "cookie3", Author: "a1", Rating: 1, Tags: []string{"cookie", "easter"}},
	}

	ambrosiaDB := memory.AmbrosiaStorageMemory{recipesDB}

	testRunner := func(path string, expected api.Recipes) {
		router := NewRouter(&ambrosiaDB, negroni.New())
		request, _ := http.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBodyTests(t, response.Body.Bytes(), expected)
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

func TestGETRecipes(t *testing.T) {
	recipesDB := api.Recipes{
		&api.Recipe{ID: 0, Name: "cake1", Author: "a1", Rating: 1, Tags: []string{"cake", "easter"}},
		&api.Recipe{ID: 3, Name: "pie1", Author: "a1", Rating: 1, Tags: []string{"pie", "chocolate"}},
		&api.Recipe{ID: 7, Name: "cookie2", Author: "a2", Rating: 1, Tags: []string{"cookie", "xmas", "nye"}},
	}

	ambrosiaDB := memory.AmbrosiaStorageMemory{recipesDB}

	testRunner := func(path string, expected api.Recipe) {
		router := NewRouter(&ambrosiaDB, negroni.New())
		request, _ := http.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBodyTest(t, response.Body.Bytes(), expected)
	}

	t.Run("query succeeds", func(t *testing.T) {
		testRunner("/api/v1/recipe/3", *recipesDB[1])
	})
}
