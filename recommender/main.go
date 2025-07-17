package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path)
	})
}

var db *sqlx.DB

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

	//////////////////////////////////////////
	// Database setup
	//////////////////////////////////////////

	pass := os.Getenv("RECSIS_RECOMMENDER_DB_PASS")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		conf.Postgres.Host, conf.Postgres.Port, conf.Postgres.User, pass, conf.Postgres.DBName, conf.Postgres.Schema)
	db, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	//////////////////////////////////////////
	// Handlers
	//////////////////////////////////////////

	handler := http.NewServeMux()
	handler.HandleFunc("GET /recommended", getRecommendedCourses)
	handler.HandleFunc("GET /newest", getNewestCourses)

	server := http.Server{
		Addr:    ":8002",
		Handler: logging(handler),
	}

	log.Println("Server starting ...")
	log.Println("http://localhost:8002/")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

type config struct {
	Postgres struct {
		Host   string `toml:"host"`
		Port   int    `toml:"port"`
		User   string `toml:"user"`
		DBName string `toml:"dbname"`
		Schema string `toml:"schema"`
	} `toml:"postgres"`
}

func getRecommendedCourses(w http.ResponseWriter, r *http.Request) {
	getCourses(w, r, sql_recommended)
}

func getNewestCourses(w http.ResponseWriter, r *http.Request) {
	getCourses(w, r, sql_newest)
}
