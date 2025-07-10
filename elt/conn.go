package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/meilisearch/meilisearch-go"
	go_ora "github.com/sijms/go-ora/v2"
)

func createSISConn(conf config) (*sqlx.DB, error) {
	var (
		host    = conf.SIS.Host
		port    = conf.SIS.Port
		service = conf.SIS.Service
		user    = os.Getenv("DB_USER")
		pass    = os.Getenv("DB_PASS")
	)
	if len(user) == 0 || len(pass) == 0 {
		return nil, fmt.Errorf("DB_USER and DB_PASS environment variables must be set")
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
		// schema = conf.RecSIS.Schema
		user = os.Getenv("RECSIS_USER")
		pass = os.Getenv("RECSIS_PASS")
	)
	if len(user) == 0 || len(pass) == 0 {
		return nil, fmt.Errorf("RECSIS_USER and RECSIS_PASS environment variables must be set")
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
	conn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open RecSIS connection: %w", err)
	}
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping RecSIS database: %w", err)
	}
	log.Println("✅ Connection to RecSIS successfull")
	return conn, err
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
