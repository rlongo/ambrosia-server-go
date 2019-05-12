package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rlongo/ambrosia/api"
)

// The SearchRecipes retrieves metadata about all recipes matching the query
func SearchRecipes(storage api.StorageServiceRecipes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		tags := query["tag"]
		author := query.Get("author")

		if recipes, err := storage.GetRecipes(tags, author); err == nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(recipes)
		} else {
			setErrorResponse(w, http.StatusNotFound, err)
		}
	}
}

func GetRecipe(storage api.StorageServiceRecipes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setErrorResponse(w, http.StatusNotFound, fmt.Errorf("Not Implemented"))
	}
}

func AddRecipe(storage api.StorageServiceRecipes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setErrorResponse(w, http.StatusNotFound, fmt.Errorf("Not Implemented"))
	}
}

func UpdateRecipe(storage api.StorageServiceRecipes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setErrorResponse(w, http.StatusNotFound, fmt.Errorf("Not Implemented"))
	}
}
