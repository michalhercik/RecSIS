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

	blueprint.SetDataManager(blueprint.DBManager{DB: db})
	coursedetail.SetDataManager(coursedetail.DBManager{DB: db})
	courses.SetDataManager(courses.DBManager{DB: db})

	//////////////////////////////////////////
	// Handlers
	//////////////////////////////////////////

	// Icon
	router.Handle("/favicon.ico", http.FileServer(http.Dir(".")))

	// Styles
	router.Handle("GET /style.css", http.FileServer(http.Dir("utils/")))

	// Home
	router.HandleFunc("GET /{$}", home.HandlePage)
	router.HandleFunc("GET /home", home.HandlePage)

	// Courses
	router.HandleFunc("GET /courses", courses.HandlePage)
	router.HandleFunc("GET /courses/page", courses.HandlePaging)
	router.HandleFunc("GET /courses/search", courses.HandleSearch)
	router.HandleFunc("POST /courses/blueprint/{code}", courses.HandleBlueprintAddition)
	router.HandleFunc("DELETE /courses/blueprint/{year}/{semester}/{code}", courses.HandleBlueprintRemoval)

	// Degree plan
	router.HandleFunc("GET /degreeplan", degreeplan.HandlePage)

	// Blueprint
	router.HandleFunc("GET /blueprint", blueprint.HandlePage)
	router.HandleFunc("POST /blueprint/year", blueprint.HandleYearAddition)
	router.HandleFunc("DELETE /blueprint/year", blueprint.HandleYearRemoval)
	router.HandleFunc("POST /blueprint/course/{code}", blueprint.HandleCourseAddition)
	router.HandleFunc("PATCH /blueprint/course/{id}", blueprint.HandleCourseMovement)
	router.HandleFunc("DELETE /blueprint/course/{id}", blueprint.HandleCourseRemoval)

	// Course detail
	router.HandleFunc("GET /course/{code}", coursedetail.HandlePage)
	router.HandleFunc("POST /course/{code}/comment", coursedetail.HandleCommentAddition)
	router.HandleFunc("POST /course/{code}/like", coursedetail.HandleLike)
	router.HandleFunc("POST /course/{code}/dislike", coursedetail.HandleDislike)
	router.HandleFunc("POST /course/{code}/blueprint", coursedetail.HandleBlueprintAddition)
	router.HandleFunc("DELETE /course/{code}/blueprint", coursedetail.HandleBlueprintRemoval)

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
