package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rlongo/ambrosia/api"
	"github.com/rlongo/ambrosia/app"
	"github.com/rlongo/ambrosia/storage/memory"
	"github.com/urfave/negroni"
)

func getOSVariable(key string) string {
	if v := os.Getenv(key); len(v) > 0 {
		return v
	}

	panic(fmt.Sprintf("env %s isn't set!", key))
}

func main() {

	port := getOSVariable("PORT")

	recipesDB := api.Recipes{}

	storageService := memory.AmbrosiaStorageMemory{RecipesDB: recipesDB}
	middleware := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)
	router := app.NewRouter(&storageService, middleware)

	log.Printf("listening on IPv4 address \"0.0.0.0\", port %s", port)
	log.Printf("listening on IPv6 address \"::\", port %s", port)

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		WriteTimeout:   time.Second * 15,
		ReadTimeout:    time.Second * 15,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
