package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rlongo/ambrosia/app"
	"github.com/rlongo/ambrosia/storage"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

func getOSVariable(key string) string {
	if v := os.Getenv(key); len(v) > 0 {
		return v
	}

	panic(fmt.Sprintf("env %s isn't set!", key))
}

func main() {
	db := getOSVariable("DB_URL")
	port := getOSVariable("PORT")

	storageService, err := storage.OpenMongo(db, storage.RecipesCollection)
	if err != nil {
		log.Fatal(err)
	}

	defer storageService.Close()

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		Debug:            true,
	})

	middleware := negroni.Classic()
	middleware.Use(corsConfig)

	router := app.NewRouter(storageService)
	middleware.UseHandler(router)

	middleware.Run(":" + port)
}
