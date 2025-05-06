package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/michalhercik/RecSIS/auth"
	"github.com/michalhercik/RecSIS/components/bpbtn"
	"github.com/michalhercik/RecSIS/components/page"
	"github.com/michalhercik/RecSIS/components/searchbar"
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
	meiliClient := meilisearch.New(hostMeili, meilisearch.WithAPIKey(searchKey))

	//////////////////////////////////////////
	// Handlers
	//////////////////////////////////////////
	pageTempl := page.Page{
		Home: "/home",
		NavItems: []page.NavItem{
			{Title: language.MakeLangString("Domů", "Home"), Path: "/"},
			{Title: language.MakeLangString("Hledání", "Search"), Path: "/courses/"},
			{Title: language.MakeLangString("Blueprint", "Blueprint"), Path: "/blueprint/"},
			{Title: language.MakeLangString("Studijní plán", "Degree plan"), Path: "/degreeplan/"},
		},
		QuickSearchPath: "/quicksearch",
		SearchBar: searchbar.MeiliSearch{
			Client:            meiliClient,
			Index:             "courses",
			Limit:             5,
			Param:             "search",
			FiltersSelector:   "#filter-form",
			SearchEndpoint:    "/courses",
			QuickEndpoint:     "/page/quicksearch",
			SearchBarView:     searchbar.SearchBar,
			SearchResultsView: searchbar.QuickResults,
			ResultsDetailEndpoint: func(code string) string {
				return fmt.Sprintf("/course/%s", code)
			},
		},
	}
	pageTempl.Init()
	home := home.Server{
		Page: page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
	}
	home.Init()
	blueprint := blueprint.Server{
		Data: blueprint.DBManager{DB: db},
		Auth: auth.UserIDFromContext{},
		Page: page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
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
		Auth: auth.UserIDFromContext{},
		Page: page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
	}
	coursedetail.Init()
	courses := courses.Server{
		Data:   courses.DBManager{DB: db},
		Search: courses.MeiliSearch{Client: meiliClient, Courses: meilisearch.IndexConfig{Uid: "courses"}},
		Auth:   auth.UserIDFromContext{},
		Page:   pageTempl,
		BpBtn: bpbtn.Add{
			DB:    db,
			Templ: bpbtn.AddBtn,
			Options: bpbtn.Options{
				HxPostBase: "/courses",
			},
		},
	}
	courses.Init()
	degreePlan := degreeplan.Server{
		Data: degreeplan.DBManager{DB: db},
		Auth: auth.UserIDFromContext{},
		BpBtn: bpbtn.Add{
			DB:    db,
			Templ: bpbtn.PlusSignBtn,
			Options: bpbtn.Options{
				HxPostBase: "/degreeplan",
			},
		},
		Page: page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
	}
	degreePlan.Init()
	static := http.FileServer(http.Dir("static"))

	router.Handle("/", home.Router())
	handle(router, "/blueprint/", blueprint.Router())
	handle(router, "/course/", coursedetail.Router())
	handle(router, "/courses/", courses.Router())
	handle(router, "/degreeplan/", degreePlan.Router())
	handle(router, "/page/", pageTempl.Router())
	router.Handle("GET /favicon-256x256.ico", static)
	router.Handle("GET /logo.svg", static)
	router.Handle("GET /style.css", static)
	router.Handle("GET /js/", static)

	//////////////////////////////////////////
	// Server setup
	//////////////////////////////////////////
	authentication := auth.Authentication{
		Authenticate: auth.DBManager{DB: db}.Authenticate,
	}
	var handler http.Handler
	handler = router
	handler = language.SetAndStripLanguage(handler)
	// handler = authentication.AuthenticateHTTP(handler)
	_ = authentication
	handler = auth.NoAuth(handler)
	handler = logging(handler)
	server := http.Server{
		Addr:    ":8000", // DOCKER, PRODUCTION: when run as docker container remove localhost
		Handler: handler,
	}

	log.Println("Server starting ...")
	log.Println("http://localhost:8000/")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
