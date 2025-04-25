package cas

import (
	"encoding/json"
	"net/http"
        "fmt"
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
        ticket := r.FormValue("ticket")
        fmt.Fprintln(w,"ticket: ", ticket)
        url := "https://acheron.ms.mff.cuni.cz:42050/cas/?format=json&ticket=" + ticket
        res, err := http.Get(url) 
        // fmt.Fprintln(w, "Validate through: ", url)
        if err != nil {
                fmt.Fprintln(w, err)
                return
        }
        var data map[string]any
        err = json.NewDecoder(res.Body).Decode(&data)
        if err != nil {
                fmt.Fprintln(w, err)
                return
        }
        fmt.Println(data)
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(r.URL.Query())
	json.NewEncoder(w).Encode(data)
}

func Cas() {
	// res, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	// if err != nil {
	// 	panic(err)
	// }
	// defer res.Body.Close()
}
