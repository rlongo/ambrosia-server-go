package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urfave/negroni"

	"github.com/rlongo/ambrosia/api"
	"github.com/rlongo/ambrosia/storage/mock"
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

func TestGETRecipesTest(t *testing.T) {
	recipesDB := api.Recipes{
		&api.Recipe{ID: 0, Name: "cake1", Author: "a1", Rating: 1, Tags: []string{"cake", "easter"}},
	}

	ambrosiaDB := mock.MockAmbrosiaStorage{recipesDB}

	t.Run("no query params", func(t *testing.T) {
		router := NewRouter(&ambrosiaDB, negroni.New())
		request, _ := http.NewRequest(http.MethodGet, "/api/v1/dojang/tests", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBodyTests(t, response.Body.Bytes(), recipesDB)
	})

}
