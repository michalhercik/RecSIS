package main

import (
	"log"
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path)
	})
}

func main() {

	router := http.NewServeMux()

	//////////////////////////////////////////
	// Handlers
	//////////////////////////////////////////
	// eg.: router.HandleFunc("GET /", handleGetHello)

	//////////////////////////////////////////
	// Server setup
	//////////////////////////////////////////
	server := http.Server{
		Addr:    ":8000",
		Handler: Logging(router),
	}

	log.Println("Server starting ...")
	log.Println("http://localhost:8000/")

	server.ListenAndServe()
}
