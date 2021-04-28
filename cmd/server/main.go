package main

import (
	"fmt"
	"github.com/namsral/flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/marcsj/standardnotes-extensions/controller"
	"github.com/marcsj/standardnotes-extensions/definition"
)

var (
	listenPort = flag.Int(
		"listen_port",
		80,
		"main listening port")
	baseURL = flag.String(
		"base_url",
		"https://extensions.your.domain",
		"the base url for this service")
	updateInterval = flag.Int(
		"update_interval",
		4320,
		"update interval in minutes")
	reposDir = flag.String(
		"repos_dir",
		"/repos",
		"directory to pull repositories from")
	definitionsDir = flag.String(
		"def_dir",
		"/definitions",
		"directory to pull definitions from")

)

func main() {
	flag.Parse()

	ctrl := &controller.Controller{
		BaseURL:        *baseURL,
		ReposDir:       *reposDir,
		DefinitionsDir: *definitionsDir,
		Packages:       map[string]*definition.Package{},
	}
	go updatePackages(ctrl.UpdatePackages, time.Duration(*updateInterval)*time.Minute)
	//
	r := mux.NewRouter()
	r.HandleFunc("/index.json", ctrl.ServeIndex)
	r.HandleFunc("/{id}/index.json", ctrl.ServePackageIndex)
	r.PathPrefix("/{id}/{version}/").HandlerFunc(ctrl.ServePackage)
	withCORS := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"DNT", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Range"}),
		handlers.MaxAge(1728000),
	)
	withLogging := handlers.LoggingHandler
	server := http.Server{
		Addr:    fmt.Sprintf(":%v", *listenPort),
		Handler: withLogging(os.Stdout, withCORS(r)),
	}
	log.Printf("listening on %v", *listenPort)
	log.Fatalf("error starting server: %v", server.ListenAndServe())
}

func updatePackages(fn func() error, wait time.Duration) {
	for {
		if err := fn(); err != nil {
			log.Printf("error updating packages: %v", err)
		}
		time.Sleep(wait)
		log.Print("finished updating packages")
	}
}
