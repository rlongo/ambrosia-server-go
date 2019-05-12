package handler

import (
	"fmt"
	"net/http"

	"github.com/rlongo/ambrosia/api"
)

// The SearchRecipes retrieves metadata about all recipes matching the query
func SearchRecipes(storage api.StorageServiceRecipes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setErrorResponse(w, http.StatusNotFound, fmt.Errorf("Not Implemented"))
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
