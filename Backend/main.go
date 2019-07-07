package main

import (
	"net/http"
	"os"
	"fmt"
	"log"
	"strings"
	"github.com/go-pg/pg"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-chi/chi/middleware"
	"github.com/alistairfink/Steak/Backend/Utilities"
	"github.com/alistairfink/Steak/Backend/Data/DatabaseConnection"
	"github.com/alistairfink/Steak/Backend/Domain/Middleware"
)

func main() {
	// Read Config
	var config *Utilities.Config
	if _, err := os.Stat("./config.json"); err == nil {
		config = Utilities.GetConfig(".", "config")
	} else {
		config = Utilities.GetConfig("/go/src/github.com/alistairfink/Steak/.", "config")
	}

	// Open DB
	db := DatabaseConnection.Connect(config)
	defer DatabaseConnection.Close(db)
	
	// Router
	router := Routes(db, config)
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf(" %-10s%-10s\n", method, strings.Replace(route, "/*", "", -1))
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}

	log.Fatal(http.ListenAndServe(":" + config.Port, router))
}

func Routes(db *pg.DB, config *Utilities.Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use (
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		Middleware.CorsMiddleware,
	)

	// Controllers

	// Paths
	router.Route("/", func(routes chi.Router) {
	})

	return router
}