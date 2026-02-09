package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/meilisearch/meilisearch-go"
	go_ora "github.com/sijms/go-ora/v2"
)

func createSISConn(conf config) (*sqlx.DB, error) {
	var (
		host    = conf.SIS.Host
		port    = conf.SIS.Port
		service = conf.SIS.Service
		user    = os.Getenv("SIS_DB_USER")
		pass    = os.Getenv("SIS_DB_PASS")
	)
	if len(pass) == 0 {
		return nil, fmt.Errorf("SIS_DB_USER and SIS_DB_PASS environment variables must be set")
	}
	connStr := go_ora.BuildUrl(host, port, service, user, pass, nil)
	conn, err := sqlx.Open("oracle", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open SIS connection: %w", err)
	}
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SIS database: %w", err)
	}
	log.Println("✅ Connection to SIS successfull")
	return conn, err
}

func createRecSISConn(conf config) (*sqlx.DB, error) {
	var (
		host   = conf.RecSIS.Host
		port   = conf.RecSIS.Port
		dbname = conf.RecSIS.DBName
		user   = conf.RecSIS.User
		pass   = os.Getenv("RECSIS_ELT_DB_PASS")
	)
	if len(pass) == 0 {
		return nil, fmt.Errorf("RECSIS_ELT_DB_PASS environment variables must be set")
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
	conn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open RecSIS connection: %w", err)
	}
	applyRecSISPoolSettings(conn)
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping RecSIS database: %w", err)
	}
	log.Println("✅ Connection to RecSIS successfull")
	return conn, err
}

func applyRecSISPoolSettings(conn *sqlx.DB) {
	maxOpen := envInt("ELT_MAX_OPEN_CONNS", 10)
	maxIdle := envInt("ELT_MAX_IDLE_CONNS", 5)
	maxLifetime := envDuration("ELT_CONN_MAX_LIFETIME", 30*time.Minute)
	conn.SetMaxOpenConns(maxOpen)
	conn.SetMaxIdleConns(maxIdle)
	conn.SetConnMaxLifetime(maxLifetime)
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if parsed, err := time.ParseDuration(v); err == nil {
			return parsed
		}
	}
	return fallback
}

func createMeilisearchConn(conf config) (meilisearch.ServiceManager, error) {
	url := fmt.Sprintf("http://%s:%d", conf.MeiliSearch.Host, conf.MeiliSearch.Port)
	apiKey := os.Getenv("MEILI_MASTER_KEY")
	meili := meilisearch.New(url, meilisearch.WithAPIKey(apiKey))
	if !meili.IsHealthy() {
		return nil, fmt.Errorf("failed to connect to Meilisearch at %s", url)
	}
	return meili, nil
}
