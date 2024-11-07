package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path)
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

	router.HandleFunc("GET /", htmxRouter(homePage, homeContent))
	router.HandleFunc("GET /courses", htmxRouter(coursesPage, coursesContent))

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
