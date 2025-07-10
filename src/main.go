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
	configPath := getConfigPath()
	conf := getConfig(configPath)
	handler := setupHandler(conf)

	//////////////////////////////////////////
	// Server setup
	//////////////////////////////////////////

	var err error
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

func getConfigPath() string {
	configPath := flag.String("config", "", "Path to the config file")
	flag.Parse()
	return *configPath
}

func getConfig(configPath string) config {
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
		// Allow self-signed certificates
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		log.Println("WARNING: Insecure TLS configuration for development environment.")
	default:
		log.Fatalf("Invalid environment: %s", conf.Environment)
	}

	return conf
}

func setupHandler(conf config) http.Handler {
	db := setupDB(conf)
	meiliClient := setupMeiliSearch(conf)

	errorHandler := setupErrorHandler()

	pageTempl := setupPage(conf, errorHandler, meiliClient)

	errorHandler.Page = page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl}

	homeServer := setupHomeServer(conf, errorHandler, pageTempl)
	blueprintServer := setupBlueprintServer(db, errorHandler, pageTempl)
	coursedetailServer := setupCourseDetailServer(db, errorHandler, pageTempl, meiliClient)
	coursesServer := setupCoursesServer(db, errorHandler, pageTempl, meiliClient)
	degreePlanServer := setupDegreePlanServer(db, errorHandler, pageTempl, meiliClient)

	exePath, _ := os.Executable()
	static := http.FileServer(http.Dir(filepath.Join(filepath.Dir(exePath), "static")))

	protectedHandler := setupProtectedHandler(homeServer, pageTempl, blueprintServer, coursedetailServer, coursesServer, degreePlanServer, static, db, errorHandler, conf)
	unprotectedHandler := setupUnprotectedHandler(protectedHandler, static)

	return unprotectedHandler
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

func setupMeiliSearch(conf config) meilisearch.ServiceManager {
	ms := meilisearch.New(conf.MeiliSearch.Host, meilisearch.WithAPIKey(conf.MeiliSearch.Key))
	if !ms.IsHealthy() {
		log.Fatalf("MeiliSearch connection failed")
	}
	return ms
}

func setupErrorHandler() errorx.ErrorHandler {
	errorHandler := errorx.ErrorHandler{
		// Initialize error handler with logging and rendering capabilities
	}
	return errorHandler
}

func setupPage(conf config, errorHandler page.Error, meiliClient meilisearch.ServiceManager) page.Page {
	pageTempl := page.Page{
		Error: errorHandler,
		Home:  "/home/",
		NavItems: []page.NavItem{
			{Title: language.MakeLangString("Domů", "Home"), Path: "/home/", Skeleton: home.Skeleton, Indicator: "#home-skeleton"},
			{Title: language.MakeLangString("Hledání", "Search"), Path: "/courses/", Skeleton: courses.Skeleton, Indicator: "#courses-skeleton"},
			{Title: language.MakeLangString("Blueprint", "Blueprint"), Path: "/blueprint/", Skeleton: blueprint.Skeleton, Indicator: "#blueprint-skeleton"},
			{Title: language.MakeLangString("Studijní plán", "Degree plan"), Path: "/degreeplan/", Skeleton: degreeplan.Skeleton, Indicator: "#degreeplan-skeleton"},
		},
		Search: page.MeiliSearch{
			Client: meiliClient,
			Index:  "courses",
			Limit:  5,
		},
		Param:          "search",
		SearchEndpoint: "/courses/",
		ResultsDetailEndpoint: func(code string) string {
			return fmt.Sprintf("/course/%s", code)
		},
	}
	pageTempl.Init()

	return pageTempl
}

func setupHomeServer(conf config, errorHandler home.Error, pageTempl page.Page) http.Handler {
	home := home.Server{
		Auth:        cas.UserIDFromContext{},
		Error:       errorHandler,
		Page:        page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
		Recommender: fmt.Sprintf("http://%s:%d", conf.Recommender.Host, conf.Recommender.Port),
	}
	home.Init()
	return home.Router()
}

func setupBlueprintServer(db *sqlx.DB, errorHandler blueprint.Error, pageTempl page.Page) http.Handler {
	blueprint := blueprint.Server{
		Auth:  cas.UserIDFromContext{},
		Data:  blueprint.DBManager{DB: db},
		Error: errorHandler,
		Page:  page.PageWithNoFiltersAndForgetsSearchQueryOnRefresh{Page: pageTempl},
	}
	blueprint.Init()
	return blueprint.Router()
}

func setupCourseDetailServer(db *sqlx.DB, errorHandler coursedetail.Error, pageTempl page.Page, meiliClient meilisearch.ServiceManager) http.Handler {
	coursedetail := coursedetail.Server{
		Auth: cas.UserIDFromContext{},
		BpBtn: bpbtn.Add{
			DB:    db,
			Templ: bpbtn.AddBtn,
			Options: bpbtn.Options{
				HxPostBase: "/course",
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

func setupCoursesServer(db *sqlx.DB, errorHandler courses.Error, pageTempl page.Page, meiliClient meilisearch.ServiceManager) http.Handler {
	courses := courses.Server{
		Auth: cas.UserIDFromContext{},
		BpBtn: bpbtn.Add{
			DB:    db,
			Templ: bpbtn.AddBtn,
			Options: bpbtn.Options{
				HxPostBase: "/courses",
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

func setupDegreePlanServer(db *sqlx.DB, errorHandler degreeplan.Error, pageTempl page.Page, meiliClient meilisearch.ServiceManager) http.Handler {
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
	return degreePlan.Router()
}

func setupProtectedHandler(homeServer http.Handler, pageTempl page.Page, blueprintServer, coursedetailServer, coursesServer, degreePlanServer, static http.Handler, db *sqlx.DB, errorHandler cas.Error, conf config) http.Handler {
	protectedRouter := http.NewServeMux()
	protectedRouter.Handle("/", homeServer)
	handle(protectedRouter, "/page/", pageTempl.Router())
	handle(protectedRouter, "/blueprint/", blueprintServer)
	handle(protectedRouter, "/course/", coursedetailServer)
	handle(protectedRouter, "/courses/", coursesServer)
	handle(protectedRouter, "/degreeplan/", degreePlanServer)
	protectedRouter.Handle("GET /logo.svg", static)
	protectedRouter.Handle("GET /style.css", static)
	protectedRouter.Handle("GET /js/", static)

	authentication := cas.Authentication{
		Data:           cas.DBManager{DB: db},
		Error:          errorHandler,
		CAS:            cas.CAS{Host: conf.CAS.Host},
		AfterLoginPath: "/",
	}
	var protectedHandler http.Handler
	protectedHandler = protectedRouter
	protectedHandler = authentication.AuthenticateHTTP(protectedHandler)

	return protectedHandler
}

func setupUnprotectedHandler(protectedHandler http.Handler, static http.Handler) http.Handler {
	unprotectedRouter := http.NewServeMux()
	unprotectedRouter.Handle("/", protectedHandler)
	unprotectedRouter.Handle("GET /favicon.ico", static)
	unprotectedRouter.Handle("GET /logo.svg", static)

	var unprotectedHandler http.Handler
	unprotectedHandler = unprotectedRouter
	unprotectedHandler = language.SetAndStripLanguage(unprotectedHandler)
	unprotectedHandler = logging(unprotectedHandler)
	return unprotectedHandler
}

// ===========================

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
