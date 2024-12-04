package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	// pages
	"github.com/michalhercik/RecSIS/blueprint"
	"github.com/michalhercik/RecSIS/coursedetail"
	"github.com/michalhercik/RecSIS/courses"
	"github.com/michalhercik/RecSIS/degreeplan"
	"github.com/michalhercik/RecSIS/home"

	// database
	_ "github.com/lib/pq"
	"github.com/michalhercik/RecSIS/mockdb"

	// TODO potentially import real database

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

type generator func(http.ResponseWriter, *http.Request) templ.Component

func htmxRouter(page generator, content generator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		component := page(w, r)
		if r.Header.Get("hx-request") == "true" {
			component = content(w, r)
		}
		component.Render(r.Context(), w)
	}
}

func main() {

	router := http.NewServeMux()

	//////////////////////////////////////////
	// Database setup
	//////////////////////////////////////////

	const (
		host     = "localhost"
		port     = 5432
		user     = "recsis"
		password = "recsis"
		dbname   = "recsis"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	mockDB := mockdb.NewDB()
	blueprint.SetDataManager(blueprint.DBManager{DB: db})
	coursedetail.SetDatabase(coursedetail.DBCourseReader{DB: db})
	courses.SetDatabase(courses.CreateDB(mockDB))

	//////////////////////////////////////////
	// Handlers
	//////////////////////////////////////////

	// Icon
	router.Handle("/favicon.ico", http.FileServer(http.Dir(".")))

	// Home
	router.HandleFunc("GET /{$}", htmxRouter(home.HandlePage, home.HandleContent))
	router.HandleFunc("GET /home", htmxRouter(home.HandlePage, home.HandleContent))

	// Courses
	router.HandleFunc("GET /courses", htmxRouter(courses.HandlePage, courses.HandleContent))

	// Degree plan
	router.HandleFunc("GET /degreeplan", htmxRouter(degreeplan.HandlePage, degreeplan.HandleContent))

	// Blueprint
	router.HandleFunc("GET /blueprint", htmxRouter(blueprint.HandlePage, blueprint.HandleContent))
	router.HandleFunc("POST /blueprint/year", blueprint.HandleYearAddition)
	router.HandleFunc("DELETE /blueprint/year", blueprint.HandleYearRemoval)
	router.HandleFunc("DELETE /blueprint/{year}/{semester}/{code}", blueprint.HandleCourseRemoval)

	// Course detail
	router.HandleFunc("GET /course/{code}", htmxRouter(coursedetail.HandlePage, coursedetail.HandleContent))
	router.HandleFunc("POST /course/{code}/add-comment", coursedetail.HandleAddingComment)
	router.HandleFunc("POST /course/{code}/like", coursedetail.HandleLike)
	router.HandleFunc("POST /course/{code}/dislike", coursedetail.HandleDislike)

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
