package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	go_ora "github.com/sijms/go-ora/v2"
)

func main() {
	errp := func(e error) {
		if e != nil {
			log.Panicln(e)
		}
	}

	sis, err := createSISConn()
	errp(err)
	defer sis.Close()
	recsis, err := createRecSISConn()
	errp(err)
	defer recsis.Close()

	povinn := &extractPOVINN{
		from: sis,
		to:   recsis,
	}
	p := extractProcess{
		this: povinn,
	}
	_ = p.run()

}

type ExtractError struct {
	err  error
	proc *extractProcess
}

func (e ExtractError) Error() string {
	return fmt.Sprintf("ExtractError: %s", e.err.Error())
}

type operation interface {
	name() string
	selectData() error
	insertData() error
}

type extractProcess struct {
	this operation
	next *extractProcess
}

func (ep *extractProcess) run() ExtractError {
	log.Printf("Extracting %s from SIS", ep.this.name())
	err := ep.this.selectData()
	if err != nil {
		log.Println("ExtractPOVINN:", err)
		return ExtractError{err: err, proc: ep}
	}
	log.Println("Extracting POVINN to RecSIS")
	err = ep.this.insertData()
	if err != nil {
		log.Println("ExtractPOVINN:", err)
		return ExtractError{err: err, proc: ep}
	}
	if ep.next != nil {
		return ep.next.run()
	}
	return ExtractError{}
}

type extractPOVINN struct {
	from *sqlx.DB
	to   *sqlx.DB
	data struct {
		Code string `db:"POVINN"`
	}
}

func (ep *extractPOVINN) name() string {
	return "POVINN"
}

func (ep *extractPOVINN) selectData() error {
	query := "SELECT Povinn FROM POVINN FETCH FIRST 1 ROW ONLY"
	err := ep.from.Select(&ep.data, query)
	return err
}

func (ep *extractPOVINN) insertData() error {
	query := "INSERT INTO POVINN (Code) VALUES (:POVINN)"
	_, err := ep.to.NamedExec(query, ep.data)
	return err
}

func createSISConn() (*sqlx.DB, error) {
	log.Println("Connecting to SIS")
	const (
		host = "localhost"
		port = 10502
	)
	var (
		user    = os.Getenv("DB_USER")
		pass    = os.Getenv("DB_PASS")
		service = "studuk.prod"
	)
	connStr := go_ora.BuildUrl(host, port, service, user, pass, nil)
	conn, err := sqlx.Open("oracle", connStr)
	return conn, err
}

func createRecSISConn() (*sqlx.DB, error) {
	log.Println("Connecting to RecSIS")
	const (
		host   = "localhost"
		port   = 5432
		dbname = "recsis"
	)
	var (
		user = os.Getenv("RECSIS_USER")
		pass = os.Getenv("RECSIS_PASS")
	)
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
	conn, err := sqlx.Open("postgres", connStr)
	return conn, err
}
