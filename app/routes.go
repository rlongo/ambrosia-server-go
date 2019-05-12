package app

import (
	"net/http"

	"github.com/rlongo/ambrosia/api"
	"github.com/rlongo/ambrosia/app/handler"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(api.StorageServiceRecipes) http.HandlerFunc
}

type routes []route

var ambrosiaRoutes = routes{
	route{"Recipes List", "GET", "/recipes", handler.SearchRecipes},
	route{"Single Recipe", "GET", "/recipe/{id}", handler.GetRecipe},
	route{"Add Recipe", "POST", "/recipe", handler.AddRecipe},
	route{"Update Recipe", "PUT", "/recipe/{id}", handler.UpdateRecipe},
}
