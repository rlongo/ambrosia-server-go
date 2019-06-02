package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	corsConfig := cors.Default()

	middleware := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		corsConfig,
	)
	router := app.NewRouter(storageService, middleware)

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
