package cas

import (
	"encoding/json"
	"log"
	"net/http"
)

type Server struct {
	router *http.ServeMux
}

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("/", s.cas)
	s.router = router
}

func (s *Server) Router() http.Handler {
	return s.router
}

func (s *Server) cas(w http.ResponseWriter, r *http.Request) {
	var data map[string]any
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("cas:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Cas() {
	// res, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	// if err != nil {
	// 	panic(err)
	// }
	// defer res.Body.Close()
}
