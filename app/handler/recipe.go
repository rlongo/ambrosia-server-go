package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	guid "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rlongo/ambrosia-server-go/api"
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
		params := mux.Vars(r)
		var recipeID = params["id"]
		fmt.Printf("Asked for recipe: %s\n", recipeID)
		if recipeID, err := guid.Parse(recipeID); err == nil {
			if r, err := storage.GetRecipe(api.RecipeID(recipeID)); err == nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(r)
			} else {
				setErrorResponse(w, http.StatusNotFound, err)
			}
		} else {
			setErrorResponse(w, http.StatusNotFound, fmt.Errorf("ID Wasn't Valid"))
		}
	}
}

func AddRecipe(storage api.StorageServiceRecipes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recipe api.Recipe

		if r.Body == nil {
			setErrorResponse(w, http.StatusBadRequest, fmt.Errorf("Can't process an empty response"))
			return
		}

		// Read out portion of body to avoid spurious crashes
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, PostMaxSize))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		if err := json.Unmarshal(body, &recipe); err != nil {
			setErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		if err := storage.AddRecipe(&recipe); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
		} else {
			setErrorResponse(w, http.StatusBadRequest, err)
		}
	}
}

func UpdateRecipe(storage api.StorageServiceRecipes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recipe api.Recipe

		if r.Body == nil {
			setErrorResponse(w, http.StatusBadRequest, fmt.Errorf("Can't process an empty response"))
			return
		}

		// Read out portion of body to avoid spurious crashes
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, PostMaxSize))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		if err := json.Unmarshal(body, &recipe); err != nil {
			setErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		if err := storage.UpdateRecipe(&recipe); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
		} else {
			setErrorResponse(w, http.StatusBadRequest, err)
		}
	}
}
