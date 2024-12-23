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
	// template
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		if r.Method != "HEAD" {
			log.Println(r.Method, r.URL.Path)
		}
	})
}

func main() {

	router := http.NewServeMux()

	//////////////////////////////////////////
	// Database setup
	//////////////////////////////////////////

	const (
		host     = "postgres"
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

	//////////////////////////////////////////
	// Handlers
	//////////////////////////////////////////

	home.Server{}.Register(router)
	blueprint.Server{
		Data: blueprint.DBManager{DB: db},
	}.Register(router, "/blueprint")
	coursedetail.Server{
		Data: coursedetail.DBManager{DB: db},
	}.Register(router, "/course")
	courses.Server{
		Data: courses.DBManager{DB: db},
	}.Register(router, "/courses")
	degreeplan.Server{}.Register(router, "/degreeplan")

	// Icon
	router.Handle("/favicon.ico", http.FileServer(http.Dir(".")))

	// Styles
	router.Handle("GET /style.css", http.FileServer(http.Dir("utils/")))

	//////////////////////////////////////////
	// Server setup
	//////////////////////////////////////////
	server := http.Server{
		Addr:    ":8000", // when run as docker container remove localhost
		Handler: logging(router),
	}

	log.Println("Server starting ...")
	log.Println("http://localhost:8000/")

	server.ListenAndServe()
}
