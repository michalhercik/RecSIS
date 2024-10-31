package main

import (
	"log"
	"net/http"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path)
	})
}

func comingSoon(w http.ResponseWriter, r *http.Request) {
	component := comingSoonPage()
	component.Render(r.Context(), w)
}

func main() {

	router := http.NewServeMux()

	//////////////////////////////////////////
	// Handlers
	//////////////////////////////////////////

	router.HandleFunc("GET /", comingSoon)

	//////////////////////////////////////////
	// Server setup
	//////////////////////////////////////////
	server := http.Server{
		Addr:    ":8000",
		Handler: logging(router),
	}

	log.Println("Server starting ...")
	log.Println("http://localhost:8000/")

	server.ListenAndServe()
}
