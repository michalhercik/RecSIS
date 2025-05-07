package cas

import (
	"encoding/json"
	"net/http"
        "net/url"
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
        // fmt.Fprintln(w,"ticket: ", ticket)
        fmt.Println("ticket: ", ticket)
        // url := "https://acheron.ms.mff.cuni.cz:42050/cas/?format=json&ticket=" + ticket
        // url := https://cas.cuni.cz/cas/serviceValidate?service=https%3A%2F%2Facheron.ms.mff.cuni.cz%3A42050%2Fcas%2F&format=json&ticket=ST-90478-nyhtZ0UOtTBx2OvMh4-0-LkcF9A-idp2
        validateReq := url.URL{
                Scheme: "https",
                Host: "cas.cuni.cz",
                Path: "/cas/serviceValidate",
                RawQuery: url.Values{
                        "format": []string{"json"},
                        "service": []string{"https://acheron.ms.mff.cuni.cz:42050/cas/"},
                        "ticket": []string{ticket},
                }.Encode(),
        }
        fmt.Println("validateReq: ", validateReq.String())
        res, err := http.Get(validateReq.String()) 
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
        // fmt.Println(data)
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
