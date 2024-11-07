package main

import (
	"log"
	"net/http"

	"github.com/michalhercik/RecSIS/blueprint"
	"github.com/michalhercik/RecSIS/courses"
	"github.com/michalhercik/RecSIS/degree_plan"
	"github.com/michalhercik/RecSIS/home"

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
	// Handlers
	//////////////////////////////////////////

	router.HandleFunc("GET /", htmxRouter(home.Page, home.Content))
	router.HandleFunc("GET /home", htmxRouter(home.Page, home.Content))
	router.HandleFunc("GET /courses", htmxRouter(courses.HandlePage, courses.HandleContent))
	router.HandleFunc("GET /blueprint", htmxRouter(blueprint.HandlePage, blueprint.HandleContent))
	router.HandleFunc("GET /degree_plan", htmxRouter(degree_plan.Page, degree_plan.Content))

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
