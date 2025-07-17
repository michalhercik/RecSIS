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
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/filters"
	"github.com/michalhercik/RecSIS/language"

	"github.com/michalhercik/RecSIS/blueprint"
	"github.com/michalhercik/RecSIS/coursedetail"
	"github.com/michalhercik/RecSIS/courses"
	"github.com/michalhercik/RecSIS/degreeplan"
	"github.com/michalhercik/RecSIS/home"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/meilisearch/meilisearch-go"
)

func main() {
	configPath := configPath()
	conf := configFrom(configPath)
	handler := setupHandler(conf)

	log.Println("Server starting ...")
	log.Printf("port: %d", conf.Server.HTTPS.Port)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.HTTPS.Port),
		Handler: handler,
	}
	err := server.ListenAndServeTLS(conf.SSL.Certificate, conf.SSL.Key)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func configPath() string {
	configPath := flag.String("config", "", "Path to the config file")
	flag.Parse()
	return *configPath
}

func configFrom(configPath string) config {
	if len(configPath) == 0 {
		log.Fatal("Config file path is required")
	}

	var conf config
	_, err := toml.DecodeFile(configPath, &conf)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	switch conf.Environment {
	case productionEnvironment:
		log.Println("INFO: Running in production mode.")
	case developmentEnvironment:
		log.Println("INFO: Running in development mode.")
		// Allow communication with servers with self-signed certificates e.g. Mock CAS
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		log.Println("WARNING: Insecure TLS configuration for development environment.")
	default:
		log.Fatalf("Invalid environment: %s", conf.Environment)
	}

	return conf
}

func setupHandler(conf config) http.Handler {
	db := setupDB(conf)
	meiliClient := meiliServiceManager(conf)

	errorHandler := errorx.ErrorHandler{}

	pageTempl := pageTemplate(errorHandler, meiliClient)

	errorHandler.Page = page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl}

	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	s := servers{
		pageTempl:          pageTempl.Router(),
		homeServer:         homeServer(conf, errorHandler, pageTempl),
		blueprintServer:    blueprintServer(db, errorHandler, pageTempl),
		coursedetailServer: courseDetailServer(db, errorHandler, pageTempl, meiliClient),
		coursesServer:      coursesServer(db, errorHandler, pageTempl, meiliClient),
		degreePlanServer:   degreePlanServer(db, errorHandler, pageTempl, meiliClient),
		static:             http.FileServer(http.Dir(filepath.Join(filepath.Dir(exePath), "static"))),
	}
	handler := protectedHandler(s)
	handler = authenticationHandler(handler, db, errorHandler, conf)
	handler = unprotectedHandler(handler, s.static)
	return handler
}

func setupDB(conf config) *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Postgres.Host, conf.Postgres.Port, conf.Postgres.User, conf.Postgres.Password, conf.Postgres.DBName)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	return db
}

func meiliServiceManager(conf config) meilisearch.ServiceManager {
	ms := meilisearch.New(conf.MeiliSearch.Host, meilisearch.WithAPIKey(conf.MeiliSearch.Key))
	if !ms.IsHealthy() {
		log.Fatalf("MeiliSearch connection failed")
	}
	return ms
}

func pageTemplate(errorHandler page.Error, meiliClient meilisearch.ServiceManager) page.Page {
	pageTempl := page.Page{
		Error: errorHandler,
		Home:  homeRoot,
		NavItems: []page.NavItem{
			{Title: language.MakeLangString("Domů", "Home"), Path: homeRoot, Skeleton: home.Skeleton, Indicator: "#home-skeleton"},
			{Title: language.MakeLangString("Hledání", "Search"), Path: coursesRoot, Skeleton: courses.Skeleton, Indicator: "#courses-skeleton"},
			{Title: language.MakeLangString("Blueprint", "Blueprint"), Path: blueprintRoot, Skeleton: blueprint.Skeleton, Indicator: "#blueprint-skeleton"},
			{Title: language.MakeLangString("Studijní plán", "Degree plan"), Path: degreePlanRoot, Skeleton: degreeplan.Skeleton, Indicator: "#degreeplan-skeleton"},
		},
		Search: page.MeiliSearch{
			Client: meiliClient,
			Index:  "courses",
			Limit:  5,
		},
		Param:          "search",
		SearchEndpoint: coursesRoot,
		ResultsDetailEndpoint: func(code string) string {
			return courseDetailRoot + code
		},
	}
	pageTempl.Init()

	return pageTempl
}

func homeServer(conf config, errorHandler home.Error, pageTempl page.Page) http.Handler {
	home := home.Server{
		Auth:        cas.UserIDFromContext{},
		Error:       errorHandler,
		Page:        page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
		Recommender: fmt.Sprintf("http://%s:%d", conf.Recommender.Host, conf.Recommender.Port),
	}
	home.Init()
	return home.Router()
}

func blueprintServer(db *sqlx.DB, errorHandler blueprint.Error, pageTempl page.Page) http.Handler {
	blueprint := blueprint.Server{
		Auth:  cas.UserIDFromContext{},
		Data:  blueprint.DBManager{DB: db},
		Error: errorHandler,
		Page:  page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
	}
	blueprint.Init()
	return blueprint.Router()
}

func courseDetailServer(db *sqlx.DB, errorHandler coursedetail.Error, pageTempl page.Page, meiliClient meilisearch.ServiceManager) http.Handler {
	coursedetail := coursedetail.Server{
		Auth: cas.UserIDFromContext{},
		BpBtn: bpbtn.Add{
			DB:    db,
			Templ: bpbtn.AddBtn,
			Options: bpbtn.Options{
				HxPostBase: courseDetailRoot,
			},
		},
		Data:    coursedetail.DBManager{DB: db},
		Error:   errorHandler,
		Filters: filters.MakeFilters(db, "course-survey"),
		Page:    page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
		Search: coursedetail.Search{
			Client: meiliClient,
			Survey: meilisearch.IndexConfig{Uid: "survey"},
		},
	}
	coursedetail.Init()
	return coursedetail.Router()
}

func coursesServer(db *sqlx.DB, errorHandler courses.Error, pageTempl page.Page, meiliClient meilisearch.ServiceManager) http.Handler {
	courses := courses.Server{
		Auth: cas.UserIDFromContext{},
		BpBtn: bpbtn.Add{
			DB:    db,
			Templ: bpbtn.AddBtn,
			Options: bpbtn.Options{
				HxPostBase: coursesRoot,
			},
		},
		Data:    courses.DBManager{DB: db},
		Error:   errorHandler,
		Filters: filters.MakeFilters(db, "courses"),
		Page:    pageTempl,
		Search: courses.MeiliSearch{
			Client:  meiliClient,
			Courses: meilisearch.IndexConfig{Uid: "courses"},
		},
	}
	courses.Init()
	return courses.Router()
}

func degreePlanServer(db *sqlx.DB, errorHandler degreeplan.Error, pageTempl page.Page, meiliClient meilisearch.ServiceManager) http.Handler {
	degreePlan := degreeplan.Server{
		Auth: cas.UserIDFromContext{},
		BpBtn: bpbtn.AddWithTwoTemplComponents{
			Add: bpbtn.Add{
				DB:    db,
				Templ: bpbtn.PlusSignBtn,
				Options: bpbtn.Options{
					HxPostBase: degreePlanRoot,
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
	return degreePlan.Router()
}

func protectedHandler(s servers) http.Handler {
	protectedRouter := http.NewServeMux()
	protectedRouter.Handle(homeRoot, s.homeServer)
	handle(protectedRouter, pageRoot, s.pageTempl)
	handle(protectedRouter, blueprintRoot, s.blueprintServer)
	handle(protectedRouter, courseDetailRoot, s.coursedetailServer)
	handle(protectedRouter, coursesRoot, s.coursesServer)
	handle(protectedRouter, degreePlanRoot, s.degreePlanServer)
	protectedRouter.Handle("GET /logo.svg", s.static)
	protectedRouter.Handle("GET /style.css", s.static)
	protectedRouter.Handle("GET /js/", s.static)

	return protectedRouter
}

func authenticationHandler(prev http.Handler, db *sqlx.DB, errorHandler cas.Error, conf config) http.Handler {
	authentication := cas.Authentication{
		Data:           cas.DBManager{DB: db},
		Error:          errorHandler,
		CAS:            cas.CAS{Host: conf.CAS.Host},
		AfterLoginPath: homeRoot,
	}
	var authenticationHandler http.Handler
	authenticationHandler = prev
	authenticationHandler = authentication.AuthenticateHTTP(authenticationHandler)

	return authenticationHandler
}

func unprotectedHandler(prev http.Handler, static http.Handler) http.Handler {
	unprotectedRouter := http.NewServeMux()
	unprotectedRouter.Handle("/", prev)
	unprotectedRouter.Handle("GET /favicon.ico", static)
	unprotectedRouter.Handle("GET /logo.svg", static)
	unprotectedRouter.Handle("GET /help/", static)

	var unprotectedHandler http.Handler
	unprotectedHandler = unprotectedRouter
	unprotectedHandler = language.SetAndStripLanguageHandler(unprotectedHandler)
	unprotectedHandler = logging(unprotectedHandler)
	return unprotectedHandler
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path)
	})
}

func handle(router *http.ServeMux, prefix string, handler http.Handler) {
	router.Handle(prefix, http.StripPrefix(prefix[:len(prefix)-1], handler))
}

type servers struct {
	homeServer         http.Handler
	pageTempl          http.Handler
	blueprintServer    http.Handler
	coursedetailServer http.Handler
	coursesServer      http.Handler
	degreePlanServer   http.Handler
	static             http.Handler
}

const (
	pageRoot         = "/page/"
	homeRoot         = "/"
	blueprintRoot    = "/blueprint/"
	courseDetailRoot = "/course/"
	coursesRoot      = "/courses/"
	degreePlanRoot   = "/degreeplan/"
)

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
