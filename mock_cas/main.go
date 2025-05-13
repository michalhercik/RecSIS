package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path)
	})
}

func main() {
	const ticket = "ST-104269-lHanfdR85pp2PaEsfwv2VWe9UGM-idp2"
	var userID string
	handler := http.NewServeMux()
	handler.HandleFunc("GET /cas/login", func(w http.ResponseWriter, r *http.Request) {
		service := r.FormValue("service")
		if len(service) == 0 {
			http.Error(w, "Parameter service is missing.", http.StatusNotAcceptable)
			return
		}
		LoginForm(service).Render(r.Context(), w)
	})
	handler.HandleFunc("POST /cas/login", func(w http.ResponseWriter, r *http.Request) {
		userID = r.FormValue("user-id")
		serviceURL, err := url.Parse(r.FormValue("service"))
		if err != nil {
			http.Error(w, "Invalid service URL", http.StatusBadRequest)
			return
		}
		serviceURL.RawQuery = "ticket=" + ticket
		http.Redirect(w, r, serviceURL.String(), http.StatusFound)
	})
	handler.HandleFunc("GET /cas/serviceValidate/", func(w http.ResponseWriter, r *http.Request) {
		if ticket != r.FormValue("ticket") {
			http.Error(w, "Invalid ticket", http.StatusBadRequest)
			return
		}
		if r.FormValue("format") != "json" {
			http.Error(w, "Parameter format is missing or is not set to json.", http.StatusNotAcceptable)
			return
		}
		response := validationResponse{
			ServiceResponse: serviceResponse{
				AuthenticationSuccess: authenticationSuccess{
					Attributes: attributes{
						EduPersonScopedAffiliation: []string{"student@mff.cuni.cz"},
					},
					User: userID,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		payload, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		w.Write(payload)
	})

	server := http.Server{
		Addr:    ":8001", // DOCKER, PRODUCTION: when run as docker container remove localhost
		Handler: logging(handler),
	}

	log.Println("Server starting ...")
	log.Println("http://localhost:8001/")

	// err = server.ListenAndServeTLS("recsis-cert/fullchain.pem", "recsis-cert/privkey.pem")
	err := server.ListenAndServeTLS("../src/server.crt", "../src/server.key")
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

type validationResponse struct {
	ServiceResponse serviceResponse `json:"serviceResponse"`
}
type serviceResponse struct {
	AuthenticationSuccess authenticationSuccess `json:"authenticationSuccess"`
}
type authenticationSuccess struct {
	Attributes attributes `json:"attributes"`
	User       string     `json:"user"`
}
type attributes struct {
	// AuthenticationDate                     []float64 `json:"authenticationDate"`
	// AuthenticationMethod                   []string  `json:"authenticationMethod"`
	// ClientIPAddress                        []string  `json:"clientIpAddress"`
	// CN                                     []string  `json:"cn"`
	// CredentialType                         []string  `json:"credentialType"`
	// CuniPersonalID                         []string  `json:"cunipersonalid"`
	EduPersonScopedAffiliation []string `json:"edupersonscopedaffiliation"`
	// GivenName                              []string  `json:"givenname"`
	// IsFromNewLogin                         []bool    `json:"isFromNewLogin"`
	// LongTermAuthenticationRequestTokenUsed []bool    `json:"longTermAuthenticationRequestTokenUsed"`
	// Mail                                   []string  `json:"mail"`
	// SAMLAuthenticationStatementAuthMethod  []string  `json:"samlAuthenticationStatementAuthMethod"`
	// ServerIPAddress                        []string  `json:"serverIpAddress"`
	// SN                                     []string  `json:"sn"`
	// SuccessfulAuthenticationHandlers       []string  `json:"successfulAuthenticationHandlers"`
	// UID                                    []string  `json:"uid"`
	// UserAgent                              []string  `json:"userAgent"`
}
