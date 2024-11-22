package main

import (
	"log"
	"net/http"

	// pages
	"github.com/michalhercik/RecSIS/blueprint"
	"github.com/michalhercik/RecSIS/courses"
	"github.com/michalhercik/RecSIS/degree_plan"
	"github.com/michalhercik/RecSIS/home"

	// database
	"github.com/michalhercik/RecSIS/database"
	"github.com/michalhercik/RecSIS/database/mock"
	//"github.com/michalhercik/RecSIS/database/real"

	// template
	"github.com/a-h/templ"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		if r.Method != "HEAD" {
			log.Println(r.Method, r.URL.Path)
		}
	})
}

type generator func() templ.Component

func htmxRouter(page generator, content generator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		component := page()
		if r.Header.Get("hx-request") == "true" {
			component = content()
		}
		component.Render(r.Context(), w)
	}
}

func main() {

	router := http.NewServeMux()

	//////////////////////////////////////////
	// Database setup
	//////////////////////////////////////////

	database.SetDatabase(mock.NewMockDB())
	//database.SetDatabase(real.NewRealDB())

	//////////////////////////////////////////
	// Handlers
	//////////////////////////////////////////

	router.HandleFunc("GET /", htmxRouter(home.Page, home.Content))
	router.HandleFunc("GET /home", htmxRouter(home.Page, home.Content))
	router.HandleFunc("GET /courses", htmxRouter(courses.HandlePage, courses.HandleContent))
	router.HandleFunc("GET /degree_plan", htmxRouter(degree_plan.Page, degree_plan.Content))

	// Blueprint
	router.HandleFunc("GET /blueprint", htmxRouter(blueprint.HandlePage, blueprint.HandleContent))
	//router.HandleFunc("POST /blueprint/remove-year/{year}", blueprint.HandleLastYearRemoval)
	router.HandleFunc("DELETE /blueprint/remove-unassigned/{id}", blueprint.HandleBLueprintUnassignedRemoval)

	//////////////////////////////////////////
	// Server setup
	//////////////////////////////////////////
	server := http.Server{
		Addr:    "localhost:8000", // when run as docker container remove localhost
		Handler: logging(router),
	}

	log.Println("Server starting ...")
	log.Println("http://localhost:8000/")

	server.ListenAndServe()
}
