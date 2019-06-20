package app

import (
	"net/http"

	"github.com/rlongo/ambrosia-server-go/api"
	"github.com/rlongo/ambrosia-server-go/app/handler"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(api.StorageServiceRecipes) http.HandlerFunc
}

type routes []route

var ambrosiaRoutes = routes{
	route{"Recipes List", http.MethodGet, "/recipes", handler.SearchRecipes},
	route{"Single Recipe", http.MethodGet, "/recipe/{id}", handler.GetRecipe},
	route{"Add Recipe", http.MethodPost, "/recipe", handler.AddRecipe},
	route{"Update Recipe", http.MethodPut, "/recipe/{id}", handler.UpdateRecipe},
}
