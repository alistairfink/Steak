package main

import (
	"fmt"
	"github.com/alistairfink/Steak/Backend/Data/DatabaseConnection"
	"github.com/alistairfink/Steak/Backend/Domain/Controllers"
	"github.com/alistairfink/Steak/Backend/Domain/Middleware"
	"github.com/alistairfink/Steak/Backend/Utilities"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-pg/pg"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
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
	localAddrs, _ := net.InterfaceAddrs()
	ip, _ := localAddrs[1].(*net.IPNet)
	println("=============================== Serving On ===============================")
	fmt.Printf(" %-12s%-12s\n", "Local", "localhost:"+config.Port)
	fmt.Printf(" %-12s%-12s\n", "Network", ip.IP.String()+":"+config.Port)
	println("==========================================================================\n")
	router := Routes(db, config)
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf(" %-10s%-10s\n", method, strings.Replace(route, "/*", "", -1))
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}

	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

func Routes(db *pg.DB, config *Utilities.Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		Middleware.CorsMiddleware,
	)

	// Controllers
	recipeController := Controllers.NewRecipeController(db, config)

	// Paths
	router.Route("/", func(routes chi.Router) {
		routes.Mount("/recipe", recipeController.Routes())
	})

	return router
}
