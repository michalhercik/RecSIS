package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/michalhercik/RecSIS/cas"
	meilicomments "github.com/michalhercik/RecSIS/internal/course/comments/meilisearch"
	"github.com/michalhercik/RecSIS/internal/course/comments/meilisearch/params"
	"github.com/michalhercik/RecSIS/internal/course/comments/meilisearch/urlparser"
	"github.com/michalhercik/RecSIS/language"

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

func handle(router *http.ServeMux, prefix string, handler http.Handler) {
	router.Handle(prefix, http.StripPrefix(prefix[:len(prefix)-1], handler))
}

func main() {
	// DANGER: this is a test code, remove it
	//===============================================================================
	// TODO: remove this
	//===============================================================================
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// DANGER
	//===============================================================================

	//////////////////////////////////////////
	// Database setup
	//////////////////////////////////////////

	// Postgres
	const (
		// host     = "postgres" // DOCKER, PRODUCTION: when run as docker container change to network name
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
		// hostMeili = "http://meilisearch:7700" // "http://localhost:7700"
		hostMeili = "http://localhost:7700"
		searchKey = "MASTER_KEY"
	)
	meiliClient := meilisearch.New(hostMeili, meilisearch.WithAPIKey(searchKey))

	//////////////////////////////////////////
	// Handlers
	//////////////////////////////////////////

	home := home.Server{}
	home.Init()
	blueprint := blueprint.Server{
		Data: blueprint.DBManager{DB: db},
		Auth: cas.UserIDFromContext{},
	}
	blueprint.Init()
	coursedetail := coursedetail.Server{
		Data: coursedetail.DBManager{DB: db},
		CourseComments: meilicomments.MeiliSearch{
			Client:        meiliClient,
			CommentsIndex: meilisearch.IndexConfig{Uid: "courses-comments"},
			UrlToFilter: urlparser.FilterParser{
				ParamPrefix: "parf",
				IDToParam:   params.IdToParam,
			},
			TeacherParam: params.TeacherCode,
			CourseParam:  params.CourseCode,
		},
		Auth: cas.UserIDFromContext{},
	}
	coursedetail.Init()
	courses := courses.Server{
		Data:   courses.DBManager{DB: db},
		Search: courses.MeiliSearch{Client: meiliClient, Courses: meilisearch.IndexConfig{Uid: "courses"}},
		Auth:   cas.UserIDFromContext{},
	}
	courses.Init()
	degreePlan := degreeplan.Server{
		Data: degreeplan.DBManager{DB: db},
		Auth: cas.UserIDFromContext{},
	}
	degreePlan.Init()

	static := http.FileServer(http.Dir("static"))

	protectedRouter := http.NewServeMux()
	protectedRouter.Handle("/", home.Router())
	handle(protectedRouter, "/blueprint/", blueprint.Router())
	handle(protectedRouter, "/course/", coursedetail.Router())
	handle(protectedRouter, "/courses/", courses.Router())
	handle(protectedRouter, "/degreeplan/", degreePlan.Router())
	protectedRouter.Handle("GET /logo.svg", static)
	protectedRouter.Handle("GET /style.css", static)

	//////////////////////////////////////////
	// Server setup
	//////////////////////////////////////////
	authentication := cas.Authentication{
		Data: cas.DBManager{DB: db},
		CAS:  cas.CAS{Host: "localhost:8001"},
		// CAS:            cas.CAS{Host: "cas.cuni.cz"},
		AfterLoginPath: "/",
	}
	var handler, unprotectedHandler, protectedHandler http.Handler
	protectedHandler = protectedRouter
	protectedHandler = authentication.AuthenticateHTTP(protectedHandler)
	// handler = auth.NoAuth(handler)
	unprotectedRouter := http.NewServeMux()
	unprotectedRouter.Handle("/", protectedHandler)
	unprotectedRouter.Handle("GET /favicon.ico", static)
	unprotectedRouter.Handle("GET /logo.svg", static)

	unprotectedHandler = unprotectedRouter
	unprotectedHandler = language.SetAndStripLanguage(unprotectedHandler)
	unprotectedHandler = logging(unprotectedHandler)
	handler = unprotectedHandler

	server := http.Server{
		Addr:    ":8000", // DOCKER, PRODUCTION: when run as docker container remove localhost
		Handler: handler,
	}

	log.Println("Server starting ...")
	log.Println("http://localhost:8000/")

	// err = server.ListenAndServeTLS("recsis-cert/fullchain.pem", "recsis-cert/privkey.pem")
	err = server.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
