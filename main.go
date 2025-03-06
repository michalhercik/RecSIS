package main

import (
	"fmt"
	"log"
	"net/http"

	// pages
	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/blueprint"
	"github.com/michalhercik/RecSIS/coursedetail"
	"github.com/michalhercik/RecSIS/courses"
	"github.com/michalhercik/RecSIS/degreeplan"
	"github.com/michalhercik/RecSIS/home"

	// database
	_ "github.com/lib/pq"
	"github.com/meilisearch/meilisearch-go"
	// template
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		// PRODUCTION: remove condition -> log everything
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

	// Postgres
	const (
		host     = "localhost" // DOCKER, PRODUCTION: when run as docker container change to network name
		port     = 5432
		user     = "recsis"
		password = "recsis"
		dbname   = "recsis"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// MeiliSearch
	const (
		hostMeili = "http://localhost:7700"
		searchKey = "MASTER_KEY"
	)
	meiliClient := meilisearch.New("http://localhost:7700", meilisearch.WithAPIKey(searchKey))

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
		Data:   courses.DBManager{DB: db},
		Search: courses.MeiliSearch{Client: meiliClient},
	}.Register(router, "/courses")
	degreeplan.Server{
		Data: degreeplan.DBManager{DB: db},
	}.Register(router, "/degreeplan")

	static := http.FileServer(http.Dir("static"))
	router.Handle("/favicon.ico", static)
	router.Handle("GET /style.css", static)

	//////////////////////////////////////////
	// Server setup
	//////////////////////////////////////////
	server := http.Server{
		Addr:    "localhost:8000", // DOCKER, PRODUCTION: when run as docker container remove localhost
		Handler: logging(router),
	}

	log.Println("Server starting ...")
	log.Println("http://localhost:8000/")

	server.ListenAndServe()
}
