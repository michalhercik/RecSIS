package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/michalhercik/RecSIS/cas"
	"github.com/michalhercik/RecSIS/components/bpbtn"
	"github.com/michalhercik/RecSIS/components/page"
	"github.com/michalhercik/RecSIS/components/searchbar"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/filters"
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

func main() {
	configPath := flag.String("config", "", "Path to the config file")
	flag.Parse()
	if len(*configPath) == 0 {
		log.Fatal("Config file path is required")
	}

	var conf config
	_, err := toml.DecodeFile(*configPath, &conf)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	switch conf.Environment {
	case productionEnvironment:
		log.Println("INFO: Running in production mode.")
	case developmentEnvironment:
		log.Println("INFO: Running in development mode.")
		// Allow self-signed certificates
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		log.Println("WARNING: Insecure TLS configuration for development environment.")
	default:
		log.Fatalf("Invalid environment: %s", conf.Environment)
	}

	//////////////////////////////////////////
	// Database setup
	//////////////////////////////////////////

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Postgres.Host, conf.Postgres.Port, conf.Postgres.User, conf.Postgres.Password, conf.Postgres.DBName)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	meiliClient := meilisearch.New(conf.MeiliSearch.Host, meilisearch.WithAPIKey(conf.MeiliSearch.Key))
	if !meiliClient.IsHealthy() {
		log.Fatalf("MeiliSearch connection failed")
	}

	//////////////////////////////////////////
	// Error handling
	//////////////////////////////////////////

	errorHandler := errorx.ErrorHandler{
		// Initialize error handler with logging and rendering capabilities
	}

	//////////////////////////////////////////
	// Page template setup
	//////////////////////////////////////////

	pageTempl := page.Page{
		Error: errorHandler,
		Home:  "/home/",
		NavItems: []page.NavItem{
			{Title: language.MakeLangString("Domů", "Home"), Path: "/home/", Skeleton: home.Skeleton, Indicator: "#home-skeleton"},
			{Title: language.MakeLangString("Hledání", "Search"), Path: "/courses/", Skeleton: courses.Skeleton, Indicator: "#courses-skeleton"},
			{Title: language.MakeLangString("Blueprint", "Blueprint"), Path: "/blueprint/", Skeleton: blueprint.Skeleton, Indicator: "#blueprint-skeleton"},
			{Title: language.MakeLangString("Studijní plán", "Degree plan"), Path: "/degreeplan/", Skeleton: degreeplan.Skeleton, Indicator: "#degreeplan-skeleton"},
		},
		QuickSearchPath: "/quicksearch",
		SearchBar: searchbar.MeiliSearch{
			Client:            meiliClient,
			Index:             "courses",
			Limit:             5,
			Param:             "search",
			FiltersSelector:   "#filter-form",
			SearchEndpoint:    "/courses/",
			QuickEndpoint:     "/page/quicksearch",
			SearchBarView:     searchbar.SearchBar,
			SearchResultsView: searchbar.QuickResults,
			ResultsDetailEndpoint: func(code string) string {
				return fmt.Sprintf("/course/%s", code)
			},
		},
	}
	pageTempl.Init()

	errorHandler.Page = page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl}

	//////////////////////////////////////////
	// Handlers
	//////////////////////////////////////////

	home := home.Server{
		Auth:        cas.UserIDFromContext{},
		Error:       errorHandler,
		Page:        page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
		Recommender: fmt.Sprintf("http://%s:%d", conf.Recommender.Host, conf.Recommender.Port),
	}
	home.Init()

	blueprint := blueprint.Server{
		Auth:  cas.UserIDFromContext{},
		Data:  blueprint.DBManager{DB: db},
		Error: errorHandler,
		Page:  page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
	}
	blueprint.Init()

	coursedetail := coursedetail.Server{
		Auth: cas.UserIDFromContext{},
		BpBtn: bpbtn.Add{
			DB:    db,
			Templ: bpbtn.AddBtn,
			Options: bpbtn.Options{
				HxPostBase: "/course",
			},
		},
		Data:  coursedetail.DBManager{DB: db},
		Error: errorHandler,
		Filters: filters.Filters{
			DB:     db,
			Filter: "course-survey",
		},
		Page: page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
		Search: coursedetail.Search{
			Client: meiliClient,
			Survey: meilisearch.IndexConfig{Uid: "survey"},
		},
	}
	coursedetail.Init()

	courses := courses.Server{
		Auth: cas.UserIDFromContext{},
		BpBtn: bpbtn.Add{
			DB:    db,
			Templ: bpbtn.AddBtn,
			Options: bpbtn.Options{
				HxPostBase: "/courses",
			},
		},
		Data:  courses.DBManager{DB: db},
		Error: errorHandler,
		Filters: filters.Filters{
			DB:     db,
			Filter: "courses",
		},
		Page: pageTempl,
		Search: courses.MeiliSearch{
			Client:  meiliClient,
			Courses: meilisearch.IndexConfig{Uid: "courses"},
		},
	}
	courses.Init()

	degreePlan := degreeplan.Server{
		Auth: cas.UserIDFromContext{},
		BpBtn: bpbtn.DoubleAdd{
			Add: bpbtn.Add{
				DB:    db,
				Templ: bpbtn.PlusSignBtn,
				Options: bpbtn.Options{
					HxPostBase: "/degreeplan",
				},
			},
			TemplSecond: bpbtn.PlusSignBtnChecked,
		},
		Data: degreeplan.DBManager{DB: db},
		DPSearch: degreeplan.MeiliSearch{
			Client:      meiliClient,
			DegreePlans: meilisearch.IndexConfig{Uid: "degree-plans"},
		},
		Error: errorHandler,
		Page:  page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
	}
	degreePlan.Init()

	exePath, err := os.Executable()
	static := http.FileServer(http.Dir(filepath.Join(filepath.Dir(exePath), "static")))

	protectedRouter := http.NewServeMux()
	protectedRouter.Handle("/", home.Router())
	handle(protectedRouter, "/page/", pageTempl.Router())
	handle(protectedRouter, "/blueprint/", blueprint.Router())
	handle(protectedRouter, "/course/", coursedetail.Router())
	handle(protectedRouter, "/courses/", courses.Router())
	handle(protectedRouter, "/degreeplan/", degreePlan.Router())
	protectedRouter.Handle("GET /logo.svg", static)
	protectedRouter.Handle("GET /style.css", static)
	protectedRouter.Handle("GET /js/", static)

	//////////////////////////////////////////
	// Server setup
	//////////////////////////////////////////

	authentication := cas.Authentication{
		Data:           cas.DBManager{DB: db},
		Error:          errorHandler,
		CAS:            cas.CAS{Host: conf.CAS.Host},
		AfterLoginPath: "/",
	}
	var handler, unprotectedHandler, protectedHandler http.Handler
	protectedHandler = protectedRouter
	protectedHandler = authentication.AuthenticateHTTP(protectedHandler)

	unprotectedRouter := http.NewServeMux()
	unprotectedRouter.Handle("/", protectedHandler)
	unprotectedRouter.Handle("GET /favicon.ico", static)
	unprotectedRouter.Handle("GET /logo.svg", static)

	unprotectedHandler = unprotectedRouter
	unprotectedHandler = language.SetAndStripLanguage(unprotectedHandler)
	unprotectedHandler = logging(unprotectedHandler)
	handler = unprotectedHandler

	// Redirect http to https
	go func() {
		httpServer := http.Server{
			Addr:    fmt.Sprintf(":%d", conf.Server.HTTP.Port),
			Handler: http.HandlerFunc(redirectToTLS(conf.Server.HTTPS.Port)),
		}
		if err = httpServer.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	log.Println("Server starting ...")
	if conf.Environment == developmentEnvironment {
		log.Printf("https://localhost:%d/", conf.Server.HTTPS.Port)
	}
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.HTTPS.Port),
		Handler: handler,
	}
	err = server.ListenAndServeTLS(conf.SSL.Certificate, conf.SSL.Key)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func redirectToTLS(port int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("https://%s:%d%s", r.Host, port, r.URL.String())
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}

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

const (
	productionEnvironment  = "production"
	developmentEnvironment = "development"
)

type config struct {
	Environment string `toml:"environment"`
	Server      struct {
		HTTP struct {
			Port int `toml:"port"`
		} `toml:"http"`
		HTTPS struct {
			Port int `toml:"port"`
		} `toml:"https"`
	} `toml:"server"`
	Postgres struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		DBName   string `toml:"dbname"`
	} `toml:"postgres"`
	MeiliSearch struct {
		Host string `toml:"host"`
		Key  string `toml:"key"`
	} `toml:"meilisearch"`
	Recommender struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
	} `toml:"recommender"`
	CAS struct {
		Host string `toml:"host"`
	} `toml:"cas"`
	SSL struct {
		Certificate string `toml:"certificate"`
		Key         string `toml:"key"`
	} `toml:"ssl"`
}
